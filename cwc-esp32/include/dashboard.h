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
#include <Arduino.h>

#define DS_ZERO 0
#define DS_BOOTED 1
#define DS_GOT_WIFI 2
#define DS_WIFI_AP 3
#define DS_REF_SEEK 4
#define DS_REF_SYNC 5
#define DS_REF_LOST 6
#define DS_ERROR 7

#define DS_MAX_STATES 8

void dash_setup();

void dash_set_state(int state);

void dash_set_key(String key, String value);

void dash_unset_key(String key);

void debug_kvp(); 
