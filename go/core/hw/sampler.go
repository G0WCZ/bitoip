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
)
import (
	"context"
)

//
// A sampler turns bit sampling into an event stream.
// It takes bit data and produces a channel event stream.  It samples from a morse_in
// Examples:  a straight key, a keyer

type Sampler interface {
	Open(config *config.Config, ctx context.Context, morseEvents chan MorseEvents, morseIn MorseIn) error
	ConfigChanged() error
	Close()
}
