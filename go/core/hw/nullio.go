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

import "github.com/G0WCZ/cwc/config"

// Null I/O type
// useful for testing, doesn't aim to do anything
// with actual hardware

type NullIO struct {
	config *config.Config
	state  State
}

type State struct {
	Bitin   bool
	Bitout  bool
	Toneout bool
}

func NewNullIO(config *config.Config) *NullIO {
	return &NullIO{
		config: config,
		state:  State{false, false, false},
	}
}

func (n *NullIO) Open() error {
	return nil
}

func (n *NullIO) Config() *config.Config {
	return n.config
}

func (n *NullIO) SetState(state State) {
	n.state = state
}

func (n *NullIO) State() State {
	return n.state
}

func (n *NullIO) Bit() bool {
	return n.state.Bitin
}
func (n *NullIO) Dot() bool {
	return false
}
func (n *NullIO) Dash() bool {
	return false
}

func (n *NullIO) SetBit(b bool) {
	n.state.Bitout = b
}

func (n *NullIO) SetToneOut(b bool) {
	n.state.Toneout = b
}

func (n *NullIO) SetStatusLED(s bool) {}

func (*NullIO) Close() {}
