package loraserver

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/brocaar/loraserver/models"
	"github.com/brocaar/lorawan"
	"github.com/brocaar/lorawan/band"
)

var (
	errDoesNotExist = errors.New("object does not exist")
)

// Server represents a LoRaWAN network-server.
type Server struct {
	ctx Context
	wg  sync.WaitGroup
}

// NewServer creates a new server.
func NewServer(ctx Context) *Server {
	return &Server{
		ctx: ctx,
	}
}

// Start starts the server.
func (s *Server) Start() error {
	log.WithFields(log.Fields{
		"net_id": s.ctx.NetID,
		"nwk_id": hex.EncodeToString([]byte{s.ctx.NetID.NwkID()}),
	}).Info("starting loraserver")

	go func() {
		s.wg.Add(1)
		handleRXPackets(s.ctx)
		s.wg.Done()
	}()
	go func() {
		s.wg.Add(1)
		handleTXPayloads(s.ctx)
		s.wg.Done()
	}()
	return nil
}

// Stop closes the gateway and application backends and waits for the
// server to complete the pending models.
func (s *Server) Stop() error {
	if err := s.ctx.Gateway.Close(); err != nil {
		return err
	}
	if err := s.ctx.Application.Close(); err != nil {
		return err
	}

	log.Info("waiting for pending packets to complete")
	s.wg.Wait()
	return nil
}

func handleTXPayloads(ctx Context) {
	var wg sync.WaitGroup
	for txPayload := range ctx.Application.TXPayloadChan() {
		go func(txPayload models.TXPayload) {
			wg.Add(1)
			if err := addTXPayloadToQueue(ctx.RedisPool, txPayload); err != nil {
				log.Errorf("could not add TXPayload to queue: %s", err)
			}
			wg.Done()
		}(txPayload)
	}
	wg.Wait()
}

func handleRXPackets(ctx Context) {
	var wg sync.WaitGroup
	for rxPacket := range ctx.Gateway.RXPacketChan() {
		go func(rxPacket models.RXPacket) {
			wg.Add(1)
			if err := handleRXPacket(ctx, rxPacket); err != nil {
				log.Errorf("error while processing RXPacket: %s", err)
			}
			wg.Done()
		}(rxPacket)
	}
	wg.Wait()
}

func handleRXPacket(ctx Context, rxPacket models.RXPacket) error {
	switch rxPacket.PHYPayload.MHDR.MType {
	case lorawan.JoinRequest:
		return validateAndCollectJoinRequestPacket(ctx, rxPacket)
	case lorawan.UnconfirmedDataUp, lorawan.ConfirmedDataUp:
		return validateAndCollectDataUpRXPacket(ctx, rxPacket)
	default:
		return fmt.Errorf("unknown MType: %v", rxPacket.PHYPayload.MHDR.MType)
	}
}

func validateAndCollectDataUpRXPacket(ctx Context, rxPacket models.RXPacket) error {
	// MACPayload must be of type *lorawan.MACPayload
	macPL, ok := rxPacket.PHYPayload.MACPayload.(*lorawan.MACPayload)
	if !ok {
		return fmt.Errorf("expected *lorawan.MACPayload, got: %T", rxPacket.PHYPayload.MACPayload)
	}

	// get the session data
	ns, err := getNodeSession(ctx.RedisPool, macPL.FHDR.DevAddr)
	if err != nil {
		return fmt.Errorf("could not get node-session: %s", err)
	}

	// validate and get the full int32 FCnt
	fullFCnt, ok := ns.ValidateAndGetFullFCntUp(macPL.FHDR.FCnt)
	if !ok {
		log.WithFields(log.Fields{
			"packet_fcnt": macPL.FHDR.FCnt,
			"server_fcnt": ns.FCntUp,
		}).Warning("invalid FCnt")
		return errors.New("invalid FCnt or too many dropped frames")
	}
	macPL.FHDR.FCnt = fullFCnt

	// validate MIC
	micOK, err := rxPacket.PHYPayload.ValidateMIC(ns.NwkSKey)
	if err != nil {
		return err
	}
	if !micOK {
		return errors.New("invalid MIC")
	}

	if macPL.FPort != nil {
		if *macPL.FPort == 0 {
			// decrypt FRMPayload with NwkSKey when FPort == 0
			if err := rxPacket.PHYPayload.DecryptFRMPayload(ns.NwkSKey); err != nil {
				return err
			}
		} else {
			if err := rxPacket.PHYPayload.DecryptFRMPayload(ns.AppSKey); err != nil {
				return err
			}
		}
	}

	return collectAndCallOnce(ctx.RedisPool, rxPacket, func(rxPackets RXPackets) error {
		return handleCollectedDataUpPackets(ctx, rxPackets)
	})
}

