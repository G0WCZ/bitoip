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
#include <WiFi.h>
#include "debug.h"
#include "config.h"
#include "dashboard.h"
#include "network.h"
#include "station.h"



#define LED_BUILTIN 2

const char *ssid = "$NETWORK";
const char *password = "$PASSWORD";

void setup() {
  config_setup();
  dash_setup();
  dash_set_state(DS_BOOTED);
  debug_begin(9600);
  debug_println("Starting debugging");

  network_setup((char *)ssid, (char *)password);
  station_setup();
  debug_kvp();
}

void loop() {
  // put your main code here, to run repeatedly:
  //digitalWrite(LED_BUILTIN, HIGH);
  //delay(1000);
  //digitalWrite(LED_BUILTIN, LOW);
  //delay(1000);
}