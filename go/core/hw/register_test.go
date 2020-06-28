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
	"fmt"
	"github.com/G0WCZ/cwc/config"
	"gotest.tools/assert"
	"testing"
)

func TestSplitDescriptors(t *testing.T) {
	result := SplitDescriptors([]string{"a", "b", "x:y"})
	expected := [][]string{{"a", ""}, {"b", ""}, {"y", "x"}}
	assert.DeepEqual(t, expected, result)
}

func TestParseInputs(t *testing.T) {
	config := config.Config{}
	config.MorseInHardware = []string{"gpio"}
	inputs := ParseInputs(&config)
	it := fmt.Sprintf("%T", inputs[0])
	assert.Equal(t, "*hw.GPIOIn", it)
}
func TestParseInputsWithAdaptor(t *testing.T) {
	config := config.Config{}
	config.MorseInHardware = []string{"keyer:gpio"}
	inputs := ParseInputs(&config)
	it := fmt.Sprintf("%T", inputs[0])
	assert.Equal(t, "*hw.GPIOIn", it)
	assert.Equal(t, "keyer", inputs[0].(*GPIOIn).adapterName)
}

func TestParseOutputs(t *testing.T) {
	config := config.Config{}
	config.MorseOutHardware = []string{"gpio"}
	inputs := ParseOutputs(&config)
	it := fmt.Sprintf("%T", inputs[0])
	assert.Equal(t, "*hw.GPIOOut", it)
}
func TestParseOutputsWithAdaptor(t *testing.T) {
	config := config.Config{}
	config.MorseOutHardware = []string{"xxx:gpio"}
	inputs := ParseOutputs(&config)
	it := fmt.Sprintf("%T", inputs[0])
	assert.Equal(t, "*hw.GPIOOut", it)
	assert.Equal(t, "xxx", inputs[0].(*GPIOOut).adapterName)
}

func TestParseGenerals(t *testing.T) {
	config := config.Config{}
	config.GeneralHardware = []string{"gpio"}
	inputs := ParseGeneralIOs(&config)
	it := fmt.Sprintf("%T", inputs[0])
	assert.Equal(t, "*hw.GPIOGeneral", it)
}
func TestParseGeneralsWithAdaptor(t *testing.T) {
	config := config.Config{}
	config.GeneralHardware = []string{"yyy:gpio"}
	inputs := ParseGeneralIOs(&config)
	it := fmt.Sprintf("%T", inputs[0])
	assert.Equal(t, "*hw.GPIOGeneral", it)
	assert.Equal(t, "yyy", inputs[0].(*GPIOGeneral).adapterName)
}