func handleCollectedDataUpPackets(ctx Context, rxPackets RXPackets) error {
	if len(rxPackets) == 0 {
		return errors.New("packet collector returned 0 packets")
	}
	rxPacket := rxPackets[0]

	var macs []string
	for _, p := range rxPackets {
		macs = append(macs, p.RXInfo.MAC.String())
	}

	rxPacketsCount := uint8(len(rxPackets))
	log.WithFields(log.Fields{
		"gw_count": rxPacketsCount,
		"gw_macs":  strings.Join(macs, ", "),
		"mtype":    rxPackets[0].PHYPayload.MHDR.MType,
	}).Info("packet(s) collected")

	macPL, ok := rxPacket.PHYPayload.MACPayload.(*lorawan.MACPayload)
	if !ok {
		return fmt.Errorf("expected *lorawan.MACPayload, got: %T", rxPacket.PHYPayload.MACPayload)
	}

	ns, err := getNodeSession(ctx.RedisPool, macPL.FHDR.DevAddr)
	if err != nil {
		return fmt.Errorf("could not get node-session: %s", err)
	}

	if macPL.FPort != nil && *macPL.FPort > 0 {
		var data []byte

		// it is possible that the FRMPayload is empty, in this case only
		// the FPort will be send
		if len(macPL.FRMPayload) == 1 {
			dataPL, ok := macPL.FRMPayload[0].(*lorawan.DataPayload)
			if !ok {
				return errors.New("FRMPayload must be of type *lorawan.DataPayload")
			}
			data = dataPL.Bytes
		}

		err = ctx.Application.Send(ns.DevEUI, ns.AppEUI, models.RXPayload{
			DevEUI:       ns.DevEUI,
			GatewayCount: rxPacketsCount,
			FPort:        *macPL.FPort,
			RSSI:         rxPacket.RXInfo.RSSI,
			Data:         data,
		})
		if err != nil {
			return fmt.Errorf("could not send RXPacket to application: %s", err)
		}
	}

	// sync counter with that of the device
	ns.FCntUp = macPL.FHDR.FCnt
	if err := saveNodeSession(ctx.RedisPool, ns); err != nil {
		return fmt.Errorf("could not update node-session: %s", err)
	}

	frmPayload, fOpts, err := handleMacCommands(ctx, ns, macPL, rxPacket.RXInfo.LoRaSNR, uint8(rxPacketsCount))
	if err != nil {
		return fmt.Errorf("Could not process mac commands: %s", err)
	}

	// handle downlink (ACK)
	time.Sleep(CollectDataDownWait)
	return handleDataDownReply(ctx, rxPacket, ns, frmPayload, fOpts)
}

func updateNodeSession(ctx Context, ns models.NodeSession, macCommandsWithResponse *MacCommandsWithResponse) error {
	if macCommandsWithResponse == nil {
		return nil
	}
	if macCommandsWithResponse.Response.RXTimingSetup {
		ns.TXParams.TXDelay1 = macCommandsWithResponse.MacCommands.RXTimingSetup.Delay
		if err := saveNodeSession(ctx.RedisPool, ns); err != nil {
			return fmt.Errorf("Could not update node session with TXDelay1: %s", err)
		}
	}

	rx2Setup := macCommandsWithResponse.Response.RX2Setup
	haveNewRX2Setup := rx2Setup != nil && rx2Setup.RX1DRoffsetACK && rx2Setup.RX2DataRateACK && rx2Setup.ChannelACK
	if haveNewRX2Setup {
		ns.TXParams.Set = true
		ns.TXParams.TX2Frequency = macCommandsWithResponse.MacCommands.RX2Setup.Frequency * 100
		ns.TXParams.TX1DRoffset = macCommandsWithResponse.MacCommands.RX2Setup.DLsettings.RX1DRoffset
		ns.TXParams.TX2DataRate = macCommandsWithResponse.MacCommands.RX2Setup.DLsettings.RX2DataRate
		if err := saveNodeSession(ctx.RedisPool, ns); err != nil {
			return fmt.Errorf("Could not update node session with RX2Setup: %s", err)
		}
	}

	return nil
}

