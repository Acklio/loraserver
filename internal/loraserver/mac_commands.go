package loraserver

import (
	"bytes"
	"encoding/gob"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/brocaar/lorawan"
	"github.com/garyburd/redigo/redis"
)

const (
	macCommandsKeyTempl          = "mac_cmds_2_be_sent_%s"
	macCommandsWithResponseTempl = "mac_cmds_last_with_response_%s"
)

//
type MacCommands struct {
	DevEUI        lorawan.EUI64
	LinkADR       lorawan.LinkADRReqPayload
	DutyCycle     lorawan.DutyCycleReqPayload
	RX2Setup      lorawan.RX2SetupReqPayload
	DevStatus     bool
	NewChannel    lorawan.NewChannelReqPayload
	RXTimingSetup lorawan.RXTimingSetupReqPayload
}

type MacCommandsResponse struct {
	LinkADR       *lorawan.LinkADRAnsPayload
	DutyCycle     bool
	RX2Setup      *lorawan.RX2SetupAnsPayload
	DevStatus     *lorawan.DevStatusAnsPayload
	NewChannel    *lorawan.NewChannelAnsPayload
	RXTimingSetup bool
	nonEmpty      bool
}

type MacCommandsWithResponse struct {
	MacCommands MacCommands
	Response    MacCommandsResponse
}

type NextAndLastMacCommands struct {
	Next *MacCommands
	Last *MacCommandsWithResponse
}

func (m *MacCommands) MacCommandSlice() []lorawan.Payload {
	var macCommands []lorawan.Payload
	if m.LinkADR != (lorawan.LinkADRReqPayload{}) {
		macCommands = append(macCommands, &lorawan.MACCommand{
			CID:     lorawan.LinkADRReq,
			Payload: &m.LinkADR,
		})
	}
	if m.DutyCycle != (lorawan.DutyCycleReqPayload{}) {
		macCommands = append(macCommands, &lorawan.MACCommand{
			CID:     lorawan.DutyCycleReq,
			Payload: &m.DutyCycle,
		})
	}
	if m.RX2Setup != (lorawan.RX2SetupReqPayload{}) {
		macCommands = append(macCommands, &lorawan.MACCommand{
			CID:     lorawan.RXParamSetupReq,
			Payload: &m.RX2Setup,
		})
	}
	if m.DevStatus {
		macCommands = append(macCommands, &lorawan.MACCommand{CID: lorawan.DevStatusReq})
	}
	if m.NewChannel != (lorawan.NewChannelReqPayload{}) {
		macCommands = append(macCommands, &lorawan.MACCommand{
			CID:     lorawan.NewChannelReq,
			Payload: &m.NewChannel,
		})
	}
	if m.RXTimingSetup != (lorawan.RXTimingSetupReqPayload{}) {
		macCommands = append(macCommands, &lorawan.MACCommand{
			CID:     lorawan.RXTimingSetupReq,
			Payload: &m.RXTimingSetup,
		})
	}

	return macCommands
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

func extractLinkADRAns(payload lorawan.MACCommandPayload) (*lorawan.LinkADRAnsPayload, error) {
	linkADRAnsPayload, ok := payload.(*lorawan.LinkADRAnsPayload)

	if !ok {
		return linkADRAnsPayload, fmt.Errorf("Unexpected payload given to extractLinkADRAns: %T", payload)
	}

	return linkADRAnsPayload, nil
}

func extractRXParamSetupAns(payload lorawan.MACCommandPayload) (*lorawan.RX2SetupAnsPayload, error) {
	rx2SetupAnsPayload, ok := payload.(*lorawan.RX2SetupAnsPayload)

	if !ok {
		return rx2SetupAnsPayload, fmt.Errorf("Unexpected payload given to extractRXParamSetupAns: %T", payload)
	}

	return rx2SetupAnsPayload, nil
}

func extractDevStatusAns(payload lorawan.MACCommandPayload) (*lorawan.DevStatusAnsPayload, error) {
	devStatusAnsPayload, ok := payload.(*lorawan.DevStatusAnsPayload)

	if !ok {
		return devStatusAnsPayload, fmt.Errorf("Unexpected payload given to extractDevStatusAns: %T", payload)
	}

	return devStatusAnsPayload, nil
}

func extractNewChannelAns(payload lorawan.MACCommandPayload) (*lorawan.NewChannelAnsPayload, error) {
	newChannelAnsPayload, ok := payload.(*lorawan.NewChannelAnsPayload)

	if !ok {
		return newChannelAnsPayload, fmt.Errorf("Unexpected payload given to extractNewChannelAns: %T", payload)
	}

	return newChannelAnsPayload, nil
}

// Convert a slice of lorawanMACCommand to a MacCommandsResponse and possibly a reply for LinkCheckReq
func fromMacCommandSlice(macCommands []lorawan.MACCommand, rxPacketsCount uint8, margin uint8) (MacCommandsResponse, []lorawan.MACCommand, error) {

	var repliesToSent []lorawan.MACCommand
	var macCommandsResponse = MacCommandsResponse{}
	var err error
	for _, macCommand := range macCommands {
		switch macCommand.CID {
		case lorawan.LinkCheckReq:
			repliesToSent = append(repliesToSent, handleMACCommandLinkCheckReq(rxPacketsCount, margin))
		case lorawan.DutyCycleAns:
			macCommandsResponse.DutyCycle = true
			macCommandsResponse.nonEmpty = true
		case lorawan.RXTimingSetupAns:
			macCommandsResponse.RXTimingSetup = true
			macCommandsResponse.nonEmpty = true
		case lorawan.LinkADRAns:
			if macCommandsResponse.LinkADR, err = extractLinkADRAns(macCommand.Payload); err != nil {
				return macCommandsResponse, nil, err
			}
			macCommandsResponse.nonEmpty = true
		case lorawan.RXParamSetupAns:
			if macCommandsResponse.RX2Setup, err = extractRXParamSetupAns(macCommand.Payload); err != nil {
				return macCommandsResponse, nil, err
			}
			macCommandsResponse.nonEmpty = true
		case lorawan.DevStatusAns:
			if macCommandsResponse.DevStatus, err = extractDevStatusAns(macCommand.Payload); err != nil {
				return macCommandsResponse, nil, err
			}
			macCommandsResponse.nonEmpty = true
		case lorawan.NewChannelAns:
			if macCommandsResponse.NewChannel, err = extractNewChannelAns(macCommand.Payload); err != nil {
				return macCommandsResponse, nil, err
			}
			macCommandsResponse.nonEmpty = true
		default:
			log.Warn("MAC Command %d not implemented", macCommand.CID)
		}
	}

	return macCommandsResponse, repliesToSent, nil
}

func markMacCommandsAsSent(p *redis.Pool, devEUI lorawan.EUI64, macCommands MacCommands) error {
	macCommandsWithResponse := MacCommandsWithResponse{
		MacCommands: macCommands,
	}

	return saveInRedisMacCommandsWithResponce(p, &macCommandsWithResponse)
}

func updateSentCommandResponse(p *redis.Pool, devEUI lorawan.EUI64, response MacCommandsResponse) error {
	macCommandsWithResponse, err := getLastMacCommandsWithResponse(p, devEUI)
	if err != nil {
		return err
	}

	if macCommandsWithResponse == nil {
		return errDoesNotExist
	}

	macCommandsWithResponse.Response = response

	return saveInRedisMacCommandsWithResponce(p, macCommandsWithResponse)
}

func saveInRedisMacCommandsWithResponce(p *redis.Pool, macCommandsWithResponse *MacCommandsWithResponse) error {
	devEUI := macCommandsWithResponse.MacCommands.DevEUI
	key := fmt.Sprintf(macCommandsWithResponseTempl, devEUI)

	c := p.Get()
	defer c.Close()

	return saveInRedis(c, key, macCommandsWithResponse)
}

func getLastMacCommandsWithResponse(p *redis.Pool, devEUI lorawan.EUI64) (*MacCommandsWithResponse, error) {
	key := fmt.Sprintf(macCommandsWithResponseTempl, devEUI)

	c := p.Get()
	defer c.Close()

	b, err := redis.Bytes(c.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	if err == redis.ErrNil {
		return nil, nil
	}

	var macCommandsWithResponse MacCommandsWithResponse

	if err = gob.NewDecoder(bytes.NewReader(b)).Decode(&macCommandsWithResponse); err != nil {
		return nil, err
	}

	return &macCommandsWithResponse, nil
}

func saveInRedis(conn redis.Conn, key string, data interface{}) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return err
	}

	_, err := redis.String(conn.Do("SET", key, buf.Bytes()))

	return err
}

