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
package main

import "fmt"
import "../bitoip"

/*
 * Protocol Version using semantic versioning
 * See: https://semver.org/
 */
var stationVersion = "3.2.0.alpha-0"

func ReflectorVersion() string {
	return stationVersion
}

func DisplayVersion() string {
	return fmt.Sprintf("CWC Station %s / Protocol %s", stationVersion, bitoip.ProtocolVersionString())
}