func handleMACCommandLinkCheckReq(rxPacketsCount uint8, margin uint8) lorawan.MACCommand {
	return lorawan.MACCommand{
		CID: lorawan.LinkCheckAns,
		Payload: &lorawan.LinkCheckAnsPayload{
			Margin: margin,
			GwCnt:  rxPacketsCount,
		},
	}
}

// TODO: Remove the hard-coded values from here once they are available from
// lorawan lib.
func calculateMacCommandsSize(macCommands []lorawan.MACCommand) (int, error) {
	sizes := map[byte]int{
		byte(lorawan.LinkCheckAns):     2,
		byte(lorawan.LinkADRReq):       4,
		byte(lorawan.DutyCycleReq):     1,
		byte(lorawan.RXParamSetupReq):  4,
		byte(lorawan.NewChannelReq):    4,
		byte(lorawan.RXTimingSetupReq): 1,
	}

	size := 0
	for _, mc := range macCommands {
		size += sizes[byte(mc.CID)]
		size++
	}

	return size, nil
}

func handleMacCommands(ctx Context,
	ns models.NodeSession,
	macPL *lorawan.MACPayload,
	loraSNR float64,
	rxPacketsCount uint8) ([]lorawan.Payload, []lorawan.MACCommand, error) {
	var macCommandsSlice []lorawan.MACCommand
	if macPL.FPort != nil && *macPL.FPort == 0 {
		for _, payload := range macPL.FRMPayload {
			macCommand, ok := payload.(*lorawan.MACCommand)
			if !ok {
				return nil, nil, fmt.Errorf("Unexpected payload found with FPort == 0: %s", payload)
			}
			macCommandsSlice = append(macCommandsSlice, *macCommand)
		}
	}
	var margin uint8
	if loraSNR > 0 {
		// Maybe this is not the correct way to handle
		margin = uint8(loraSNR)
	}

	macCommandsSlice = append(macCommandsSlice, macPL.FHDR.FOpts...)

	macCommandsResponse, unparsedMacCommands, err := parseMacCommandsResponse(macCommandsSlice)
	if err != nil {
		return nil, nil, err
	}

	if macCommandsResponse.nonEmpty {
		if err = updateSentCommandResponse(ctx.RedisPool, ns.DevEUI, macCommandsResponse); err != nil {
			return nil, nil, err
		}
	}

	macCommandsWithResponse, err := getLastMacCommandsWithResponse(ctx.RedisPool, ns.DevEUI)
	if err != nil {
		return nil, nil, err
	}

	if err = updateNodeSession(ctx, ns, macCommandsWithResponse); err != nil {
		return nil, nil, err
	}

	macCommands, err := getMacCommandsToBeSent(ctx.RedisPool, ns.DevEUI)
	if err != nil && err != errDoesNotExist {
		return nil, nil, fmt.Errorf("Could not get MAC commands: %s", err)
	}

	var macCommandsResponses []lorawan.MACCommand
	if err != errDoesNotExist {
		macCommandsResponses = macCommands.ToMACCommands()

		if err = markMacCommandsAsSent(ctx.RedisPool, ns.DevEUI, macCommands); err != nil {
			return nil, nil, fmt.Errorf("Could not mark MAC commands as waiting response: %s", err)
		}
	}
	for _, macCommand := range unparsedMacCommands {
		if macCommand.CID == lorawan.LinkCheckReq {
			macCommandsResponses = append(macCommandsResponses, handleMACCommandLinkCheckReq(rxPacketsCount, margin))
		} else {
			log.WithFields(log.Fields{
				"macCommand": macCommand,
			}).Warn("Unknown mac command received")
		}
	}

	forceZeroFPort := macCommands.RequestEncryption
	if !forceZeroFPort {
		var macCommandsSize int
		if macCommandsSize, err = calculateMacCommandsSize(macCommandsResponses); err != nil {
			return nil, nil, err
		}
		// TODO: extract in constant
		forceZeroFPort = macCommandsSize > 15
	}

	var frmPayload []lorawan.Payload
	if forceZeroFPort {
		for i := range macCommandsResponses {
			frmPayload = append(frmPayload, &macCommandsResponses[i])
		}

		macCommandsResponses = nil
	}

	return frmPayload, macCommandsResponses, nil
}

