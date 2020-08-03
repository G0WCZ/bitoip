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
// Basic control of CWC station

#include <Arduino.h>
#include "debug.h"
#include "config.h"
#include "dashboard.h"
#include "messages.h"
#include "timesync.h"
#include "udptransport.h"


void listen_confirm(void * payload) {
    ListenConfirmPayload *lcp = (ListenConfirmPayload *)payload;
    debug_printf("Listen confirm for %d with carrier key %d\n", lcp->channel, lcp->carrier_key);
}

void channel_setup() {
    char callsign[CALLSIGN_SIZE];

    set_handler(LISTEN_CONFIRM, listen_confirm);

    get_config("Callsign").toCharArray(callsign, CALLSIGN_SIZE);
    
    debug_printf("callsign is %s", callsign);

    listen_request(0, callsign);
}


void station_setup() {
    dash_set_state(DS_REF_SEEK); // set to "seeking reflector state"
    timesync_setup();
    udp_transport_setup();
    channel_setup();
}

void station_run() {
    udp_transport_run();
    time_sync();
    delay(5000);
}