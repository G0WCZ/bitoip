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

//
// General hardware support .. for things like LEDs and knobs and stuff
// for now prety basic to support status LED
//
const Sidetone = "sidetone"
const StatusLED = "status_led"
const On = "on"
const Off = "off"

type (
	GeneralIO interface {
		Open() error
		ConfigChanged() error
		SetStatus(string, string)
		GetStatus(string) string
		Close() error
		Name() string
	}
)
