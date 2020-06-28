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
	"time"
)

const (
	CHECK       int = 0
	PREDOT      int = 1
	PREDASH     int = 2
	SENDDOT     int = 3
	SENDDASH    int = 4
	DOTDELAY    int = 5
	DASHDELAY   int = 6
	DOTHELD     int = 7
	DASHHELD    int = 8
	LETTERSPACE int = 9
	EXITLOOP    int = 10
)

type KeyerState struct {
	state      int
	dotMemory  bool
	dashMemory bool
	kDelay     int
	dotDelay   int
	dashDelay  int
	out        bool
	weight     int
	mode       int
	speed      int
	reverse    bool
	spacing    bool
}

func (k *KeyerState) clearMemory() {
	k.dotMemory = false
	k.dashMemory = false
}

var keyerStates []*KeyerState

func InitKeyer(config *config.Config) {
	dotDelay := 1200 / config.Keyer.Speed

	keyerStates = append(keyerStates, &KeyerState{
		state:      CHECK,
		dotMemory:  false,
		dashMemory: false,
		kDelay:     0,
		dotDelay:   dotDelay,
		dashDelay:  dotDelay * 3 * config.Keyer.Weight / 50,
		out:        false,
		weight:     config.Keyer.Weight,
		mode:       config.Keyer.Mode,
		speed:      config.Keyer.Speed,
		reverse:    config.Keyer.Reverse,
		spacing:    config.Keyer.LetterSpace,
	})
}

func ResetKeyers() {
	keyerStates = []*KeyerState{}
}

// Bit samplers including reverse
func SampleDitPaddle(reverse bool, input MorseIn) bool {
	if reverse {
		return input.Dah()
	} else {
		return input.Dit()
	}
}

func SampleDahPaddle(reverse bool, input MorseIn) bool {
	return SampleDitPaddle(!reverse, input)
}

func KeyerSampleInput(t time.Time, inputIdx int, input MorseIn) bool {
	k := keyerStates[inputIdx]

	switch k.state {
	case CHECK:
		if input.Dit() {
			k.state = PREDOT
		} else if input.Dah() {
			k.state = PREDASH
		}

	case PREDOT:
		glog.V(2).Infof("PREDOT")
		k.clearMemory()
		k.state = SENDDOT

	case PREDASH:
		glog.V(2).Infof("PREDASH")
		k.clearMemory()
		k.state = SENDDASH

	// dot paddle  pressed so set keyer_out high for time dependant on speed
	// also check if dash paddle is pressed during this time
	case SENDDOT:
		glog.V(2).Infof("SENDDOT")
		k.out = true
		if k.kDelay == k.dotDelay {
			k.kDelay = 0
			k.out = false
			k.state = DOTDELAY // add inter-character spacing of one dot length
		} else {
			k.kDelay++
		}

		// if Mode A and both paddels are relesed then clear dash memory
		if k.mode == 0 {
			if !SampleDitPaddle(k.reverse, input) && !SampleDahPaddle(k.reverse, input) {
				k.dashMemory = false
			} else {
				if SampleDahPaddle(k.reverse, input) {
					k.dashMemory = true
				}
			}
		}

	// dash paddle pressed so set keyer_out high for time dependant on 3 x dot delay and weight
	// also check if dot paddle is pressed during this time
	case SENDDASH:
		glog.V(2).Infof("SENDDASH")
		k.out = true
		if k.kDelay == k.dashDelay {
			k.kDelay = 0
			k.out = false
			k.state = DASHDELAY // add inter-character spacing of one dash length
		} else {
			k.kDelay++
		}

		// if Mode A and both padles are relesed then clear dot memory
		if k.mode == 0 {
			if !SampleDitPaddle(k.reverse, input) && !SampleDahPaddle(k.reverse, input) {
				k.dotMemory = false
			} else {
				if SampleDitPaddle(k.reverse, input) {
					k.dotMemory = true
				}
			}
		}

	// add dot delay at end of the dot and check for dash memory, then check if paddle still held
	case DOTDELAY:
		if k.kDelay == k.dotDelay {
			k.kDelay = 0
			if k.dashMemory { // dash has been set during the dot so service
				k.state = PREDASH
			} else {
				k.state = DOTHELD // dot is still active so service
			}
		} else {
			k.kDelay++
		}

		// set dash memory
		if SampleDahPaddle(k.reverse, input) {
			k.dashMemory = true
		}

	// add dot delay at end of the dash and check for dot memory, then check if paddle still held
	case DASHDELAY:
		if k.kDelay == k.dotDelay {
			k.kDelay = 0
			if k.dotMemory { // dot has been set during the dash so service
				k.state = PREDOT
			} else {
				k.state = DASHHELD // dash is still active so service
			}
		} else {
			k.kDelay++
		}

		// set dot memory
		if SampleDitPaddle(k.reverse, input) {
			k.dotMemory = true
		}

	// check if dot paddle is still held, if so repeat the dot. Else check if Letter space is required
	case DOTHELD:
		glog.V(2).Infof("DOTHELD")
		if SampleDitPaddle(k.reverse, input) { // dot has been set during the dash so service
			k.state = PREDOT
		} else {
			if SampleDahPaddle(k.reverse, input) { // has dash paddle been pressed
				k.state = PREDASH
			} else {
				if k.spacing { // Letter space enabled so clear any pending dots or dashes
					k.clearMemory()
					k.state = LETTERSPACE
				} else {
					k.state = EXITLOOP
				}
			}
		}

	// check if dash paddle is still held, if so repeat the dash. Else check if Letter space is required
	case DASHHELD:
		glog.V(2).Infof("DASHHELD")
		if SampleDahPaddle(k.reverse, input) { // dash has been set during the dot so service
			k.state = PREDASH
		} else {
			if SampleDahPaddle(k.reverse, input) { // has dot paddle been pressed
				k.state = PREDOT
			} else {
				if k.spacing { // Letter space enabled so clear any pending dots or dashes
					k.clearMemory()
					k.state = LETTERSPACE
				} else {
					k.state = EXITLOOP
				}
			}
		}

	// Add letter space (3 x dot delay) to end of character and check if a paddle is pressed during this time.
	// Actually add 2 x dot_delay since we already have a dot delay at the end of the character.
	case LETTERSPACE:
		glog.V(2).Infof("LETTERSPACE")
		if k.kDelay == 2*k.dotDelay {
			k.kDelay = 0
			if k.dotMemory { // check if a dot or dash paddle was pressed during the delay
				k.state = PREDOT
			} else {
				if k.dashMemory {
					k.state = PREDASH
				} else {
					k.state = EXITLOOP // no memories set so restart
				}
			}
		} else {
			k.kDelay++
		}

		// save any key presses during the letter space delay
		if SampleDitPaddle(k.reverse, input) {
			k.dotMemory = true
		}
		if SampleDahPaddle(k.reverse, input) {
			k.dashMemory = true
		}

	case EXITLOOP:
		glog.V(2).Infof("EXITLOOP")
		k.state = CHECK

	default:
		glog.V(2).Infof("default case")
		k.state = EXITLOOP
	}

	return k.out
}
