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

const OnPeriods = uint32(1)
const CycleLength = uint32(32)

type GPIOGeneral struct {
	Config      *config.Config
	adapterName string
	usePwm      bool
	pwmOutState string
	pwmOut      rpio.Pin
	useStatus   bool
	status      rpio.Pin
	name        string
}

func (G *GPIOGeneral) Open() error {
	return G.ConfigChanged()
}

func (G *GPIOGeneral) ConfigChanged() error {
	return G.initGPIO()
}

func (G *GPIOGeneral) initGPIO() error {
	_ = rpio.Open()

	sFreq := G.Config.SidetoneFrequency
	glog.Infof("sidetone frequency set to %dHz", sFreq)

	// PCM output
	if sFreq > 0 {
		pcmPinNo := G.Config.GPIOPins.PWMA

		G.usePwm = true
		G.pwmOut = rpio.Pin(pcmPinNo)
		G.pwmOut.Mode(rpio.Pwm)
		G.pwmOut.Freq(int(uint32(sFreq) * CycleLength))
		G.pwmOut.DutyCycle(0, CycleLength)
	}

	statusLEDPin := G.Config.GPIOPins.StatusLED
	if statusLEDPin > 0 {
		G.useStatus = true
		G.status = rpio.Pin(statusLEDPin)
		G.status.Output()
		G.status.Low()
	} else {
		G.useStatus = false
	}

	return nil
}

func (G *GPIOGeneral) SetStatus(name string, value string) {
	if name == Sidetone {
		if G.usePwm {
			if value == On {
				G.pwmOut.DutyCycle(OnPeriods, CycleLength)
			} else {
				G.pwmOut.DutyCycle(0, CycleLength)
			}
		}
	} else if name == StatusLED && G.useStatus {
		if value == On {
			G.status.High()
		} else {
			G.status.Low()
		}
	}
}

func (G *GPIOGeneral) GetStatus(name string) string {
	if name == Sidetone {
		return G.pwmOutState
	} else if name == StatusLED {
		return Off
	} else {
		return Off
	}
}

func (G *GPIOGeneral) Close() error {
	G.SetStatus(Sidetone, Off)
	G.SetStatus(StatusLED, Off)
	return nil
}
func (G *GPIOGeneral) Name() string {
	return G.name
}

func NewGPIOGeneral(config *config.Config, name string, adapterName string) GeneralIO {
	return &GPIOGeneral{Config: config, adapterName: adapterName,
		usePwm: false, pwmOutState: Off, name: name}
}
