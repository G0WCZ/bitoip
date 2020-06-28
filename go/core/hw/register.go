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
	"strings"
)

var Inputs = map[string]func(*config.Config, string, string) MorseIn{
	"gpio":   NewGPIOIn,
	"mock":   NewMockIn,
	"serial": NewSerialIn,
}

var Outputs = map[string]func(*config.Config, string, string) MorseOut{
	"gpio":   NewGPIOOut,
	"mock":   NewMockOut,
	"serial": NewSerialOut,
}

var Generals = map[string]func(*config.Config, string, string) GeneralIO{
	"gpio": NewGPIOGeneral,
	"mock": NewMockGeneral,
}

func SplitDescriptors(descriptors []string) [][]string {
	var parsed [][]string

	for _, s := range descriptors {
		splits := strings.Split(s, ":")
		if len(splits) == 2 {
			parsed = append(parsed, []string{strings.TrimSpace(splits[1]), strings.TrimSpace(splits[0])})
		} else {
			parsed = append(parsed, []string{strings.TrimSpace(splits[0]), ""})
		}
	}

	return parsed
}

func ParseInputs(config *config.Config) []MorseIn {
	var inputs []MorseIn

	for _, i := range SplitDescriptors(config.MorseInHardware) {
		morseIn, ok := Inputs[i[0]]
		if ok {
			inputs = append(inputs, morseIn(config, i[0], i[1]))
		} else {
			glog.Errorf("Unknown input %s", i[0])
		}
	}
	return inputs
}

func ParseOutputs(config *config.Config) []MorseOut {
	var outputs []MorseOut

	for _, i := range SplitDescriptors(config.MorseOutHardware) {
		morseOut, ok := Outputs[i[0]]
		if ok {
			outputs = append(outputs, morseOut(config, i[0], i[1]))
		} else {
			glog.Errorf("Unknown output %s", i[0])
		}
	}
	return outputs
}

func ParseGeneralIOs(config *config.Config) []GeneralIO {
	var generals []GeneralIO

	for _, i := range SplitDescriptors(config.GeneralHardware) {
		general, ok := Generals[i[0]]
		if ok {
			generals = append(generals, general(config, i[0], i[1]))
		} else {
			glog.Errorf("Unknown general io %s", i[0])
		}
	}
	return generals
}
