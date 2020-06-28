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

package hw

import (
	"github.com/G0WCZ/cwc/bitoip"
	"time"
)

type LocalBitEvent uint8

// Bit masks for LocalBitEvents
// NOT the same as the bitoip.BitEvent constants!!
const (
	// Mostly for local paddle before keyer sampling
	BitRightOn  LocalBitEvent = 0x20 // Right, Ring
	BitRightOff LocalBitEvent = 0x10 // Right, Ring
	BitLeftOn   LocalBitEvent = 0x08 // Left, Tip
	BitLeftOff  LocalBitEvent = 0x04 // Left, Tip
	// Normal keying
	BitOn  LocalBitEvent = 0x02 // Straight On
	BitOff LocalBitEvent = 0x01 // Straight Off
	// Control
	LastEvent LocalBitEvent = 0x80 // high bit set to indicate last one
)

// Absolute time morse event
type TimedBitEvent struct {
	StartTime time.Time
	BitEvent  bitoip.BitEvent
}

// slice of events
type MorseEvents []TimedBitEvent

//TODO conversion to/from messages
