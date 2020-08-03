/*
Copyright (C) 2020 Graeme Sutherland, Nodestone Limited


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
#define DEBUG true // flag to turn on/off debugging

#define debug_begin(...) do { if (DEBUG) { Serial.begin(__VA_ARGS__); while(!Serial); }} while (0)
#define debug_print(...) do { if (DEBUG) Serial.print(__VA_ARGS__); } while (0)
#define debug_write(...) do { if (DEBUG) Serial.write(__VA_ARGS__); } while (0)
#define debug_println(...) do { if (DEBUG) Serial.println(__VA_ARGS__); } while (0)
#define debug_printf(...) do { if (DEBUG) Serial.printf(__VA_ARGS__); } while (0)