func handleDataDownReply(ctx Context,
	rxPacket models.RXPacket,
	ns models.NodeSession,
	frmPayload []lorawan.Payload,
	fOpts []lorawan.MACCommand) error {
	macPL, ok := rxPacket.PHYPayload.MACPayload.(*lorawan.MACPayload)
	if !ok {
		return fmt.Errorf("expected *lorawan.MACPayload, got: %T", rxPacket.PHYPayload.MACPayload)
	}

	// the last payload was received by the node
	if macPL.FHDR.FCtrl.ACK {
		if err := clearInProcessTXPayload(ctx.RedisPool, ns.DevEUI); err != nil {
			return fmt.Errorf("could not clear in-process TXPayload: %s", err)
		}
		ns.FCntDown++
		if err := saveNodeSession(ctx.RedisPool, ns); err != nil {
			return fmt.Errorf("could not update node-session: %s", err)
		}
	}

	var remaining bool

	var confirmed = frmPayload != nil || fOpts != nil
	var fPort *uint8
	var err error

	// set ACK to true when received packet needs an ACK
	ack := rxPacket.PHYPayload.MHDR.MType == lorawan.ConfirmedDataUp

	if frmPayload != nil {
		tmp := uint8(0)
		fPort = &tmp
	} else {
		// check if there are payloads pending in the queue
		var txPayload models.TXPayload
		txPayload, remaining, err = getTXPayloadAndRemainingFromQueue(ctx.RedisPool, ns.DevEUI)

		// errDoesNotExist should not be handled as an error, it just means there
		// is no queue / the queue is empty
		if err != nil && err != errDoesNotExist {
			return fmt.Errorf("could not get TXPayload from queue: %s", err)
		}

		if err == nil {
			// we have a pending packet
			frmPayload = append(frmPayload, &lorawan.DataPayload{Bytes: txPayload.Data})
			confirmed = txPayload.Confirmed
			fPort = &txPayload.FPort
		} else if fOpts == nil && !ack {
			// nothing pending in the queues and no need to ACK RXPacket
			return nil
		}
	}

	if frmPayload == nil && fOpts == nil && !ack {
		return nil
	}

	macPLReply, err := createMACPayload(frmPayload, fOpts, remaining, fPort, ns, ack)
	if err != nil {
		return fmt.Errorf("Error constructing MACPayload: %v", err)
	}

	if err = send(macPLReply, rxPacket, ns, ctx, confirmed); err != nil {
		return err
	}

	// remove the payload from the queue when not confirmed
	if !confirmed {
		if err := clearInProcessTXPayload(ctx.RedisPool, ns.DevEUI); err != nil {
			return err
		}
	}

	return nil
}

func createMACPayload(frmPayload []lorawan.Payload,
	fOpts []lorawan.MACCommand,
	remaining bool,
	fPort *uint8,
	ns models.NodeSession,
	ack bool) (*lorawan.MACPayload, error) {
	macPL := &lorawan.MACPayload{
		FHDR: lorawan.FHDR{
			DevAddr: ns.DevAddr,
			FCtrl: lorawan.FCtrl{
				FPending: remaining,
				ACK:      ack,
			},
			FCnt:  ns.FCntDown,
			FOpts: fOpts,
		},
		FPort:      fPort,
		FRMPayload: frmPayload,
	}

	return macPL, nil
}