func getMacCommandsToBeSent(p *redis.Pool, devEUI lorawan.EUI64) (MacCommands, error) {
	var macCommands = MacCommands{}

	c := p.Get()
	defer c.Close()
	key := fmt.Sprintf(macCommandsKeyTempl, devEUI)

	b, err := redis.Bytes(c.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		return macCommands, err
	}
	if err == redis.ErrNil {
		return macCommands, errDoesNotExist
	}

	macCommands.DevEUI = devEUI
	err = gob.NewDecoder(bytes.NewReader(b)).Decode(&macCommands)

	return macCommands, err
}

func getNextAndLastMacCommands(p *redis.Pool, devEUI lorawan.EUI64) (NextAndLastMacCommands, error) {
	nextCommands, err := getMacCommandsToBeSent(p, devEUI)
	if err != nil && err != errDoesNotExist {
		return NextAndLastMacCommands{}, err
	}
	lastMacCommands, err := getLastMacCommandsWithResponse(p, devEUI)
	if err != nil {
		return NextAndLastMacCommands{}, err
	}

	result := NextAndLastMacCommands{
		Next: &nextCommands,
		Last: lastMacCommands,
	}

	return result, nil
}

func updateMacCommandsToBeSent(p *redis.Pool, data MacCommands) (lorawan.EUI64, error) {
	c := p.Get()
	defer c.Close()
	if err := saveInRedis(c, fmt.Sprintf(macCommandsKeyTempl, data.DevEUI), data); err != nil {
		return lorawan.EUI64{}, err
	}

	return data.DevEUI, nil
}

//
type NodeMacCommandsAPI struct {
	ctx Context
}

//
func NewNodeMacCommandsAPI(ctx Context) *NodeMacCommandsAPI {
	return &NodeMacCommandsAPI{
		ctx: ctx,
	}
}

//
func (a *NodeMacCommandsAPI) GetMacCmdsToBeSent(devEUI lorawan.EUI64, macCommands *NextAndLastMacCommands) error {
	var err error
	*macCommands, err = getNextAndLastMacCommands(a.ctx.RedisPool, devEUI)
	return err
}

func (a *NodeMacCommandsAPI) Update(data MacCommands, devEUI *lorawan.EUI64) error {
	var err error
	*devEUI, err = updateMacCommandsToBeSent(a.ctx.RedisPool, data)

	return err
}
