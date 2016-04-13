package loraserver

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/brocaar/lorawan"
	"github.com/garyburd/redigo/redis"
)

const (
	macCommandsKeyTempl          = "mac_cmds_2_be_sent_%s"
	macCommandsWithResponseTempl = "mac_cmds_last_with_response_%s"
)

// MacCommands groups together the MAC commands that can be sent to an
// end-device.
type MacCommands struct {
	DevEUI        lorawan.EUI64
	LinkADR       lorawan.LinkADRReqPayload
	DutyCycle     lorawan.DutyCycleReqPayload
	RX2Setup      lorawan.RX2SetupReqPayload
	DevStatus     bool
	NewChannel    lorawan.NewChannelReqPayload
	RXTimingSetup lorawan.RXTimingSetupReqPayload
}

// MacCommandsResponse groups together the answers that an end-device can give
// to MacCommands.
type MacCommandsResponse struct {
	LinkADR       *lorawan.LinkADRAnsPayload
	DutyCycle     bool
	RX2Setup      *lorawan.RX2SetupAnsPayload
	DevStatus     *lorawan.DevStatusAnsPayload
	NewChannel    *lorawan.NewChannelAnsPayload
	RXTimingSetup bool
	nonEmpty      bool
}

// MacCommandsWithResponse groups together a MacCommands (request towards an
// end-device) and MacCommandsResponse (answer to the request).
type MacCommandsWithResponse struct {
	MacCommands MacCommands
	Response    MacCommandsResponse
}

// NextAndLastMacCommands groups together the next MacCommands to be sent to an
// end device and the last sent one with its reply (next and last might contain
// the same macCommands if there is still no answer).
type NextAndLastMacCommands struct {
	Next *MacCommands
	Last *MacCommandsWithResponse
}

// Payloads converts the MacCommands into a slice of Payloads.
func (m *MacCommands) Payloads() []lorawan.Payload {
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

// Convert a slice of lorawan.MACCommand to a MacCommandsResponse and possibly unparsed mac commands
func parseMacCommandsResponse(macCommands []lorawan.MACCommand) (MacCommandsResponse, []lorawan.MACCommand, error) {

	var unparsed []lorawan.MACCommand
	var macCommandsResponse = MacCommandsResponse{}
	var err error
	for _, macCommand := range macCommands {
		switch macCommand.CID {
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
			unparsed = append(unparsed, macCommand)
		}
	}

	return macCommandsResponse, unparsed, nil
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
	keyMacCommandsToBeSent := fmt.Sprintf(macCommandsKeyTempl, devEUI)
	keyLastMacCommands := fmt.Sprintf(macCommandsWithResponseTempl, devEUI)

	c.Send("GET", keyMacCommandsToBeSent)
	c.Send("GET", keyLastMacCommands)
	c.Flush()
	b, err := redis.Bytes(c.Receive())
	if err != nil && err != redis.ErrNil {
		return macCommands, err
	}
	if err == redis.ErrNil {
		return macCommands, errDoesNotExist
	}

	macCommands.DevEUI = devEUI
	err = gob.NewDecoder(bytes.NewReader(b)).Decode(&macCommands)
	if err != nil {
		return macCommands, err
	}

	// Check that the mac commands are not already answered.
	b, err = redis.Bytes(c.Receive())
	if err == nil {
		var macCommandsWithResponse MacCommandsWithResponse
		if err = gob.NewDecoder(bytes.NewReader(b)).Decode(&macCommandsWithResponse); err != nil {
			return macCommands, err
		}
		emptyResponse := MacCommandsResponse{}
		hasUnansweredMacCommands := (macCommands != macCommandsWithResponse.MacCommands ||
			macCommandsWithResponse.Response == emptyResponse)
		if !hasUnansweredMacCommands {
			return macCommands, errDoesNotExist
		}
	}

	return macCommands, nil
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

// NodeMacCommandsAPI exports the MacCommands related functions.
type NodeMacCommandsAPI struct {
	ctx Context
}

// NewNodeMacCommandsAPI crestes a new NodeMacCommandsAPI.
func NewNodeMacCommandsAPI(ctx Context) *NodeMacCommandsAPI {
	return &NodeMacCommandsAPI{
		ctx: ctx,
	}
}

// GetMacCmdsToBeSent gets the next to be sent and last sent MacCommands.
// The latter has also an response that might or might not be nil.
func (a *NodeMacCommandsAPI) GetMacCmdsToBeSent(devEUI lorawan.EUI64, macCommands *NextAndLastMacCommands) error {
	var err error
	*macCommands, err = getNextAndLastMacCommands(a.ctx.RedisPool, devEUI)
	return err
}

// Update updates the value of the next MacCommands to be sent.
func (a *NodeMacCommandsAPI) Update(data MacCommands, devEUI *lorawan.EUI64) error {
	var err error
	*devEUI, err = updateMacCommandsToBeSent(a.ctx.RedisPool, data)

	return err
}
