package loraserver

import (
	"errors"
	"testing"
	"time"

	"github.com/brocaar/loraserver/models"
	"github.com/brocaar/lorawan"
	"github.com/brocaar/lorawan/band"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHandleDataUpPackets(t *testing.T) {
	conf := getConfig()

	Convey("Given a clean state", t, func() {
		app := &testApplicationBackend{
			rxPayloadChan: make(chan models.RXPayload, 1),
			txPayloadChan: make(chan models.TXPayload),
		}
		gw := &testGatewayBackend{
			rxPacketChan: make(chan models.RXPacket),
			txPacketChan: make(chan models.TXPacket, 2),
		}
		p := NewRedisPool(conf.RedisURL)
		mustFlushRedis(p)

		ctx := Context{
			RedisPool:   p,
			Gateway:     gw,
			Application: app,
		}

		Convey("Given a stored node session", func() {
			ns := models.NodeSession{
				DevAddr:  [4]byte{1, 2, 3, 4},
				DevEUI:   [8]byte{1, 2, 3, 4, 5, 6, 7, 8},
				AppSKey:  [16]byte{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
				NwkSKey:  [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
				FCntUp:   8,
				FCntDown: 5,
				AppEUI:   [8]byte{8, 7, 6, 5, 4, 3, 2, 1},
			}
			So(saveNodeSession(p, ns), ShouldBeNil)

			Convey("Given an UnconfirmedDataUp packet", func() {
				fPort := uint8(1)
				phy := lorawan.PHYPayload{
					MHDR: lorawan.MHDR{
						MType: lorawan.UnconfirmedDataUp,
						Major: lorawan.LoRaWANR1,
					},
					MACPayload: &lorawan.MACPayload{
						FHDR: lorawan.FHDR{
							DevAddr: ns.DevAddr,
							FCnt:    10,
						},
						FPort: &fPort,
						FRMPayload: []lorawan.Payload{
							&lorawan.DataPayload{Bytes: []byte("hello!")},
						},
					},
				}

				So(phy.EncryptFRMPayload(ns.AppSKey), ShouldBeNil)
				So(phy.SetMIC(ns.NwkSKey), ShouldBeNil)

				rxPacket := models.RXPacket{
					PHYPayload: phy,
				}

				Convey("Given that the application backend returns an error", func() {
					app.err = errors.New("BOOM")

					Convey("Then handleRXPacket returns an error", func() {
						So(handleRXPacket(ctx, rxPacket), ShouldNotBeNil)

						Convey("Then the FCntUp has not been incremented", func() {
							nsUpdated, err := getNodeSession(p, ns.DevAddr)
							So(err, ShouldBeNil)
							So(nsUpdated, ShouldResemble, ns)
						})
					})
				})

				Convey("When the frame-counter is invalid", func() {
					ns.FCntUp = 11
					So(saveNodeSession(p, ns), ShouldBeNil)

					Convey("Then handleRXPacket returns a frame-counter related error", func() {
						err := handleRXPacket(ctx, rxPacket)
						So(err, ShouldResemble, errors.New("invalid FCnt or too many dropped frames"))
					})
				})

				Convey("When the MIC is invalid", func() {
					rxPacket.PHYPayload.MIC = [4]byte{1, 2, 3, 4}

					Convey("Then handleRXPacket returns a MIC related error", func() {
						err := handleRXPacket(ctx, rxPacket)
						So(err, ShouldResemble, errors.New("invalid MIC"))
					})
				})

				Convey("When calling handleRXPacket", func() {
					So(handleRXPacket(ctx, rxPacket), ShouldBeNil)

					Convey("Then the packet is correctly received by the application backend", func() {
						packet := <-app.rxPayloadChan
						So(packet, ShouldResemble, models.RXPayload{
							DevEUI:       ns.DevEUI,
							FPort:        1,
							GatewayCount: 1,
							Data:         []byte("hello!"),
						})
					})

					Convey("Then the FCntUp must be synced", func() {
						nsUpdated, err := getNodeSession(p, ns.DevAddr)
						So(err, ShouldBeNil)
						So(nsUpdated.FCntUp, ShouldEqual, 10)
					})

					Convey("Then no downlink data was sent to the gateway", func() {
						var received bool
						select {
						case <-gw.txPacketChan:
							received = true
						case <-time.After(time.Second):
							// nothing to do
						}
						So(received, ShouldEqual, false)
					})
				})

				Convey("Given an enqueued TXPayload (unconfirmed)", func() {
					txPayload := models.TXPayload{
						Confirmed: false,
						DevEUI:    ns.DevEUI,
						FPort:     5,
						Data:      []byte("hello back!"),
					}
					So(addTXPayloadToQueue(p, txPayload), ShouldBeNil)

					Convey("When calling handleRXPacket", func() {
						So(handleRXPacket(ctx, rxPacket), ShouldBeNil)

						Convey("Then the packet is received by the application backend", func() {
							_ = <-app.rxPayloadChan
						})

						Convey("Then two identical packets are sent to the gateway (two receive windows)", func() {
							txPacket1 := <-gw.txPacketChan
							txPacket2 := <-gw.txPacketChan
							So(txPacket1.PHYPayload, ShouldResemble, txPacket2.PHYPayload)

							macPL, ok := txPacket1.PHYPayload.MACPayload.(*lorawan.MACPayload)
							So(ok, ShouldBeTrue)

							Convey("Then these packets contain the expected values", func() {
								So(txPacket1.PHYPayload.MHDR.MType, ShouldEqual, lorawan.UnconfirmedDataDown)
								So(macPL.FHDR.FCnt, ShouldEqual, ns.FCntDown)
								So(macPL.FHDR.FCtrl.ACK, ShouldBeFalse)
								So(*macPL.FPort, ShouldEqual, 5)

								So(txPacket1.PHYPayload.DecryptFRMPayload(ns.AppSKey), ShouldBeNil)
								So(macPL.FRMPayload, ShouldHaveLength, 1)
								pl, ok := macPL.FRMPayload[0].(*lorawan.DataPayload)
								So(ok, ShouldBeTrue)
								So(pl.Bytes, ShouldResemble, []byte("hello back!"))
							})
						})

						Convey("Then the FCntDown was incremented", func() {
							ns2, err := getNodeSession(p, ns.DevAddr)
							So(err, ShouldBeNil)
							So(ns2.FCntDown, ShouldEqual, ns.FCntDown+1)
						})

						Convey("Then the TXPayload queue is empty", func() {
							_, _, err := getTXPayloadAndRemainingFromQueue(p, ns.DevEUI)
							So(err, ShouldEqual, errDoesNotExist)
						})
					})

					Convey("Given a node specific TX params", func() {
						delay := uint8(3)
						tx2Frequency := uint32(868400000)
						ns.TXParams = models.TXParams{
							Set:          true,
							TXDelay1:     delay,
							TX1DRoffset:  1,
							TX2DataRate:  3,
							TX2Frequency: tx2Frequency,
						}
						So(saveNodeSession(p, ns), ShouldBeNil)

						rxPacket.RXInfo.DataRate = band.DataRateConfiguration[2]
						Convey("When calling handleRXPacket", func() {
							So(handleRXPacket(ctx, rxPacket), ShouldBeNil)

							Convey("Then two tx packets are sent to the gateway with correct params", func() {
								txPacket1 := <-gw.txPacketChan
								txPacket2 := <-gw.txPacketChan

								delayDuration1 := uint32(time.Duration(delay) * time.Second / time.Microsecond)
								So(txPacket1.TXInfo.DataRate, ShouldResemble, band.DataRateConfiguration[1])
								So(txPacket1.TXInfo.Frequency, ShouldEqual, rxPacket.RXInfo.Frequency)
								So(txPacket1.TXInfo.Timestamp, ShouldEqual, rxPacket.RXInfo.Timestamp+delayDuration1)

								delayDuration2 := uint32(time.Duration(delay+1) * time.Second / time.Microsecond)
								So(txPacket2.TXInfo.DataRate, ShouldResemble, band.DataRateConfiguration[3])
								So(txPacket2.TXInfo.Frequency, ShouldEqual, tx2Frequency)
								So(txPacket2.TXInfo.Timestamp, ShouldEqual, rxPacket.RXInfo.Timestamp+delayDuration2)
							})
						})
					})
				})

				Convey("Given an enqueued TXPayload (confirmed)", func() {
					txPayload := models.TXPayload{
						Confirmed: true,
						DevEUI:    ns.DevEUI,
						FPort:     5,
						Data:      []byte("hello back!"),
					}
					So(addTXPayloadToQueue(p, txPayload), ShouldBeNil)

					Convey("When calling handleRXPacket", func() {
						So(handleRXPacket(ctx, rxPacket), ShouldBeNil)

						Convey("Then the packet is received by the application backend", func() {
							_ = <-app.rxPayloadChan
						})

						Convey("Then two identical packets are sent to the gateway (two receive windows)", func() {
							txPacket1 := <-gw.txPacketChan
							txPacket2 := <-gw.txPacketChan
							So(txPacket1.PHYPayload, ShouldResemble, txPacket2.PHYPayload)

							macPL, ok := txPacket1.PHYPayload.MACPayload.(*lorawan.MACPayload)
							So(ok, ShouldBeTrue)

							Convey("Then these packets contain the expected values", func() {
								So(txPacket1.PHYPayload.MHDR.MType, ShouldEqual, lorawan.ConfirmedDataDown)
								So(macPL.FHDR.FCnt, ShouldEqual, ns.FCntDown)
								So(macPL.FHDR.FCtrl.ACK, ShouldBeFalse)
								So(*macPL.FPort, ShouldEqual, 5)

								So(txPacket1.PHYPayload.DecryptFRMPayload(ns.AppSKey), ShouldBeNil)
								So(len(macPL.FRMPayload), ShouldEqual, 1)
								pl, ok := macPL.FRMPayload[0].(*lorawan.DataPayload)
								So(ok, ShouldBeTrue)
								So(pl.Bytes, ShouldResemble, []byte("hello back!"))
							})
						})

						Convey("Then the TXPayload is still in the queue", func() {
							tx, _, err := getTXPayloadAndRemainingFromQueue(p, ns.DevEUI)
							So(err, ShouldBeNil)
							So(tx, ShouldResemble, txPayload)
						})

						Convey("Then the FCntDown was not incremented", func() {
							ns2, err := getNodeSession(p, ns.DevAddr)
							So(err, ShouldBeNil)
							So(ns2.FCntDown, ShouldEqual, ns.FCntDown)

							Convey("Given the node sends an ACK", func() {
								phy := lorawan.PHYPayload{
									MHDR: lorawan.MHDR{
										MType: lorawan.UnconfirmedDataUp,
										Major: lorawan.LoRaWANR1,
									},
									MACPayload: &lorawan.MACPayload{
										FHDR: lorawan.FHDR{
											DevAddr: ns.DevAddr,
											FCnt:    11,
											FCtrl: lorawan.FCtrl{
												ACK: true,
											},
										},
									},
								}
								So(phy.SetMIC(ns.NwkSKey), ShouldBeNil)

								rxPacket := models.RXPacket{
									PHYPayload: phy,
								}

								So(handleRXPacket(ctx, rxPacket), ShouldBeNil)

								Convey("Then the FCntDown was incremented", func() {
									ns2, err := getNodeSession(p, ns.DevAddr)
									So(err, ShouldBeNil)
									So(ns2.FCntDown, ShouldEqual, ns.FCntDown+1)

									Convey("Then the TXPayload queue is empty", func() {
										_, _, err := getTXPayloadAndRemainingFromQueue(p, ns.DevEUI)
										So(err, ShouldEqual, errDoesNotExist)
									})
								})
							})
						})
					})
				})

				var delay uint8 = 4
				frequency := uint32(86850000)
				rx2DataRate := uint8(2)
				rx1DRoffset := uint8(1)
				macCommands := MacCommands{
					DevEUI:    ns.DevEUI,
					DevStatus: true,
					RXTimingSetup: lorawan.RXTimingSetupReqPayload{
						Delay: delay,
					},
					RX2Setup: lorawan.RX2SetupReqPayload{
						Frequency: frequency / 100,
						DLsettings: lorawan.DLsettings{
							RX1DRoffset: rx1DRoffset,
							RX2DataRate: rx2DataRate,
						},
					},
				}
				Convey("Given pending MacCommands", func() {
					_, err := updateMacCommandsToBeSent(p, macCommands)
					So(err, ShouldBeNil)

					Convey("When calling handleRXPacket", func() {
						So(handleRXPacket(ctx, rxPacket), ShouldBeNil)

						Convey("Then the packet is received by the application backend", func() {
							_ = <-app.rxPayloadChan
						})

						Convey("Then two identical packets are sent to the gateway (two receive windows)", func() {
							So(gw.txPacketChan, ShouldHaveLength, 2)

							txPacket1 := <-gw.txPacketChan
							txPacket2 := <-gw.txPacketChan
							So(txPacket1.PHYPayload, ShouldResemble, txPacket2.PHYPayload)

							Convey("Then these packets contain the expected values", func() {
								So(txPacket1.PHYPayload.MHDR.MType, ShouldEqual, lorawan.ConfirmedDataDown)
								So(txPacket1.PHYPayload.DecryptFRMPayload(ns.NwkSKey), ShouldBeNil)

								macPL, ok := txPacket1.PHYPayload.MACPayload.(*lorawan.MACPayload)
								So(ok, ShouldBeTrue)

								So(macPL.FHDR.FCnt, ShouldEqual, ns.FCntDown)
								So(*macPL.FPort, ShouldEqual, 0)
								So(macPL.FRMPayload, ShouldHaveLength, 3)

								for _, macCommandRaw := range macPL.FRMPayload {
									macCommand, ok := macCommandRaw.(*lorawan.MACCommand)
									So(ok, ShouldBeTrue)
									if macCommand.CID == lorawan.RXTimingSetupReq {
										So(macCommand.Payload, ShouldResemble, &lorawan.RXTimingSetupReqPayload{
											Delay: delay,
										})
									}
									if macCommand.CID == lorawan.RXParamSetupReq {
										expectedRXParamSetup := lorawan.RX2SetupReqPayload{
											Frequency: frequency / 100,
											DLsettings: lorawan.DLsettings{
												RX1DRoffset: rx1DRoffset,
												RX2DataRate: rx2DataRate,
											},
										}
										So(macCommand.Payload, ShouldResemble, &expectedRXParamSetup)
									}
								}
							})

							Convey("Then the mac commands are marked as sent", func() {
								data, err := getLastMacCommandsWithResponse(p, ns.DevEUI)
								So(err, ShouldBeNil)
								expectedMacCommandsWithResponse := MacCommandsWithResponse{
									MacCommands: macCommands,
								}
								So(*data, ShouldResemble, expectedMacCommandsWithResponse)
							})
						})
					})
				})

				Convey("Given sent MacCommands", func() {
					macCommandsWithResponse := &MacCommandsWithResponse{
						MacCommands: macCommands,
					}
					_, err := updateMacCommandsToBeSent(p, macCommands)
					So(err, ShouldBeNil)
					So(saveInRedisMacCommandsWithResponce(p, macCommandsWithResponse), ShouldBeNil)

					Convey("Given a packet with mac commands responses", func() {
						devStatus := lorawan.DevStatusAnsPayload{
							Battery: 5,
							Margin:  1,
						}
						rx2Setup := lorawan.RX2SetupAnsPayload{
							RX1DRoffsetACK: true,
							RX2DataRateACK: true,
							ChannelACK:     true,
						}

						responses := []lorawan.MACCommand{
							lorawan.MACCommand{
								CID:     lorawan.DevStatusAns,
								Payload: &devStatus,
							},
							lorawan.MACCommand{
								CID: lorawan.RXTimingSetupAns,
							},
							lorawan.MACCommand{
								CID:     lorawan.RXParamSetupAns,
								Payload: &rx2Setup,
							},
						}

						macCommandsResponsesHandling := func(rxPacket models.RXPacket) func() {
							return func() {
								So(handleRXPacket(ctx, rxPacket), ShouldBeNil)

								Convey("Then the mac commands replies are updated", func() {
									data, err := getLastMacCommandsWithResponse(p, ns.DevEUI)
									So(err, ShouldBeNil)
									expectedMacCommandsWithResponse := MacCommandsWithResponse{
										MacCommands: macCommands,
										Response: MacCommandsResponse{
											DevStatus:     &devStatus,
											RXTimingSetup: true,
											RX2Setup:      &rx2Setup,
										},
									}
									So(*data, ShouldResemble, expectedMacCommandsWithResponse)
								})

								Convey("Then node session is updated", func() {
									newNS, err := getNodeSession(ctx.RedisPool, ns.DevAddr)
									So(err, ShouldBeNil)
									So(newNS.TXParams.TXDelay1, ShouldEqual, delay)

									So(newNS.TXParams.TX1DRoffset, ShouldEqual, rx1DRoffset)
									So(newNS.TXParams.TX2Frequency, ShouldEqual, frequency)
									So(newNS.TXParams.TX2DataRate, ShouldEqual, rx2DataRate)
								})
							}
						}

						macPL, ok := phy.MACPayload.(*lorawan.MACPayload)
						So(ok, ShouldBeTrue)

						Convey("Given the mac command responses are in FOpts", func() {
							macPL.FHDR.FOpts = responses
							rxPacket.PHYPayload.MACPayload = macPL
							So(rxPacket.PHYPayload.SetMIC(ns.NwkSKey), ShouldBeNil)

							Convey("When calling handleRXPacket", macCommandsResponsesHandling(rxPacket))
						})

						Convey("Given the mac command responses are in FRMPayload", func() {
							fPort := uint8(0)
							macPL.FPort = &fPort
							macPL.FRMPayload = nil
							for _, data := range responses {
								macCommand := data
								macPL.FRMPayload = append(macPL.FRMPayload, &macCommand)
							}
							rxPacket.PHYPayload.MACPayload = macPL
							So(rxPacket.PHYPayload.EncryptFRMPayload(ns.NwkSKey), ShouldBeNil)
							So(rxPacket.PHYPayload.SetMIC(ns.NwkSKey), ShouldBeNil)

							Convey("When calling handleRXPacket", macCommandsResponsesHandling(rxPacket))
						})
					})
				})
			})

			Convey("Given LinkcCheckReq received", func() {
				macPL := lorawan.MACPayload{
					FHDR: lorawan.FHDR{
						DevAddr: ns.DevAddr,
						FCnt:    10,
					},
				}
				macPL.FRMPayload = nil
				macPL.FHDR.FOpts = []lorawan.MACCommand{
					lorawan.MACCommand{
						CID: lorawan.LinkCheckReq,
					},
				}
				phy := lorawan.PHYPayload{
					MACPayload: &macPL,
					MHDR: lorawan.MHDR{
						MType: lorawan.ConfirmedDataUp,
						Major: lorawan.LoRaWANR1,
					},
				}
				So(phy.SetMIC(ns.NwkSKey), ShouldBeNil)
				margin := 3.0
				rxPacket := models.RXPacket{
					PHYPayload: phy,
					RXInfo: models.RXInfo{
						LoRaSNR: margin,
					},
				}

				Convey("When calling handleRXPacket", func() {
					So(handleRXPacket(ctx, rxPacket), ShouldBeNil)

					Convey("Then LinkCheckAns should be sent", func() {
						So(gw.txPacketChan, ShouldHaveLength, 2)
						txPacket1 := <-gw.txPacketChan
						txPacket2 := <-gw.txPacketChan
						So(txPacket1.PHYPayload, ShouldResemble, txPacket2.PHYPayload)

						So(txPacket1.PHYPayload.DecryptFRMPayload(ns.NwkSKey), ShouldBeNil)
						macPL, ok := txPacket1.PHYPayload.MACPayload.(*lorawan.MACPayload)
						So(*macPL.FPort, ShouldEqual, 0)
						So(ok, ShouldBeTrue)
						Convey("Then it should have the expected value", func() {

							frmPayload := macPL.FRMPayload
							So(frmPayload, ShouldHaveLength, 1)
							linkCheckAns, ok := frmPayload[0].(*lorawan.MACCommand)
							So(ok, ShouldBeTrue)

							payload, ok := linkCheckAns.Payload.(*lorawan.LinkCheckAnsPayload)
							So(ok, ShouldBeTrue)
							So(payload.GwCnt, ShouldEqual, 1)
							So(payload.Margin, ShouldEqual, uint8(margin))
						})
					})
				})
			})

			Convey("Given a ConfirmedDataUp packet", func() {
				fPort := uint8(1)
				phy := lorawan.PHYPayload{
					MHDR: lorawan.MHDR{
						MType: lorawan.ConfirmedDataUp,
						Major: lorawan.LoRaWANR1,
					},
					MACPayload: &lorawan.MACPayload{
						FHDR: lorawan.FHDR{
							DevAddr: ns.DevAddr,
							FCnt:    10,
						},
						FPort: &fPort,
						FRMPayload: []lorawan.Payload{
							&lorawan.DataPayload{Bytes: []byte("hello!")},
						},
					},
				}

				So(phy.EncryptFRMPayload(ns.AppSKey), ShouldBeNil)
				So(phy.SetMIC(ns.NwkSKey), ShouldBeNil)

				rxPacket := models.RXPacket{
					PHYPayload: phy,
				}

				Convey("When calling handleRXPacket", func() {
					So(handleRXPacket(ctx, rxPacket), ShouldBeNil)

					Convey("Then the packet is correctly received by the application backend", func() {
						packet := <-app.rxPayloadChan
						So(packet, ShouldResemble, models.RXPayload{
							DevEUI:       ns.DevEUI,
							FPort:        1,
							GatewayCount: 1,
							Data:         []byte("hello!"),
						})
					})

					Convey("Then two ACK packets are sent to the gateway (two receive windows)", func() {
						txPacket1 := <-gw.txPacketChan
						txPacket2 := <-gw.txPacketChan

						So(txPacket1.PHYPayload.MIC, ShouldEqual, txPacket2.PHYPayload.MIC)

						macPL, ok := txPacket1.PHYPayload.MACPayload.(*lorawan.MACPayload)
						So(ok, ShouldBeTrue)

						So(macPL.FHDR.FCtrl.ACK, ShouldBeTrue)
						So(macPL.FHDR.FCnt, ShouldEqual, 5)

						Convey("Then the FCntDown counter has incremented", func() {
							ns, err := getNodeSession(ctx.RedisPool, ns.DevAddr)
							So(err, ShouldBeNil)
							So(ns.FCntDown, ShouldEqual, 6)
						})
					})
				})

			})
		})
	})
}

func TestHandleJoinRequestPackets(t *testing.T) {
	conf := getConfig()

	Convey("Given a dummy gateway and application backend and a clean Postgres and Redis database", t, func() {
		a := &testApplicationBackend{
			rxPayloadChan: make(chan models.RXPayload, 1),
		}
		g := &testGatewayBackend{
			rxPacketChan: make(chan models.RXPacket),
			txPacketChan: make(chan models.TXPacket, 2),
		}
		p := NewRedisPool(conf.RedisURL)
		mustFlushRedis(p)
		db, err := OpenDatabase(conf.PostgresDSN)
		So(err, ShouldBeNil)
		mustResetDB(db)

		ctx := Context{
			RedisPool:   p,
			Gateway:     g,
			Application: a,
			DB:          db,
		}

		Convey("Given a node and application in the database", func() {
			app := models.Application{
				AppEUI: [8]byte{1, 2, 3, 4, 5, 6, 7, 8},
				Name:   "test app",
			}
			So(createApplication(ctx.DB, app), ShouldBeNil)

			node := models.Node{
				DevEUI: [8]byte{8, 7, 6, 5, 4, 3, 2, 1},
				AppEUI: app.AppEUI,
				AppKey: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			}
			So(createNode(ctx.DB, node), ShouldBeNil)

			Convey("Given a JoinRequest packet", func() {
				phy := lorawan.PHYPayload{
					MHDR: lorawan.MHDR{
						MType: lorawan.JoinRequest,
						Major: lorawan.LoRaWANR1,
					},
					MACPayload: &lorawan.JoinRequestPayload{
						AppEUI:   app.AppEUI,
						DevEUI:   node.DevEUI,
						DevNonce: [2]byte{1, 2},
					},
				}
				So(phy.SetMIC(node.AppKey), ShouldBeNil)

				rxPacket := models.RXPacket{
					PHYPayload: phy,
				}

				Convey("When calling handleRXPacket", func() {
					So(handleRXPacket(ctx, rxPacket), ShouldBeNil)

					Convey("Then a JoinAccept was sent to the node", func() {
						txPacket := <-g.txPacketChan
						phy := txPacket.PHYPayload
						So(phy.DecryptJoinAcceptPayload(node.AppKey), ShouldBeNil)
						So(phy.MHDR.MType, ShouldEqual, lorawan.JoinAccept)

						Convey("Then the first delay is 5 sec", func() {
							So(txPacket.TXInfo.Timestamp, ShouldEqual, rxPacket.RXInfo.Timestamp+uint32(5*time.Second/time.Microsecond))
						})

						Convey("Then the second delay is 6 sec", func() {
							txPacket = <-g.txPacketChan
							So(txPacket.TXInfo.Timestamp, ShouldEqual, rxPacket.RXInfo.Timestamp+uint32(6*time.Second/time.Microsecond))
						})

						Convey("Then a node-session was created", func() {
							jaPL, ok := phy.MACPayload.(*lorawan.JoinAcceptPayload)
							So(ok, ShouldBeTrue)

							_, err := getNodeSession(ctx.RedisPool, jaPL.DevAddr)
							So(err, ShouldBeNil)
						})

						Convey("Then the dev-nonce was added to the used dev-nonces", func() {
							node, err := getNode(ctx.DB, node.DevEUI)
							So(err, ShouldBeNil)
							So([2]byte{1, 2}, ShouldBeIn, node.UsedDevNonces)
						})
					})
				})
			})
		})
	})
}

func TestDataRateOffset(t *testing.T) {
	Convey("Given valid data rates and offsets", t, func() {
		examples := []struct {
			dataRate band.DataRate
			offset   uint8
			expected band.DataRate
		}{
			{band.DataRateConfiguration[0], 0, band.DataRateConfiguration[0]},
			{band.DataRateConfiguration[2], 2, band.DataRateConfiguration[0]},
			{band.DataRateConfiguration[4], 2, band.DataRateConfiguration[2]},
			{band.DataRateConfiguration[5], 4, band.DataRateConfiguration[1]},
		}

		Convey("Then DataRate.Offset provides correct results", func() {
			for _, example := range examples {
				actual, err := dataRateOffset(example.dataRate, example.offset)
				So(err, ShouldBeNil)
				So(actual, ShouldResemble, example.expected)
			}
		})
	})

	Convey("Given invalid data rates or offsets", t, func() {
		examples := []struct {
			dataRate band.DataRate
			offset   uint8
		}{
			{band.DataRate{}, 0},
			{band.DataRateConfiguration[0], 1},
		}

		Convey("Then DataRate.Offset returns an error", func() {
			for _, example := range examples {
				_, err := dataRateOffset(example.dataRate, example.offset)
				So(err, ShouldNotBeNil)
			}
		})
	})
}
