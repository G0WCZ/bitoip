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

package bitoip

import (
	"bytes"
	"encoding/binary"
	"log"
	"reflect"

	"github.com/golang/glog"
)

// conservative UDP payload in bytes
const MaxMessageSizeInBytes = 200
const CallsignSize = 16

var byteOrder = binary.BigEndian

type (
	MessageVerb    = byte
	ChannelIdType  = uint16
	CarrierKeyType = uint16
	Callsign       = [CallsignSize]byte
	Payload        = interface{}
)

const (
	EnumerateChannels MessageVerb = 0x90
	ListChannels      MessageVerb = 0x91
	TimeSync          MessageVerb = 0x92
	TimeSyncResponse  MessageVerb = 0x93
	ListenRequest     MessageVerb = 0x94
	ListenConfirm     MessageVerb = 0x95
	Unlisten          MessageVerb = 0x96
	KeyValue          MessageVerb = 0x81
	CarrierEvent      MessageVerb = 0x82
	VersionInfo       MessageVerb = 0x97
)

// List of Channels
const MaxChannelsPerMessage int = (MaxMessageSizeInBytes - 1) / 2

type ListChannelsPayload struct {
	Channels [MaxChannelsPerMessage]uint16
}

// TimeSync
type TimeSyncPayload struct {
	CurrentTime int64
}

type TimeSyncResponsePayload struct {
	GivenTime    int64
	ServerRxTime int64
	ServerTxTime int64
}

type ListenRequestPayload struct {
	Channel  ChannelIdType
	Callsign [16]byte
}

type ListenConfirmPayload struct {
	Channel    ChannelIdType
	CarrierKey CarrierKeyType
}

type UnlistenPayload struct {
	Channel    ChannelIdType
	CarrierKey CarrierKeyType
}

type KeyValuePayload struct {
	Channel    ChannelIdType
	CarrierKey CarrierKeyType
	Key        [8]byte
	Value      [16]byte
}

type BitEvent uint8

const (
	// Mostly for local paddle before keyer sampling
	BitRightOn  BitEvent = 0x10 // Right, Ring
	BitRightOff BitEvent = 0x08 // Right, Ring
	BitLeftOn   BitEvent = 0x04 // Left, Tip
	BitLeftOff  BitEvent = 0x02 // Left, Tip
	// Normal keying
	BitOn     BitEvent = 0x01 // Straight On
	BitOff    BitEvent = 0x00 // Straight Off
	LastEvent BitEvent = 0x80 // high bit set to indicate last one
)

const MaxBitEvents = (MaxMessageSizeInBytes - 22) / 5
const MaxNsPerCarrierEvent = 2 ^ 32

// Offset allows for about 4 seconds of offset
type CarrierBitEvent struct {
	TimeOffset uint32
	BitEvent   BitEvent
}

type CarrierEventPayload struct {
	Channel        ChannelIdType
	CarrierKey     CarrierKeyType
	StartTimeStamp int64
	BitEvents      [MaxBitEvents]CarrierBitEvent
	SendTime       int64
}

type VersionInfoPayload struct {
	MyProtocolVersion   Version
	MyCodeVersion       Version
	LatestStableStation Version
}

var messagePayload = map[MessageVerb]reflect.Type{
	EnumerateChannels: nil,
	ListChannels:      reflect.TypeOf(ListChannelsPayload{}),
	TimeSync:          reflect.TypeOf(TimeSyncPayload{}),
	TimeSyncResponse:  reflect.TypeOf(TimeSyncResponsePayload{}),
	ListenRequest:     reflect.TypeOf(ListenRequestPayload{}),
	ListenConfirm:     reflect.TypeOf(ListenConfirmPayload{}),
	Unlisten:          reflect.TypeOf(UnlistenPayload{}),
	KeyValue:          reflect.TypeOf(KeyValuePayload{}),
	CarrierEvent:      reflect.TypeOf(CarrierEventPayload{}),
	VersionInfo:       reflect.TypeOf(VersionInfoPayload{}),
}

func EncodePayload(verb MessageVerb, payload Payload) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(verb)
	if payload != nil {
		err := binary.Write(buf, byteOrder, payload)
		if err != nil {
			log.Fatalf("Bad message encode for %T %v", payload, err)
		}
	}
	return buf.Bytes()
}

func DecodePacket(lineBuffer []byte) (MessageVerb, interface{}) {
	verb := MessageVerb(lineBuffer[0])
	var payloadObj interface{} = nil

	switch verb {
	case EnumerateChannels:
		break
	case ListChannels:
		payloadObj = new(ListChannelsPayload)
	case TimeSync:
		payloadObj = new(TimeSyncPayload)
	case TimeSyncResponse:
		payloadObj = new(TimeSyncResponsePayload)
	case ListenRequest:
		payloadObj = new(ListenRequestPayload)
	case ListenConfirm:
		payloadObj = new(ListenConfirmPayload)
	case Unlisten:
		payloadObj = new(UnlistenPayload)
	case KeyValue:
		payloadObj = new(KeyValuePayload)
	case CarrierEvent:
		payloadObj = new(CarrierEventPayload)
	case VersionInfo:
		payloadObj = new(VersionInfoPayload)

	}
	buffer := bytes.NewReader(lineBuffer[1:])

	if payloadObj != nil {
		err := binary.Read(buffer, byteOrder, payloadObj)
		if err != nil {
			glog.Fatalf("Error reading message for %d, %v", verb, err)
			return verb, nil
		}
	}
	return verb, payloadObj
}