func send(macPL *lorawan.MACPayload, rxPacket models.RXPacket, ns models.NodeSession, ctx Context, confirmed bool) error {
	phy := lorawan.PHYPayload{
		MHDR: lorawan.MHDR{
			MType: lorawan.UnconfirmedDataDown,
			Major: lorawan.LoRaWANR1,
		},
		MACPayload: macPL,
	}
	if confirmed {
		phy.MHDR.MType = lorawan.ConfirmedDataDown
	}
	if macPL.FPort != nil {
		if *macPL.FPort == 0 {
			if err := phy.EncryptFRMPayload(ns.NwkSKey); err != nil {
				return fmt.Errorf("could not encrypt FRMPayload(NwkSKey): %s", err)
			}
		} else {
			if err := phy.EncryptFRMPayload(ns.AppSKey); err != nil {
				return fmt.Errorf("could not encrypt FRMPayload: %s", err)
			}
		}
	}

	if err := phy.SetMIC(ns.NwkSKey); err != nil {
		return fmt.Errorf("could not set MIC: %s", err)
	}

	receiveDelay1 := band.ReceiveDelay1
	if ns.TXParams.TXDelay1 != 0 {
		receiveDelay1 = time.Duration(ns.TXParams.TXDelay1) * time.Second
	}

	dataRate := rxPacket.RXInfo.DataRate
	if ns.TXParams.Set {
		var err error
		if dataRate, err = dataRateOffset(dataRate, ns.TXParams.TX1DRoffset); err != nil {
			return err
		}
	}

	txPacket := models.TXPacket{
		TXInfo: models.TXInfo{
			MAC:       rxPacket.RXInfo.MAC,
			Timestamp: rxPacket.RXInfo.Timestamp + uint32(receiveDelay1/time.Microsecond),
			Frequency: rxPacket.RXInfo.Frequency,
			Power:     band.DefaultTXPower,
			DataRate:  dataRate,
			CodeRate:  rxPacket.RXInfo.CodeRate,
		},
		PHYPayload: phy,
	}

	// window 1
	if err := ctx.Gateway.Send(txPacket); err != nil {
		return fmt.Errorf("sending TXPacket to the gateway failed: %s", err)
	}

	// specified in the lorawan specification
	receiveDelay2 := receiveDelay1 + time.Second

	// window 2
	if ns.TXParams.Set {
		txPacket.TXInfo.Frequency = ns.TXParams.TX2Frequency
		dataRateIndex := int(ns.TXParams.TX2DataRate)
		if dataRateIndex < len(band.DataRateConfiguration) {
			txPacket.TXInfo.DataRate = band.DataRateConfiguration[dataRateIndex]
		} else {
			return fmt.Errorf("Unknown data rate numeric value: %d", dataRateIndex)
		}
	} else {
		txPacket.TXInfo.Frequency = band.RX2Frequency
		txPacket.TXInfo.DataRate = band.DataRateConfiguration[band.RX2DataRate]
	}
	txPacket.TXInfo.Timestamp = rxPacket.RXInfo.Timestamp + uint32(receiveDelay2/time.Microsecond)

	if err := ctx.Gateway.Send(txPacket); err != nil {
		return fmt.Errorf("sending TXPacket to the gateway failed: %s", err)
	}

	// increment the FCntDown when MType != ConfirmedDataDown. In case of
	// ConfirmedDataDown we increment on ACK.
	if phy.MHDR.MType != lorawan.ConfirmedDataDown {
		ns.FCntDown++
		if err := saveNodeSession(ctx.RedisPool, ns); err != nil {
			return fmt.Errorf("could not update node-session: %s", err)
		}
	}

	return nil
}

func dataRateOffset(dataRate band.DataRate, offset uint8) (band.DataRate, error) {
	dataRateIndex := -1
	for index, currentDataRate := range band.DataRateConfiguration {
		if dataRate == currentDataRate {
			dataRateIndex = index
			break
		}
	}

	return dataRateByIndex(dataRateIndex - int(offset))
}

func dataRateByIndex(index int) (band.DataRate, error) {
	if index < 0 || index >= len(band.DataRateConfiguration) {
		return band.DataRate{}, fmt.Errorf("Unknown data rate index: %d", index)
	}

	return band.DataRateConfiguration[index], nil
}

func validateAndCollectJoinRequestPacket(ctx Context, rxPacket models.RXPacket) error {
	// MACPayload must be of type *lorawan.JoinRequestPayload
	jrPL, ok := rxPacket.PHYPayload.MACPayload.(*lorawan.JoinRequestPayload)
	if !ok {
		return fmt.Errorf("expected *lorawan.JoinRequestPayload, got: %T", rxPacket.PHYPayload.MACPayload)
	}

	// get node information for this DevEUI
	node, err := getNode(ctx.DB, jrPL.DevEUI)
	if err != nil {
		return fmt.Errorf("could not get node: %s", err)
	}

	// validate the MIC
	ok, err = rxPacket.PHYPayload.ValidateMIC(node.AppKey)
	if err != nil {
		return fmt.Errorf("could not validate MIC: %s", err)
	}
	if !ok {
		return errors.New("invalid mic")
	}

	return collectAndCallOnce(ctx.RedisPool, rxPacket, func(rxPackets RXPackets) error {
		return handleCollectedJoinRequestPackets(ctx, rxPackets)
	})

}

