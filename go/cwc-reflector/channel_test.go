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
	"context"
	"github.com/G0WCZ/cwc/bitoip"
	"net"
	"sort"
	"testing"
	"time"

	"github.com/golang/glog"
	"gotest.tools/assert"
)

func TestNewChannel(t *testing.T) {
	channel := NewChannel(33)
	assert.DeepEqual(t, channel.ChannelId, uint16(33))
	assert.Equal(t, len(channel.Subscribers), 0)
	assert.Equal(t, len(channel.Addresses), 0)
	assert.Equal(t, channel.LastKey, uint16(99))
}

func TestGetChannel(t *testing.T) {
	channel1 := GetChannel(21)
	channel2 := GetChannel(21)
	assert.DeepEqual(t, channel1, channel2)
}

func TestSubscribeWhenNotSubscribed(t *testing.T) {
	channel1 := GetChannel(21)
	addr, _ := net.ResolveUDPAddr("udp", "localhost:2020")
	channel1.Subscribe(*addr, "G0WCZ")
	assert.Equal(t, len(channel1.Addresses), 1)
	assert.Equal(t, len(channel1.Subscribers), 1)
}

func TestSubscribeWhenSubscribed(t *testing.T) {
	channel1 := GetChannel(21)
	addr, _ := net.ResolveUDPAddr("udp", "localhost:2020")
	channel1.Subscribe(*addr, "G0WCZ")
	assert.Equal(t, len(channel1.Addresses), 1)
	assert.Equal(t, len(channel1.Subscribers), 1)
	channel1.Subscribe(*addr, "G0WCZ")
	assert.Equal(t, len(channel1.Addresses), 1)
	assert.Equal(t, len(channel1.Subscribers), 1)
}
func TestUnsubscribeWhenSubscribed(t *testing.T) {
	channel2 := GetChannel(22)
	addr, _ := net.ResolveUDPAddr("udp", "localhost:2020")
	channel2.Subscribe(*addr, "G0WCZ")
	assert.Equal(t, len(channel2.Addresses), 1)
	assert.Equal(t, len(channel2.Subscribers), 1)
	channel2.Unsubscribe(*addr)
	assert.Equal(t, len(channel2.Subscribers), 0)
	assert.Equal(t, len(channel2.Addresses), 0)
}

func TestUnsubscribeWhenNotSubscribed(t *testing.T) {
	channel2 := GetChannel(22)
	addr, _ := net.ResolveUDPAddr("udp", "localhost:2020")
	channel2.Unsubscribe(*addr)
	assert.Equal(t, len(channel2.Subscribers), 0)
	assert.Equal(t, len(channel2.Addresses), 0)
}

func sortSlice(sl []uint16) []uint16 {
	sort.Slice(sl, func(i, j int) bool { return sl[i] < sl[j] })
	return sl
}

func TestChannelIds(t *testing.T) {
	GetChannel(21)
	GetChannel(22)
	GetChannel(33)
	assert.DeepEqual(t, sortSlice(ChannelIds()), sortSlice([]uint16{21, 22, 33}))
}

func TestEmptyChannelIds(t *testing.T) {
	channels = make(map[uint16]*Channel)
	assert.Equal(t, len(ChannelIds()), 0)
}

func carrierEventPayload(key bitoip.CarrierKeyType) bitoip.CarrierEventPayload {
	return bitoip.CarrierEventPayload{
		1,
		key,
		time.Now().UnixNano(),
		[bitoip.MaxBitEvents]bitoip.CarrierBitEvent{
			{0, bitoip.BitOn},
			{100, bitoip.BitOff | bitoip.LastEvent},
		},
		int64(0),
	}
}

func TestBroadcastEmpty(t *testing.T) {
	channels = make(map[uint16]*Channel)
	c1 := GetChannel(1)
	ce := carrierEventPayload(100)

	c1.Broadcast(ce)
}

func TestBroadcastToSubscriber(t *testing.T) {
	c1 := GetChannel(1)
	add := "localhost:9453"
	addr, _ := net.ResolveUDPAddr("udp", add)
	glog.Infof("addr: %v", addr)
	key := c1.Subscribe(*addr, "G0WCZ")
	ce := carrierEventPayload(key)

	pc, _ := net.ListenPacket("udp", add)
	buffer := make([]byte, bitoip.MaxMessageSizeInBytes)
	doneChan := make(chan []byte, 1)

	// get one message
	go func() {
		_, _, _ = pc.ReadFrom(buffer)
		glog.Infof("Raw Rx: %d %v", len(buffer), buffer)
		doneChan <- buffer
	}()

	serverAddress, _ := net.ResolveUDPAddr("udp", "localhost:6012")
	ctx := context.Background()
	messages := make(chan bitoip.RxMSG)
	go bitoip.UDPRx(ctx, serverAddress, messages)

	// delay for connection to be established
	time.Sleep(time.Second * 2)

	// broadcast
	c1.Broadcast(ce)

	buf := <-doneChan

	verb, payload := bitoip.DecodePacket(buf)
	assert.Equal(t, verb, bitoip.CarrierEvent)
	assert.DeepEqual(t, payload, &ce)
	rxce := payload.(*bitoip.CarrierEventPayload)
	assert.Equal(t, rxce.CarrierKey, key)
}

func TestSuperviseChannelsNoSubscribers(t *testing.T) {
	channels = make(map[uint16]*Channel)
	_ = GetChannel(1)
	_ = GetChannel(2)
	r := SuperviseChannels(time.Now(), time.Duration(10*time.Minute))
	assert.Equal(t, r, 0)
}

func TestSuperviseChannelsNoneRemoved(t *testing.T) {
	channels = make(map[uint16]*Channel)
	c1 := GetChannel(1)
	c2 := GetChannel(2)
	addr, _ := net.ResolveUDPAddr("udp", "localhost:19234")
	c1.Subscribe(*addr, "A1AAA")
	c2.Subscribe(*addr, "A1BBB")
	r := SuperviseChannels(time.Now(), time.Duration(10*time.Minute))
	assert.Equal(t, r, 0)
}

func TestSuperviseChannels2Removed(t *testing.T) {
	channels = make(map[uint16]*Channel)
	c1 := GetChannel(1)
	c2 := GetChannel(2)
	addr, _ := net.ResolveUDPAddr("udp", "localhost:19234")
	c1.Subscribe(*addr, "A1AAA")
	c2.Subscribe(*addr, "A1BBB")
	r := SuperviseChannels(time.Now().Add(time.Duration(20*time.Minute)), time.Duration(10*time.Minute))
	assert.Equal(t, r, 2)
	r = SuperviseChannels(time.Now().Add(time.Duration(20*time.Minute)), time.Duration(10*time.Minute))
	assert.Equal(t, r, 0)
}

func TestStation_AddSeenOn(t *testing.T) {
	s := Station{
		"A1AAA",
		nil,
		time.Now(),
		1,
		time.Now(),
	}
	s.AddSeenOn(bitoip.ChannelIdType(1))
	assert.Equal(t, uint16(1), s.SeenOnChannels[0])
	s.AddSeenOn(bitoip.ChannelIdType(1))
	assert.Equal(t, 1, len(s.SeenOnChannels))
	s.AddSeenOn(bitoip.ChannelIdType(2))
	s.AddSeenOn(bitoip.ChannelIdType(3))
	s.AddSeenOn(bitoip.ChannelIdType(4))
	s.AddSeenOn(bitoip.ChannelIdType(5))
	assert.DeepEqual(t,
		s.SeenOnChannels,
		[]bitoip.ChannelIdType{5, 4, 3, 2, 1},
	)
}
