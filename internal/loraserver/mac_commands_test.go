package loraserver

import (
	"fmt"
	"testing"

	"github.com/brocaar/lorawan"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMacCommandsPayloads(t *testing.T) {
	Convey("Given a MacCommands", t, func() {
		macCommands := MacCommands{}

		Convey("When it is empty", func() {
			Convey("When converting to slice", func() {
				payloads := macCommands.Payloads()
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
				payloads := macCommands.Payloads()
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
				payloads := macCommands.Payloads()
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
				payloads := macCommands.Payloads()
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
				payloads := macCommands.Payloads()
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
				payloads := macCommands.Payloads()
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
				payloads := macCommands.Payloads()
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

			payloads := macCommands.Payloads()

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
			DutyCycle: lorawan.DutyCycleReqPayload{
				MaxDCCycle: 2,
			},
		}
		macCommandsWithResponse := MacCommandsWithResponse{
			MacCommands: macCommands,
		}
		otherMacCommands := MacCommands{
			DevEUI:    devEUI,
			DevStatus: true,
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

				Convey("Given it has response", func() {
					macCommandsResponse := MacCommandsResponse{
						DutyCycle: true,
					}
					So(updateSentCommandResponse(p, devEUI, macCommandsResponse), ShouldBeNil)

					Convey("Then there should be no mac commands to be sent", func() {
						_, err := getMacCommandsToBeSent(p, devEUI)
						So(err, ShouldEqual, errDoesNotExist)
					})
				})
			})

			Convey("When it is updated", func() {
				updatedDevEUI, err := updateMacCommandsToBeSent(p, otherMacCommands)
				So(err, ShouldBeNil)
				So(updatedDevEUI, ShouldEqual, devEUI)

				Convey("Then the mac commands to be sent is the excpected", func() {
					data, err := getMacCommandsToBeSent(p, devEUI)
					So(err, ShouldBeNil)
					So(data, ShouldResemble, otherMacCommands)
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
				err := saveInRedis(conn, fmt.Sprintf(macCommandsKeyTempl, devEUI), otherMacCommands)
				So(err, ShouldBeNil)

				Convey("Then both can be retrieved with getNextAndLastMacCommands", func() {
					data, err := getNextAndLastMacCommands(p, devEUI)
					So(err, ShouldBeNil)
					expectedData := NextAndLastMacCommands{
						Next: &otherMacCommands,
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

func TestFromMacCommandsSlice(t *testing.T) {
	Convey("Given an empty response slice", t, func() {
		var responseSlice []lorawan.MACCommand
		Convey("When converted to MacCommandsResponse", func() {
			response, unparsedMacCommands, err := parseMacCommandsResponse(responseSlice)
			So(err, ShouldBeNil)

			Convey("Then there should be no unparsed commands", func() {
				So(unparsedMacCommands, ShouldBeNil)
			})

			Convey("Then response has zero-value", func() {
				So(response, ShouldResemble, MacCommandsResponse{})
			})
		})
	})

	Convey("Given a response slice with all the possible mac commands", t, func() {
		devStatus := lorawan.DevStatusAnsPayload{
			Battery: 5,
			Margin:  1,
		}

		linkADR := lorawan.LinkADRAnsPayload{
			ChannelMaskACK: true,
			DataRateACK:    true,
			PowerACK:       true,
		}

		rx2Setup := lorawan.RX2SetupAnsPayload{
			ChannelACK:     true,
			RX1DRoffsetACK: true,
			RX2DataRateACK: true,
		}

		newChannel := lorawan.NewChannelAnsPayload{
			ChannelFrequencyOK: true,
			DataRateRangeOK:    true,
		}

		linkCheckReq := lorawan.MACCommand{
			CID: lorawan.LinkCheckReq,
		}
		responseSlice := []lorawan.MACCommand{
			linkCheckReq,
			lorawan.MACCommand{
				CID:     lorawan.LinkADRAns,
				Payload: &linkADR,
			},
			lorawan.MACCommand{
				CID: lorawan.DutyCycleAns,
			},
			lorawan.MACCommand{
				CID:     lorawan.RXParamSetupAns,
				Payload: &rx2Setup,
			},
			lorawan.MACCommand{
				CID:     lorawan.DevStatusAns,
				Payload: &devStatus,
			},
			lorawan.MACCommand{
				CID:     lorawan.NewChannelAns,
				Payload: &newChannel,
			},
			lorawan.MACCommand{
				CID: lorawan.RXTimingSetupAns,
			},
		}

		Convey("Then all the relevant data is extracted", func() {
			response, unparsedMacCommands, err := parseMacCommandsResponse(responseSlice)
			So(err, ShouldBeNil)

			So(*response.LinkADR, ShouldResemble, linkADR)
			So(response.DutyCycle, ShouldBeTrue)
			So(*response.RX2Setup, ShouldResemble, rx2Setup)
			So(*response.DevStatus, ShouldResemble, devStatus)
			So(*response.NewChannel, ShouldResemble, newChannel)
			So(response.RXTimingSetup, ShouldBeTrue)
			So(unparsedMacCommands, ShouldHaveLength, 1)

			So(unparsedMacCommands[0], ShouldResemble, linkCheckReq)
		})
	})
}
