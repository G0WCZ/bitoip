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

//
// Very simple config system.. to be replaced with something better
//

#include <Arduino.h>
#include <SimpleMap.h>
#include "debug.h"

SimpleMap<String, String> *config = new SimpleMap<String, String>([](String &a, String &b) -> int {
        if (a == b) return 0;
        if (a > b) return 1;
        /*if (a < b) */ return -1;
    });



String get_config(String key) {
    return config->get(key);
}

void set_config(String key, String value) {
    config->put(key, value);
}

void config_setup() {
    config->put("LocalPort", "5990"); // Local UDP Port
    config->put("ReflectorHost", "cwc.onlineradioclub.org");
    config->put("ReflectorPort", "7388");
    config->put("Callsign", "G0WCZ-32");
}