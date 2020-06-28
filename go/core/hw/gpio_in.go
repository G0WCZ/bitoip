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
	"github.com/stianeikeland/go-rpio"
)

const KEYER = "keyer"

type GPIOIn struct {
	Config      *config.Config
	adapterName string
	leftInput   rpio.Pin
	rightInput  rpio.Pin
	keyer       bool
	name        string
}

func NewGPIOIn(config *config.Config, name string, adapterName string) MorseIn {
	return &GPIOIn{
		Config:      config,
		adapterName: adapterName,
		keyer:       adapterName == KEYER,
		name:        name,
	}
}

func (g *GPIOIn) Open() error {
	return g.initPorts(g.Config)
}

func (g *GPIOIn) ConfigChanged() error {
	return g.initPorts(g.Config)
}

func (g *GPIOIn) initPorts(config *config.Config) error {
	_ = rpio.Open()
	g.leftInput = rpio.Pin(config.GPIOPins.KeyLeft)
	g.leftInput.Input()
	g.leftInput.PullUp()

	g.rightInput = rpio.Pin(config.GPIOPins.KeyRight)
	g.rightInput.Input()
	g.rightInput.PullUp()

	return nil
}

func (g *GPIOIn) Close() error {
	return nil
}

func (g *GPIOIn) Bit() bool {
	return g.Dit()
}

func (g *GPIOIn) Name() string {
	return g.name
}

// NB Active low - has a pullup. Keydown pulls to ground
func (g *GPIOIn) Dit() bool {
	if g.leftInput.Read() == rpio.High {
		return false
	} else {
		return true
	}
}

func (g *GPIOIn) Dah() bool {
	if g.rightInput.Read() == rpio.High {
		return false
	} else {
		return true
	}
}

func (g *GPIOIn) UseKeyer() bool {
	return g.keyer
}
