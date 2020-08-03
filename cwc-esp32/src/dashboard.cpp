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
#include <SimpleMap.h>
#include "debug.h"
#include "dashboard.h"

#define LED_FLASH_RANGE 12
#define LED_PIN 2

using namespace std;

// MAPPING FOR LED on/off
bool led_flash_map[DS_MAX_STATES][LED_FLASH_RANGE] = {
    {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // DS_ZERO
    {1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, // DS_BOOTED 
    {1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0}, // DS_GOT_WIFI
    {1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0}, // DS_WIFI_AP
    {1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0}, // DS_REF_SEEK
    {1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0}, // DS_REF_SYNC
    {1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0}, // DS_REF_LOST
    {1, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1}  // DS_ERROR
};

volatile int state_count = 0;
int state = DS_ZERO;
hw_timer_t *timer = NULL;

SimpleMap <String, String> *kvp = NULL;

void set_status_led(bool on) {
    digitalWrite(LED_PIN, on);
}


// triggered by timer interrupt
void IRAM_ATTR next_state() {
    set_status_led(led_flash_map[state][state_count]);
    state_count = (state_count + 1) % LED_FLASH_RANGE;
}

void IRAM_ATTR start_timer() {
    timer = timerBegin(3, 80, true); // timer_id = 3; divider=80; countUp = true;
    timerAttachInterrupt(timer, &next_state, true); // edge = true
    timerAlarmWrite(timer, 250000, true);  //250 ms
    timerAlarmEnable(timer);
}

void stop_timer() {
    timerEnd(timer);
    timer = NULL;
}


void dash_setup() {
    state = DS_ZERO;
    state_count = 0;

    kvp = new SimpleMap<String, String>([](String &a, String &b) -> int {
        if (a == b) return 0;
        if (a > b) return 1;
        /*if (a < b) */ return -1;
    });

    pinMode(LED_PIN, OUTPUT);
    set_status_led(false);
    start_timer();
}

void dash_set_state(int new_state) {
    state = new_state;
}

void dash_set_key(String key, String value) {
    kvp->put(key, value);
}

void dash_unset_key(String key) {
    kvp->remove(kvp->getIndex(key));
}

void debug_kvp() {
    for (int i=0; i<kvp->size(); i++) {
        debug_println("===");
        debug_println(kvp->getKey(i) + ": " + kvp->getData(i));
        debug_println("===");
    }
}