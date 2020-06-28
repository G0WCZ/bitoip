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
	"github.com/golang/glog"
	"github.com/stianeikeland/go-rpio"
)

/*
 * PI GPIO Hardware
 * This does the physical connection to Raspberry PI GPIO pins for
 * input, output and PWM.
 */

// PWM settings
const OnDutyCycle = uint32(1)
const PWMCycleLength = uint32(32)

type PiGPIO struct {
	config *config.Config
	output rpio.Pin
	input  rpio.Pin
	pwm    bool
	pwmOut rpio.Pin
	status rpio.Pin
}

func NewPiGPIO(config *config.Config) *PiGPIO {
	pigpio := PiGPIO{
		config: config,
		pwm:    false,
	}
	return &pigpio
}

// Set up inputs and outputs
func (g *PiGPIO) Open() error {
	err := rpio.Open()
	if err != nil {
		return err
	}
	sFreq := g.config.SidetoneFrequency

	glog.Infof("setting sidetone to %d", sFreq)

	// PCM output
	if sFreq > 0 {
		pcmPinNo := g.config.GPIOPins.PWMA

		g.pwm = true
		g.pwmOut = rpio.Pin(pcmPinNo)
		g.pwmOut.Mode(rpio.Pwm)
		g.pwmOut.Freq(sFreq * 32)
		g.pwmOut.DutyCycle(0, 32)
	}

	// sending morse to a GPIO
	outLED := g.config.GPIOPins.SignalLED

	// receiving morse from a GPIO
	inPin := g.config.GPIOPins.KeyLeft

	statusPin := g.config.GPIOPins.StatusLED

	// Pin output
	g.output = rpio.Pin(outLED)
	g.output.Output()
	g.output.Low()

	// Input pin
	g.input = rpio.Pin(inPin)
	g.input.Input()
	g.input.PullUp()

	// Status LED
	g.status = rpio.Pin(statusPin)
	g.status.Output()
	g.status.Low()

	return nil
}

// ready Morse In hardware
func (g *PiGPIO) Bit() bool {
	if g.input.Read() == rpio.High {
		return false
	} else {
		return true
	}
}

func (g *PiGPIO) Dot() bool {
	return false
}
func (g *PiGPIO) Dash() bool {
	return false
}

// Set Morse Out hardware
func (g *PiGPIO) SetBit(bit0 bool) {
	if bit0 {
		g.output.High()
		g.SetToneOut(true)
	} else {
		g.output.Low()
		g.SetToneOut(false)
	}
}

// Set PWM on/off
func (g *PiGPIO) SetToneOut(v bool) {
	if g.pwm {
		var dutyLen uint32

		if v {
			dutyLen = OnDutyCycle
		} else {
			dutyLen = 0
		}

		g.pwmOut.DutyCycle(dutyLen, PWMCycleLength)
	}
}

// Close the interface
func (g *PiGPIO) Close() {
	g.status.Low()
}

func (g *PiGPIO) SetStatusLED(s bool) {
	if s {
		g.status.High()
	} else {
		g.status.Low()
	}
}
