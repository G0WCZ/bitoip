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
package main

import (
	"github.com/G0WCZ/cwc/bitoip"
	"net"
	"strings"
	"time"

	"github.com/golang/glog"
)

// Handle an incoming message to the reflector
func Handler(serverAddress *net.UDPAddr, msg bitoip.RxMSG) {
	switch msg.Verb {
	// Channel list
	case bitoip.EnumerateChannels:
		responsePayload := new(bitoip.ListChannelsPayload)
		copy(responsePayload.Channels[:], ChannelIds())
		bitoip.UDPTx(bitoip.ListChannels,
			msg.Payload,
			&msg.SrcAddress)
	// Carrier morse data
	case bitoip.CarrierEvent:
		ce := msg.Payload.(*bitoip.CarrierEventPayload)
		glog.V(1).Infof("got carrier event %v", ce)
		channel := GetChannel(ce.Channel)
		channel.Broadcast(*ce)

	// Subscribe request
	case bitoip.ListenRequest:
		lr := msg.Payload.(*bitoip.ListenRequestPayload)
		channel := GetChannel(lr.Channel)
		key := channel.Subscribe(msg.SrcAddress, strings.Trim(string(lr.Callsign[:]), "\x00"))
		lcp := bitoip.ListenConfirmPayload{lr.Channel, key}

		bitoip.UDPTx(bitoip.ListenConfirm, lcp, &msg.SrcAddress)

	// Time sync
	case bitoip.TimeSync:
		ts := msg.Payload.(*bitoip.TimeSyncPayload)

		tsr := bitoip.TimeSyncResponsePayload{
			ts.CurrentTime,
			msg.RxTime,
			time.Now().UnixNano(),
		}

		bitoip.UDPTx(bitoip.TimeSyncResponse, tsr, &msg.SrcAddress)
	}

}
