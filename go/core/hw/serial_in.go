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

type SerialIn struct {
	Config         *config.Config
	adapterName    string
	leftInput      string
	rightInput     string
	keyer          bool
	port           serial.Port
	portDeviceName string
	name           string
}

func NewSerialIn(config *config.Config, name string, adapterName string) MorseIn {
	return &SerialIn{
		Config:         config,
		adapterName:    adapterName,
		keyer:          adapterName == KEYER,
		portDeviceName: config.Serial.Device,
		name:           name,
	}
}

func (s *SerialIn) Open() error {
	return s.ConfigChanged()
}

func (s *SerialIn) ConfigChanged() error {
	if len(s.portDeviceName) > 0 {
		ClosePort(s.portDeviceName)
		s.port = GetPort(s.Config.Serial.Device)
		s.portDeviceName = s.Config.Serial.Device
	}
	return nil
}

func (s *SerialIn) Bit() bool {
	panic("implement me")
}

func (s *SerialIn) Dit() bool {
	panic("implement me")
}

func (s *SerialIn) Dah() bool {
	panic("implement me")
}

func (s *SerialIn) Close() error {
	panic("implement me")
}

func (s *SerialIn) UseKeyer() bool {
	return s.keyer
}

func (s *SerialIn) Name() string {
	return s.name
}
