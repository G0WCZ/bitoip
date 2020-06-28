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

const (
	PWM_CYCLE     uint32 = 32
	ON_DUTY_CYCLE uint32 = 1
)

type GPIOOut struct {
	Config      *config.Config
	output      rpio.Pin
	pwmOut      rpio.Pin
	statusLED   rpio.Pin
	adapterName string
	name        string
}

func NewGPIOOut(config *config.Config, name string, adapterName string) MorseOut {
	gpio_out := GPIOOut{
		Config:      config,
		adapterName: adapterName,
		name:        name,
	}
	return &gpio_out
}

func (g *GPIOOut) Open() error {
	return g.initPorts(g.Config)
}

func (g *GPIOOut) ConfigChanged() error {
	return g.initPorts(g.Config)
}

func (g *GPIOOut) initPorts(config *config.Config) error {
	g.output = rpio.Pin(config.GPIOPins.KeyOut)
	g.output.Output()
	g.output.Low()

	// PWM setup
	g.pwmOut = rpio.Pin(config.GPIOPins.PWMA)
	g.pwmOut.Mode(rpio.Pwm)
	g.pwmOut.Freq(config.SidetoneFrequency * int(PWM_CYCLE))
	g.pwmOut.DutyCycle(0, PWM_CYCLE)

	// Status LED setup
	g.statusLED = rpio.Pin(config.GPIOPins.StatusLED)
	g.statusLED.Output()
	g.statusLED.Low()

	return nil
}

func (g *GPIOOut) Close() error {
	// Turn stuff off to be sure
	g.SetBit(false)
	g.SetStatusLED(false)
	g.SetToneOut(false)

	return nil
}

func (g *GPIOOut) Name() string {
	return g.name
}

func (g *GPIOOut) SetBit(bit bool) {
	if bit {
		g.output.Low()
	} else {
		g.output.High()
	}
}

func (g *GPIOOut) SetToneOut(bit bool) {
	if bit {
		g.pwmOut.DutyCycle(ON_DUTY_CYCLE, PWM_CYCLE)
	} else {
		g.pwmOut.DutyCycle(0, PWM_CYCLE)
	}
}

func (g *GPIOOut) SetStatusLED(bit bool) {
	if bit {
		g.statusLED.Low()
	} else {
		g.statusLED.High()
	}
}
