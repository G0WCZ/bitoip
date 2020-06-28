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
	"sort"
	"time"

	"github.com/golang/glog"
)

const LastReservedCarrierKey = 99
const StationGoneTimeout = time.Duration(5 * time.Minute)

type Subscriber struct {
	Key        bitoip.CarrierKeyType
	Address    net.UDPAddr
	LastTx     time.Time
	LastListen time.Time
	Callsign   string
}

type Channel struct {
	ChannelId   bitoip.ChannelIdType
	Subscribers map[bitoip.CarrierKeyType]Subscriber
	Addresses   map[string]Subscriber
	Callsigns   map[string]Subscriber
	LastKey     bitoip.CarrierKeyType
}

type ChannelMap map[uint16]*Channel

var channels = make(ChannelMap)

const MaxSeenOn = 5

type Station struct {
	Callsign          string
	SeenOnChannels    []bitoip.ChannelIdType
	LastActive        time.Time
	LastActiveChannel bitoip.ChannelIdType
	LastListen        time.Time
}

func (s *Station) AddSeenOn(cId bitoip.ChannelIdType) {
	seenOn := make([]bitoip.ChannelIdType, 0)

	seenOn = append(seenOn, cId)

	for i, v := range s.SeenOnChannels {
		if s.SeenOnChannels[i] != cId {
			seenOn = append(seenOn, v)
		}
		if len(seenOn) >= MaxSeenOn {
			break
		}
	}

	s.SeenOnChannels = seenOn
}

type StationMap map[string]*Station

var stations = make(StationMap)

// Create a new channel
func NewChannel(channelId bitoip.ChannelIdType) Channel {
	return Channel{
		channelId,
		make(map[bitoip.CarrierKeyType]Subscriber),
		make(map[string]Subscriber),
		make(map[string]Subscriber),
		LastReservedCarrierKey,
	}
}

// Return array of channel Ids of existing channels
func ChannelIds() []uint16 {
	keys := make([]uint16, 0, len(channels))
	for k := range channels {
		keys = append(keys, k)
	}
	return keys
}

// Get a channel by channel_id
func GetChannel(channel_id bitoip.ChannelIdType) *Channel {
	if channel, ok := channels[channel_id]; ok {
		return channel
	} else {
		nc := NewChannel(channel_id)
		channels[channel_id] = &nc
		return &nc
	}
}

func GetStation(callsign string) *Station {
	if station, ok := stations[callsign]; ok {
		return station
	} else {
		s := Station{
			Callsign:          callsign,
			SeenOnChannels:    make([]bitoip.ChannelIdType, 5),
			LastActive:        time.Time{},
			LastListen:        time.Time{},
			LastActiveChannel: 0,
		}
		stations[callsign] = &s
		return &s
	}
}

// Subscribe to this channel
// if already susscribed, then update details and LastTx
func (c *Channel) Subscribe(address net.UDPAddr, callsign string) bitoip.CarrierKeyType {
	glog.V(2).Infof("subscribe from: %v", address)

	subscriber, ok := c.Callsigns[callsign]

	if ok {
		// already have this callsign, so reuse:
		subscriber.LastListen = time.Now()
		c.Addresses[address.String()] = subscriber
		c.Subscribers[subscriber.Key] = subscriber
		c.Callsigns[callsign] = subscriber
		glog.V(2).Infof("subscribe %s existing key %d", callsign, subscriber.Key)

	} else {
		c.LastKey += 1
		subscriber = Subscriber{c.LastKey, address, *new(time.Time), time.Now(), callsign}
		c.Subscribers[c.LastKey] = subscriber
		c.Addresses[address.String()] = subscriber
		c.Callsigns[callsign] = subscriber
		glog.V(2).Infof("subscribe %s new key %d", callsign, subscriber.Key)
	}
	s := GetStation(callsign)
	s.LastListen = time.Now()
	stations[callsign] = s
	glog.V(2).Infof("subscriber is: %v key %v", subscriber, subscriber.Key)

	return subscriber.Key
}

// Unsubscribe from channel
func (c *Channel) Unsubscribe(address net.UDPAddr) {
	if subscriber, ok := c.Addresses[address.String()]; ok {
		delete(c.Subscribers, subscriber.Key)
		delete(c.Addresses, subscriber.Address.String())
		delete(c.Callsigns, subscriber.Callsign)
	}
}

// Broadcast this carrier event to all on this channel
// and always return to sender (who can ignore if they wish, or can use as net sidetone
func (c *Channel) Broadcast(event bitoip.CarrierEventPayload) {
	glog.V(2).Infof("broadcast event %v", event)
	txr, ok := c.Subscribers[event.CarrierKey]

	if !ok {
		glog.Infof("broadcast from unsubscribed key %v dropped", event.CarrierKey)
		return
	}

	for _, v := range c.Subscribers {
		if txr.LastTx.Sub(v.LastListen) < StationGoneTimeout {
			glog.V(2).Infof("sending to subs %v: %v", v.Address, event)
			bitoip.UDPTx(bitoip.CarrierEvent, event, &v.Address)
		}
	}

	txr.LastTx = time.Now()
	c.Subscribers[event.CarrierKey] = txr
	c.Addresses[txr.Address.String()] = txr
	c.Callsigns[txr.Callsign] = txr

	s := GetStation(txr.Callsign)
	s.LastActiveChannel = c.ChannelId
	s.LastActive = txr.LastTx
	s.LastListen = txr.LastTx
	s.AddSeenOn(c.ChannelId)
	stations[txr.Callsign] = s
}

func (c *Channel) GetListenSortedSubscribers() []Subscriber {
	var subs = make([]Subscriber, 0)

	for cs, _ := range c.Callsigns {
		subs = append(subs, c.Callsigns[cs])
	}

	sort.Slice(subs, func(i int, j int) bool {
		return subs[i].LastListen.After(subs[j].LastListen)
	})

	return subs
}

// Check through for subscribers that we haven't seem for a while
// and remove them.
func SuperviseChannels(t time.Time, timeout time.Duration) int {
	removed := 0
	for _, channel := range channels {
		for key, sub := range channel.Subscribers {
			if t.Sub(sub.LastListen) > timeout {
				delete(channel.Subscribers, key)
				removed += 1
				for add, sub := range channel.Addresses {
					if sub.Key == key {
						delete(channel.Addresses, add)
					}
				}
				for call, sub := range channel.Callsigns {
					if sub.Key == key {
						delete(channel.Callsigns, call)
					}
				}
			}
		}
	}

	return removed
}
