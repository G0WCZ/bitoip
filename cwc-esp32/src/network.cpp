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
// Wifi, basic network and UDP stuff
// wifi, then ntp, then...
#include <WiFi.h>
#include <time.h>
#include "debug.h"
#include "dashboard.h"


const char* ntpServer = "pool.ntp.org";
const long  gmtOffset_sec = 0;
const int   daylightOffset_sec = 0;

void network_setup(char *ssid, char *password) {
    delay(10000);
    debug_printf("Establishing connection to wifi %s\n", ssid);

    WiFi.begin(ssid, password);
    
    while (WiFi.status() != WL_CONNECTED) {
        delay(100);
    }
    dash_set_state(DS_GOT_WIFI); 
    dash_set_key(String("LocalIP"), WiFi.localIP().toString());
    debug_print("Wifi connected. IP is ");
    debug_println(WiFi.localIP());

    // NTP setup
    configTime(gmtOffset_sec, daylightOffset_sec, ntpServer);

    struct tm timeinfo;

    if(!getLocalTime(&timeinfo)){
      debug_println("Failed to obtain ntp time");
    } else {
      debug_println(&timeinfo, "Time set to: %A, %B %d %Y %H:%M:%S");
    }
    
}