func handleCollectedJoinRequestPackets(ctx Context, rxPackets RXPackets) error {
	if len(rxPackets) == 0 {
		return errors.New("packet collector returned 0 packets")
	}
	rxPacket := rxPackets[0]

	var macs []string
	for _, p := range rxPackets {
		macs = append(macs, p.RXInfo.MAC.String())
	}

	log.WithFields(log.Fields{
		"gw_count": len(rxPackets),
		"gw_macs":  strings.Join(macs, ", "),
		"mtype":    rxPackets[0].PHYPayload.MHDR.MType,
	}).Info("packet(s) collected")

	// MACPayload must be of type *lorawan.JoinRequestPayload
	jrPL, ok := rxPacket.PHYPayload.MACPayload.(*lorawan.JoinRequestPayload)
	if !ok {
		return fmt.Errorf("expected *lorawan.JoinRequestPayload, got: %T", rxPacket.PHYPayload.MACPayload)
	}

	// get node information for this DevEUI
	node, err := getNode(ctx.DB, jrPL.DevEUI)
	if err != nil {
		return fmt.Errorf("could not get node: %s", err)
	}

	// validate the given nonce
	if !node.ValidateDevNonce(jrPL.DevNonce) {
		return errors.New("the given dev-nonce has already been used before")
	}

	// get random (free) DevAddr
	devAddr, err := getRandomDevAddr(ctx.RedisPool, ctx.NetID)
	if err != nil {
		return fmt.Errorf("could not get random DevAddr: %s", err)
	}

	// get app nonce
	appNonce, err := getAppNonce()
	if err != nil {
		return fmt.Errorf("could not get AppNonce: %s", err)
	}

	// get keys
	nwkSKey, err := getNwkSKey(node.AppKey, ctx.NetID, appNonce, jrPL.DevNonce)
	if err != nil {
		return fmt.Errorf("could not get NwkSKey: %s", err)
	}
	appSKey, err := getAppSKey(node.AppKey, ctx.NetID, appNonce, jrPL.DevNonce)
	if err != nil {
		return fmt.Errorf("could not get AppSKey: %s", err)
	}

	ns := models.NodeSession{
		DevAddr:  devAddr,
		DevEUI:   jrPL.DevEUI,
		AppSKey:  appSKey,
		NwkSKey:  nwkSKey,
		FCntUp:   0,
		FCntDown: 0,

		AppEUI: node.AppEUI,
	}
	if err = saveNodeSession(ctx.RedisPool, ns); err != nil {
		return fmt.Errorf("could not save node-session: %s", err)
	}

	// update the node (with updated used dev-nonces)
	if err = updateNode(ctx.DB, node); err != nil {
		return fmt.Errorf("could not update the node: %s", err)
	}

	// construct the lorawan packet
	phy := lorawan.PHYPayload{
		MHDR: lorawan.MHDR{
			MType: lorawan.JoinAccept,
			Major: lorawan.LoRaWANR1,
		},
		MACPayload: &lorawan.JoinAcceptPayload{
			AppNonce: appNonce,
			NetID:    ctx.NetID,
			DevAddr:  devAddr,
		},
	}
	if err = phy.SetMIC(node.AppKey); err != nil {
		return fmt.Errorf("could not set MIC: %s", err)
	}
	if err = phy.EncryptJoinAcceptPayload(node.AppKey); err != nil {
		return fmt.Errorf("could not encrypt join-accept: %s", err)
	}

	txPacket := models.TXPacket{
		TXInfo: models.TXInfo{
			MAC:       rxPacket.RXInfo.MAC,
			Timestamp: rxPacket.RXInfo.Timestamp + uint32(band.JoinAcceptDelay1/time.Microsecond),
			Frequency: rxPacket.RXInfo.Frequency,
			Power:     band.DefaultTXPower,
			DataRate:  rxPacket.RXInfo.DataRate,
			CodeRate:  rxPacket.RXInfo.CodeRate,
		},
		PHYPayload: phy,
	}

	// window 1
	if err = ctx.Gateway.Send(txPacket); err != nil {
		return fmt.Errorf("sending TXPacket to the gateway failed: %s", err)
	}

	// window 2
	txPacket.TXInfo.Timestamp = rxPacket.RXInfo.Timestamp + uint32(band.JoinAcceptDelay2/time.Microsecond)
	txPacket.TXInfo.Frequency = band.RX2Frequency
	txPacket.TXInfo.DataRate = band.DataRateConfiguration[band.RX2DataRate]
	if err = ctx.Gateway.Send(txPacket); err != nil {
		return fmt.Errorf("sending TXPacket to the gateway failed: %s", err)
	}

	return nil
}
