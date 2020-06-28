/*
Copyright (C) 2019 Graeme Sutherland, Nodestone Limited


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
	"github.com/G0WCZ/cwc/bitoip"
	"github.com/G0WCZ/cwc/config"
	"github.com/G0WCZ/cwc/core/hw"
	"github.com/golang/glog"
	"sync"
	"time"
)

import (
	"context"
)

// List of inputs to work with
var inputs []hw.MorseIn

// Time ticker for keyer and straight key
var rxTicker *time.Ticker

// The last bit output
var lastBit bool

var events []hw.TimedBitEvent
var RxMutex = sync.Mutex{}
var MaxEvents = 0

var channelId bitoip.ChannelIdType

func SetChannelId(cId bitoip.ChannelIdType) {
	channelId = cId
}

func ChannelId() bitoip.ChannelIdType {
	return channelId
}

var carrierKey bitoip.CarrierKeyType

func SetCarrierKey(ck bitoip.CarrierKeyType) {
	carrierKey = ck
}

func CarrierKey() bitoip.CarrierKeyType {
	return carrierKey
}

// Clock offset from the reflector
var timeOffset = int64(0)

func SetTimeOffset(t int64) {
	timeOffset = t
}

// round trip delay calculated by the transport
var roundTrip = int64(0)

func SetRoundTrip(t int64) {
	roundTrip = t
}

// This is the morse receiver
// This polls morse input "hardware" for data

func MorseRx(ctx context.Context, morseReceived chan bitoip.CarrierEventPayload, config *config.Config) {
	OpenInputs(config)
	// Setup
	lastBit = false

	// event receive queue
	MaxEvents = config.Advanced.MaxEvents
	events = make([]hw.TimedBitEvent, 0, MaxEvents)

	// Start sampler
	rxTicker = time.NewTicker(time.Duration(config.Advanced.InputTickMs))

	for {
		select {
		case <-ctx.Done():
			rxTicker.Stop()
			CloseInputs()
			return

		case t := <-rxTicker.C:
			Sample(t, morseReceived, config)
		}
	}
}

func OpenInputs(config *config.Config) {
	inputs = hw.ParseInputs(config)
	// assume not active to start with
	lastBit = false

	for _, input := range inputs {
		input.Open()
	}
}

func CloseInputs() {
	for _, input := range inputs {
		input.Close()
	}

	lastBit = false
}

// Dual keyer and straight-key sampler. Does sampling and sends stuff out if ready
// As a 1ms tick, this needs to be efficient
// Either keyer or straight key sampling and post-processing
func Sample(t time.Time, morseRecieved chan bitoip.CarrierEventPayload,
	config *config.Config) {
	// to hold output samples to resolve
	newBit := false

	for idx, input := range inputs {
		// for each input
		if input.UseKeyer() {
			newBit = newBit || hw.KeyerSampleInput(t, idx, input)
		} else {
			// bits are active high
			newBit = newBit || input.Bit()
		}
	}

	if lastBit != newBit {
		var bit uint8 = 0

		// changed output
		if newBit {
			bit = 1
		}

		// store event
		RxMutex.Lock()
		lastBit = newBit
		events = append(events, hw.TimedBitEvent{t, bitoip.BitEvent(bit)})
		RxMutex.Unlock()

		// flush if ready
		if (len(events) >= MaxEvents-1) && (events[len(events)-1].BitEvent&bitoip.BitOn == 0) {
			events = Flush(events, morseRecieved, config)
		}
	}
}

// Flush events and place in the toSend channel to wake up the UDP sender to
// transmit the packet.
func Flush(events []hw.TimedBitEvent, toSend chan bitoip.CarrierEventPayload, config *config.Config) []hw.TimedBitEvent {
	glog.V(2).Infof("Flushing events %v", events)
	RxMutex.Lock()
	if len(events) > 0 {
		toSend <- BuildPayload(events, config)
		events = events[:0]
	}
	RxMutex.Unlock()
	return events
}

// Build a payload (CarrierEventPayload) of on and off events. Called from Flush() to
// make a packet ready to send.
func BuildPayload(events []hw.TimedBitEvent, config *config.Config) bitoip.CarrierEventPayload {
	baseTime := events[0].StartTime.UnixNano()
	// packetStartTime := baseTime + timeOffset + roundTrip/2 + MaxSendTimespan.Nanoseconds()
	packetStartTime := baseTime + timeOffset + roundTrip/2 + (config.Advanced.MaxSendTimespanMs * 1000)
	cep := bitoip.CarrierEventPayload{
		channelId,
		carrierKey,
		packetStartTime,
		[bitoip.MaxBitEvents]bitoip.CarrierBitEvent{},
		time.Now().UnixNano(),
	}
	for i, event := range events {
		bit := event.BitEvent

		// mark last event this message
		if i == (len(events) - 1) {
			bit = bit | bitoip.LastEvent
		}

		cep.BitEvents[i] = bitoip.CarrierBitEvent{
			uint32(event.StartTime.UnixNano() - baseTime),
			bit,
		}
	}
	return cep
}
