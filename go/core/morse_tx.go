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
	"context"
	"github.com/G0WCZ/cwc/bitoip"
	"github.com/G0WCZ/cwc/config"
	"github.com/G0WCZ/cwc/core/hw"
	"github.com/golang/glog"
	"sort"
	"sync"
	"time"
)

var ticker *time.Ticker

// Lock for the Morse Transmit Queue
var TxMutex = sync.Mutex{}

// Morse Transmit Queue
var TxQueue []hw.TimedBitEvent

// outputs
var outputs []hw.MorseOut

var bitOut bool

// This is the morse sender
// This sends output bits/streams

func MorseTx(ctx context.Context, config *config.Config) {
	// empty send queue
	TxQueue = []hw.TimedBitEvent{}

	// setup outputs
	OpenOutputs(config)

	//start timer
	ticker = time.NewTicker(time.Duration(config.Advanced.OutputTickMs))

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return

		case t := <-ticker.C:
			TransmitToOutputs(t, outputs)
		}
	}
}

// Resolve config into actual outputs that are enabled
// For example, might be a decoder and a serial bit output
func ConfigureMorseOutputs(config *config.Config) {

}

func OpenOutputs(config *config.Config) {
	bitOut = false

	outputs = hw.ParseOutputs(config)
	for _, output := range outputs {

		output.Open()
	}
}

func CloseOutputs() {
	bitOut = false

	for _, output := range outputs {
		output.Close()
	}
}

// Queue this stuff for sending to hardware -- LED or relay or PWM
// by adding to queue that will be sent out based on the tick timing
func QueueForOutput(carrierEvents *bitoip.CarrierEventPayload, config *config.Config) {
	if (config.RemoteEcho || (carrierEvents.CarrierKey != carrierKey)) &&
		carrierEvents.Channel == channelId {
		// compose into events
		newEvents := make([]hw.TimedBitEvent, 0)

		// remove the calculated server time offset
		start := time.Unix(0, carrierEvents.StartTimeStamp-timeOffset+(roundTrip/2))
		diff := start.UnixNano() - time.Now().UnixNano()
		if diff < 0 {
			// if we have negative time, increase offset a little to 'allow'
			start.Add(time.Duration(diff))
			timeOffset += diff
			glog.V(2).Infof("Negative time offset %v to current time", diff/1000)
		}

		for _, ce := range carrierEvents.BitEvents {
			newEvents = append(newEvents, hw.TimedBitEvent{
				start.Add(time.Duration(ce.TimeOffset)),
				ce.BitEvent,
			})
			if (ce.BitEvent & bitoip.LastEvent) > 0 {
				break
			}
		}

		// Lock and append new events
		TxMutex.Lock()
		TxQueue = append(TxQueue, newEvents...)
		// then sort the output by time (this is probably super slow)
		sort.Slice(TxQueue, func(i, j int) bool { return TxQueue[i].StartTime.Before(TxQueue[j].StartTime) })
		TxMutex.Unlock()
	} else {
		// don't re-sound our own stuff if echo isn't turned on
		glog.V(2).Infof("ignoring own carrier")
	}
	glog.V(2).Infof("TXQueue is now: %v", TxQueue)
}

// When woken up  (same timer as checking for an incoming bit change)
// check to see if an output state change is needed and do it.
func TransmitToOutputs(t time.Time, outputs []hw.MorseOut) {
	now := time.Now()

	// Lock
	TxMutex.Lock()

	// Change outputs if needed
	if len(TxQueue) > 0 && TxQueue[0].StartTime.Before(now) {
		be := TxQueue[0].BitEvent

		newBit := (be & bitoip.BitOn) != 0

		if newBit != bitOut {
			bitOut = newBit
			for _, output := range outputs {
				output.SetBit(newBit)
				output.SetToneOut(newBit)
				output.SetStatusLED(newBit)
			}
		}

		TxQueue = TxQueue[1:]
	}

	// Unlock
	TxMutex.Unlock()
}
