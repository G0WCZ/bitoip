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
#include "debug.h"
#include "messages.h"
#include "timestamp.h"

// time sync support

// buckets for holding sync calcs
#define TS_BUCKET_SIZE 5

int time_offset_index = 0;

int64_t offsets[TS_BUCKET_SIZE];
int64_t offset_sum = 0;

int64_t roundtrips[TS_BUCKET_SIZE];
int64_t roundtrip_sum = 0;

int64_t common_time_offset = 0;
int64_t common_round_trip = 0;

long time_sync_count = 0;

void handle_time_sync_response(void * payload) {
    TimeSyncResponsePayload * tsrp = (TimeSyncResponsePayload*)payload;
    uint64_t now = timestamp_64_now();

    int64_t latest_time_offset = ((tsrp->server_rx_time - tsrp->given_time) - (tsrp->server_tx_time - now)) / 2;
    int64_t latest_round_trip = (now - tsrp->given_time) - (tsrp->server_rx_time - tsrp->server_tx_time);
    
    offsets[time_offset_index] = latest_time_offset;
    roundtrips[time_offset_index] = latest_round_trip;

    time_offset_index = (time_offset_index + 1) % TS_BUCKET_SIZE;

    offset_sum = 0;
    roundtrip_sum = 0;

    if (time_sync_count < TS_BUCKET_SIZE) {
        time_sync_count++;
    }

    for (int i=0; i<time_sync_count; i++) {
        offset_sum += offsets[i];
        roundtrip_sum += roundtrips[i];
    }

    common_time_offset = offset_sum / time_sync_count;
    common_round_trip = roundtrip_sum / time_sync_count;

    debug_printf("timesync offset: %ld ", ((common_time_offset & 0xffffffff)/1000));
    debug_printf("timesync round trip: %ld\n", ((common_round_trip & 0xffffffff)/1000));
    
}

int64_t get_time_offset() {
    return common_time_offset;
}

int64_t get_roundtrip() {
    return common_round_trip;
}

void timesync_setup() {
    set_handler(TIME_SYNC_RESPONSE, handle_time_sync_response);
}
