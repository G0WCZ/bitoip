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
	"github.com/G0WCZ/cwc/config"
	"strings"

	"github.com/golang/glog"
	"go.bug.st/serial.v1"
)

/*
 * Serial Hardware
 * Use RS232 flow control signals for input and output.
 * Should be reasonably portable across Mac, PC, Linux (incl. Pi)
 */

type SerialIO struct {
	config *config.Config
	port   serial.Port
	useRTS bool
	useCTS bool
}

// Create this hardare device
func NewSerialIO(config *config.Config) *SerialIO {
	serialIO := SerialIO{
		config: config,
		port:   nil,
		useRTS: true,
		useCTS: true,
	}
	return &serialIO
}

// Open the port and set bit behaviours
func (s *SerialIO) Open() error {
	serialDevice := s.config.Serial.Device

	glog.Infof("Opening serial port %s", serialDevice)

	mode := &serial.Mode{}

	port, err := serial.Open(serialDevice, mode)
	s.port = port

	if err != nil {
		glog.Fatalf("Can not open serial port: %v", err)
	}

	s.useRTS = strings.EqualFold(s.config.Serial.KeyOut, "RTS")
	s.useCTS = strings.EqualFold(s.config.Serial.KeyLeft, "CTS")

	return nil
}

// Read a morse input bit
func (s *SerialIO) Bit() bool {
	bits, err := s.port.GetModemStatusBits()
	if err != nil {
		glog.Fatalf("Port bit read failed %v", err)
		return false
	}
	if s.useCTS {
		return bits.CTS
	} else {
		return bits.DSR
	}
}
func (s *SerialIO) Dot() bool {
	return false
}
func (s *SerialIO) Dash() bool {
	return false
}

// Send a morse output bit
func (s *SerialIO) SetBit(bit bool) {
	var err error
	if s.useRTS {
		err = s.port.SetRTS(bit)
	} else {
		err = s.port.SetDTR(bit)
	}
	if err != nil {
		glog.Fatalf("port bit set failed: %v", err)
	}
}

// No tone sending supported, so this does nothing
func (s *SerialIO) SetToneOut(_ bool) {}

func (s *SerialIO) SetStatusLED(_ bool) {}

func (s *SerialIO) Close() {
	s.port.Close()
}
