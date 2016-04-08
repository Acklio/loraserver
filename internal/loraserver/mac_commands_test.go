package loraserver

import (
	"fmt"
	"testing"

	"github.com/brocaar/lorawan"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMacCommandsToSlice(t *testing.T) {
	Convey("Given a MacCommands", t, func() {
		macCommands := MacCommands{}

		Convey("When it is empty", func() {
			Convey("When converting to slice", func() {
				payloads := macCommands.MacCommandSlice()
				So(payloads, ShouldBeNil)
			})
		})

		linkADRPayload := lorawan.LinkADRReqPayload{
			DataRate: 1,
			TXPower:  14,
			ChMask:   lorawan.ChMask{true, false, true},
			Redundancy: lorawan.Redundancy{
				ChMaskCntl: 1,
				NbRep:      3,
			},
		}

		Convey("When it has LinkADR", func() {
			macCommands.LinkADR = linkADRPayload
			Convey("When converting to slice", func() {
				payloads := macCommands.MacCommandSlice()
				So(payloads, ShouldHaveLength, 1)
				expectedPayload := &lorawan.MACCommand{
					Payload: &linkADRPayload,
					CID:     lorawan.LinkADRReq,
				}
				So(payloads[0], ShouldResemble, expectedPayload)
			})
		})

		dutyCyclePayload := lorawan.DutyCycleReqPayload{
			MaxDCCycle: 3,
		}

		Convey("When it has DutyCycle", func() {
			macCommands.DutyCycle = dutyCyclePayload
			Convey("When converting to slice", func() {
				payloads := macCommands.MacCommandSlice()
				So(payloads, ShouldHaveLength, 1)
				expectedPayload := &lorawan.MACCommand{
					Payload: &dutyCyclePayload,
					CID:     lorawan.DutyCycleReq,
				}
				So(payloads[0], ShouldResemble, expectedPayload)
			})
		})

		rx2SetupPayload := lorawan.RX2SetupReqPayload{
			Frequency: 3,
			DLsettings: lorawan.DLsettings{
				RX2DataRate: 2,
				RX1DRoffset: 3,
			},
		}

		Convey("When it has RX2Setup", func() {
			macCommands.RX2Setup = rx2SetupPayload
			Convey("When converting to slice", func() {
				payloads := macCommands.MacCommandSlice()
				So(payloads, ShouldHaveLength, 1)
				expectedPayload := &lorawan.MACCommand{
					Payload: &rx2SetupPayload,
					CID:     lorawan.RXParamSetupReq,
				}
				So(payloads[0], ShouldResemble, expectedPayload)
			})
		})

		devStatus := true

		Convey("When it has DevStatus", func() {
			macCommands.DevStatus = devStatus
			Convey("When converting to slice", func() {
				payloads := macCommands.MacCommandSlice()
				So(payloads, ShouldHaveLength, 1)
				expectedPayload := &lorawan.MACCommand{
					CID: lorawan.DevStatusReq,
				}
				So(payloads[0], ShouldResemble, expectedPayload)
			})
		})

		newChannelPayload := lorawan.NewChannelReqPayload{
			ChIndex: 15,
			Freq:    868500000,
			MaxDR:   3,
			MinDR:   1,
		}

		Convey("When it has NewChannel", func() {
			macCommands.NewChannel = newChannelPayload
			Convey("When converting to slice", func() {
				payloads := macCommands.MacCommandSlice()
				So(payloads, ShouldHaveLength, 1)
				expectedPayload := &lorawan.MACCommand{
					Payload: &newChannelPayload,
					CID:     lorawan.NewChannelReq,
				}
				So(payloads[0], ShouldResemble, expectedPayload)
			})
		})

		rxTimingSetupPayload := lorawan.RXTimingSetupReqPayload{
			Delay: 2,
		}

		Convey("When it has RXTimingSetup", func() {
			macCommands.RXTimingSetup = rxTimingSetupPayload
			Convey("When converting to slice", func() {
				payloads := macCommands.MacCommandSlice()
				So(payloads, ShouldHaveLength, 1)
				expectedPayload := &lorawan.MACCommand{
					Payload: &rxTimingSetupPayload,
					CID:     lorawan.RXTimingSetupReq,
				}
				So(payloads[0], ShouldResemble, expectedPayload)
			})
		})

		Convey("When converting to slice, if all possible MAC commands are present", func() {
			macCommands.LinkADR = linkADRPayload
			macCommands.DutyCycle = dutyCyclePayload
			macCommands.RX2Setup = rx2SetupPayload
			macCommands.DevStatus = devStatus
			macCommands.NewChannel = newChannelPayload
			macCommands.RXTimingSetup = rxTimingSetupPayload

			payloads := macCommands.MacCommandSlice()

			Convey("Then it has the appropriate length", func() {
				So(payloads, ShouldHaveLength, 6)
			})
		})
	})
}

func TestMacCommandsSaving(t *testing.T) {
	conf := getConfig()

	Convey("Given an empty database", t, func() {
		p := NewRedisPool(conf.RedisURL)
		mustFlushRedis(p)
		devEUI := lorawan.EUI64{1, 2}

		Convey("Then there are no saved mac commands to be sent", func() {
			data, err := getMacCommandsToBeSent(p, devEUI)
			So(err, ShouldEqual, errDoesNotExist)
			So(data.DevEUI, ShouldEqual, lorawan.EUI64{})
		})

		macCommands := MacCommands{
			DevEUI: devEUI,
		}
		macCommandsWithResponse := MacCommandsWithResponse{
			MacCommands: macCommands,
		}

		conn := p.Get()
		defer conn.Close()

		Convey("Given there is unsent MacCommmands in the database", func() {
			err := saveInRedis(conn, fmt.Sprintf(macCommandsKeyTempl, devEUI), macCommands)
			So(err, ShouldBeNil)

			Convey("Then it is possible to be retrieved", func() {
				data, err := getMacCommandsToBeSent(p, devEUI)
				So(err, ShouldBeNil)
				So(data, ShouldResemble, macCommands)
			})

			Convey("Given it can be marked as sent", func() {
				err := markMacCommandsAsSent(p, devEUI, macCommands)
				So(err, ShouldBeNil)

				Convey("Then it can be retrieved", func() {
					data, err := getLastMacCommandsWithResponse(p, devEUI)
					So(err, ShouldBeNil)
					So(*data, ShouldResemble, macCommandsWithResponse)
				})
			})
		})

		Convey("Then there are no last mac commands", func() {
			data, err := getLastMacCommandsWithResponse(p, devEUI)
			So(err, ShouldBeNil)
			So(data, ShouldEqual, nil)
		})

		Convey("Given there is a mac command with response", func() {
			err := saveInRedis(conn, fmt.Sprintf(macCommandsWithResponseTempl, devEUI), macCommandsWithResponse)
			So(err, ShouldBeNil)

			Convey("Then it is possible to be retrieved", func() {
				data, err := getLastMacCommandsWithResponse(p, devEUI)
				So(err, ShouldBeNil)
				So(*data, ShouldResemble, macCommandsWithResponse)
			})

			Convey("Given there is another mac command pending", func() {
				newMacCommands := MacCommands{
					DevEUI:    devEUI,
					DevStatus: true,
				}
				err := saveInRedis(conn, fmt.Sprintf(macCommandsKeyTempl, devEUI), newMacCommands)
				So(err, ShouldBeNil)

				Convey("Then both can be retrieved with getNextAndLastMacCommands", func() {
					data, err := getNextAndLastMacCommands(p, devEUI)
					So(err, ShouldBeNil)
					expectedData := NextAndLastMacCommands{
						Next: &newMacCommands,
						Last: &macCommandsWithResponse,
					}
					So(data, ShouldResemble, expectedData)
				})
			})

			Convey("When we update the response", func() {
				macCommandsResponse := MacCommandsResponse{
					DutyCycle: true,
				}
				err := updateSentCommandResponse(p, devEUI, macCommandsResponse)
				So(err, ShouldBeNil)
				Convey("Then the response the the last mac commands is the expected", func() {
					data, err := getLastMacCommandsWithResponse(p, devEUI)
					So(err, ShouldBeNil)
					So(data.Response, ShouldResemble, macCommandsResponse)
				})
			})
		})
	})
}
