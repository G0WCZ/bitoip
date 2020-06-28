/*
Copyright (C) 2020 Graeme Sutherland, Nodestone Limited


This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package core

import (
	"context"
	"fmt"
	"github.com/G0WCZ/cwc/bitoip"
	"github.com/G0WCZ/cwc/config"
	"github.com/G0WCZ/cwc/core/hw"
	"net"
	"strings"
	"time"

	"github.com/golang/glog"
)

const LocalMulticast = "224.0.0.73:%d"

// General  station client
// Can be in local mode, in which case all is local muticast on the local network
// Else the client of a reflector
func StationClient(ctx context.Context, cancel func(), cfg *config.Config) {
	var addr string

	if cfg.NetworkMode == "local" {
		addr = fmt.Sprintf(LocalMulticast, cfg.LocalPort)
		glog.Infof("Starting in local mode with local multicast address %s", addr)
	} else {
		addr = cfg.ReflectorAddress
		glog.Infof("Connecting to reflector %s", addr)
	}

	resolvedAddress, err := net.ResolveUDPAddr("udp", addr)

	if err != nil {
		glog.Errorf("Error resolving address %s %v", addr, err)
		return
	}
	// channel to send to network
	toSend := make(chan bitoip.CarrierEventPayload)

	// channel to send to hardware
	toMorse := make(chan bitoip.RxMSG)

	// channel for configChanges
	configChanges := make(chan config.ConfigChange)

	go General(ctx, cfg)

	go MorseRx(ctx, toSend, cfg)

	go MorseTx(ctx, cfg)

	localRxAddress, err := net.ResolveUDPAddr("udp", "0.0.0.0:0")

	if err != nil {
		glog.Fatalf("Can't allocate local address: %v", err)
	}

	// start the UDP Receiver
	go bitoip.UDPRx(ctx, localRxAddress, toMorse)

	var csBase [16]byte

	// Allow things to catch up. May not be needed anymore
	time.Sleep(time.Second * 1)

	// get callsign into a []byte we can send
	r := strings.NewReader(cfg.Callsign)
	_, err = r.Read(csBase[0:16])

	if err != nil {
		glog.Errorf("Callsign %s can not be encoded", cfg.Callsign)
	}

	var currentCarrierKey bitoip.CarrierKeyType
	var currentChannel bitoip.ChannelIdType = cfg.Channel

	// transmit a listen request to the configured channel
	bitoip.UDPTx(bitoip.ListenRequest, bitoip.ListenRequestPayload{
		cfg.Channel,
		csBase,
	},
		resolvedAddress,
	)

	// Do time sync
	// Set up buckets and fill the buckets with time offset and round-trip data
	const timeOffsetBucketSize = 5

	timeOffsetIndex := 0
	timeOffsets := make([]int64, timeOffsetBucketSize, timeOffsetBucketSize)
	timeOffsetSum := int64(0)
	roundTrips := make([]int64, timeOffsetBucketSize, timeOffsetBucketSize)
	roundTripSum := int64(0)

	commonTimeOffset := int64(0)
	commonRoundTrip := int64(0)

	// set up basis of keepAlive
	lastUDPSend := time.Now()

	keepAliveTick := time.Tick(20 * time.Second)

	// start off with fast time syncs, gets slowed down later
	timeSyncTick := time.Tick(5 * time.Second)

	for i := 0; i < timeOffsetBucketSize; i++ {
		bitoip.UDPTx(bitoip.TimeSync, bitoip.TimeSyncPayload{
			time.Now().UnixNano(),
		}, resolvedAddress)
	}

	timeSyncCount := 0

	// loop on the toSend (from the hardware to send on UDP) and toMorse (send to the morse hardware)
	// channels -- and the keepalive as well.
	// TODO should also redo time sync occasionally as well
	for {
		select {
		case <-ctx.Done():
			SetStatus(hw.StatusLED, hw.Off)
			return

		case cep := <-toSend:
			glog.V(2).Infof("carrier event payload to send: %v", cep)
			// TODO fill in some channel details
			bitoip.UDPTx(bitoip.CarrierEvent, cep, resolvedAddress)

		case tm := <-toMorse:

			// we have data, so turn signal LED on
			SetStatus(hw.StatusLED, hw.On)

			switch tm.Verb {
			case bitoip.CarrierEvent:
				glog.V(2).Infof("carrier events to morse: %v", tm)
				QueueForOutput(tm.Payload.(*bitoip.CarrierEventPayload), cfg)

			case bitoip.ListenConfirm:
				glog.V(2).Infof("listen confirm: %v", tm)
				lc := tm.Payload.(*bitoip.ListenConfirmPayload)
				glog.Infof("listening channel %d with carrier key %d", lc.Channel, lc.CarrierKey)
				currentCarrierKey = lc.CarrierKey
				SetCarrierKey(lc.CarrierKey)

			case bitoip.TimeSyncResponse:
				glog.V(2).Infof("time sync response %v", tm)
				tsr := tm.Payload.(*bitoip.TimeSyncResponsePayload)
				now := time.Now().UnixNano()

				// time offset and roundtrip calculation.  See how NTP does this. Basically
				// the same algorithm
				latestTimeOffset := ((tsr.ServerRxTime - tsr.GivenTime) - (tsr.ServerTxTime - now)) / 2
				roundTrip := (now - tsr.GivenTime) - (tsr.ServerRxTime - tsr.ServerTxTime)

				timeOffsets[timeOffsetIndex] = latestTimeOffset
				timeOffsetIndex = (timeOffsetIndex + 1) % timeOffsetBucketSize

				roundTrips[timeOffsetIndex] = roundTrip

				timeOffsetSum = int64(0)
				roundTripSum = int64(0)

				for i := 0; i < timeOffsetBucketSize; i++ {
					timeOffsetSum += timeOffsets[i]
					roundTripSum += roundTrips[i]
				}
				commonTimeOffset = (timeOffsetSum / timeOffsetBucketSize)
				SetTimeOffset(commonTimeOffset)
				commonRoundTrip = (roundTripSum / timeOffsetBucketSize)
				SetRoundTrip(commonRoundTrip)

				glog.V(2).Infof("timesync: offset %d µs roundtrip %d µs",
					commonTimeOffset/1000,
					commonRoundTrip/1000)
			}

		case kat := <-keepAliveTick:

			// check and send a keepalive if nothing else has happened
			if kat.Sub(lastUDPSend) > time.Duration(20*time.Second) {
				// turn off Status LED - to be turned back on by response above
				SetStatus(hw.StatusLED, hw.Off)

				lastUDPSend = kat
				p := bitoip.ListenRequestPayload{
					cfg.Channel,
					csBase,
				}
				glog.V(2).Info("sending keepalive")
				bitoip.UDPTx(bitoip.ListenRequest, p, resolvedAddress)
			}

		// do time sync
		case tst := <-timeSyncTick:
			timeSyncCount += 1
			if timeSyncCount == 5 {
				// slow down after initial syncs
				timeSyncTick = time.Tick(140 * time.Second)
			}

			// turn off Status LED - to be turned back on by response above
			SetStatus(hw.StatusLED, hw.Off)

			glog.V(2).Info("sending timesync")
			bitoip.UDPTx(bitoip.TimeSync, bitoip.TimeSyncPayload{
				tst.UnixNano(),
			}, resolvedAddress)

		case cc := <-configChanges:
			if cc == config.ConfigChangeRestart {
				cancel()
			} else if cc == config.ConfigChangeChannel {
				glog.V(2).Infof("changing channel from %d to %d", currentChannel, cfg.Channel)

				// unlisten current channel
				bitoip.UDPTx(bitoip.Unlisten, bitoip.UnlistenPayload{
					currentChannel,
					currentCarrierKey,
				},
					resolvedAddress,
				)
			}
			// transmit a listen request to the configured channel
			bitoip.UDPTx(bitoip.ListenRequest, bitoip.ListenRequestPayload{
				cfg.Channel,
				csBase,
			},
				resolvedAddress,
			)

		}
	}
}
