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
	"github.com/golang/glog"
	"go.bug.st/serial.v1"
)

/**
 * Used to allow sharing of a serial device between an input and an output driver.
 */

// Map of open port devices
// so can be used by both input and output
var devices map[string]serial.Port

// Get pointer to serial port, opening if needed
func GetPort(serialDevice string) serial.Port {
	port, ok := devices[serialDevice]
	if !ok {
		mode := &serial.Mode{}

		newPort, err := serial.Open(serialDevice, mode)

		if err != nil {
			glog.Fatalf("Can not open serial port: %v", err)
		}

		devices[serialDevice] = newPort
		return newPort
	} else {
		return port
	}
}

func ClosePort(serialDevice string) {
	port, ok := devices[serialDevice]
	if ok {
		port.Close()
		delete(devices, serialDevice)
	}
}
