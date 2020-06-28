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
	"go.bug.st/serial.v1"
)

type SerialOut struct {
	Config         *config.Config
	port           serial.Port
	portDeviceName string
	adapterName    string
	useDTR         bool
	name           string
}

func NewSerialOut(c *config.Config, name string, adapterName string) MorseOut {
	s_out := SerialOut{
		Config:         c,
		portDeviceName: c.Serial.Device,
		adapterName:    adapterName,
		port:           nil,
		useDTR:         c.Serial.KeyOut == config.SerialPinDTR,
		name:           name,
	}
	return &s_out
}

func (s *SerialOut) Open() error {
	return s.ConfigChanged()
}

func (s *SerialOut) ConfigChanged() error {
	if len(s.portDeviceName) > 0 {
		ClosePort(s.portDeviceName)
		s.port = GetPort(s.Config.Serial.Device)
		s.portDeviceName = s.Config.Serial.Device
	}
	return nil
}

func (s *SerialOut) SetBit(bit bool) {
	if s.port != nil {
		if s.Config.Serial.KeyOut == config.SerialPinDTR {
			_ = s.port.SetDTR(bit)
		} else {
			_ = s.port.SetRTS(bit)
		}
	}
}

func (s *SerialOut) SetToneOut(bool) {
	// pass
}

func (s *SerialOut) SetStatusLED(bool) {
	// pass
}

func (s *SerialOut) Close() error {
	panic("implement me")
}

func (s *SerialOut) Name() string {
	return s.name
}
