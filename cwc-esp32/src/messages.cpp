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

#include <sys/time.h>
#include <inttypes.h>
#include <string.h>
#include <WiFi.h>
#include "debug.h"
#include "timestamp.h"
#include "messages.h"

uint64_t
ntoh64(const uint64_t *input)
{
    uint64_t rval;
    uint8_t *data = (uint8_t *)&rval;

    data[0] = *input >> 56;
    data[1] = *input >> 48;
    data[2] = *input >> 40;
    data[3] = *input >> 32;
    data[4] = *input >> 24;
    data[5] = *input >> 16;
    data[6] = *input >> 8;
    data[7] = *input >> 0;

    return rval;
}

uint64_t
hton64(const uint64_t *input)
{
    return (ntoh64(input));
}

void debug_data(char *pkt, int len) {
    debug_print("[");
    for (int i=0; i<len; i++) {
        debug_printf("%x ", pkt[i]);
    }
    debug_print("]\n");
}

void (*sender)(char * pkt, int len) = NULL;
/**
 * Send message via UDP
 */
void msg_send(char * pkt, int len) {
    if (sender != NULL) {
        debug_printf("sending verb 0x%x len %d\n", pkt[0], len);
        debug_data(pkt, len);
        (*sender)(pkt, len);
    } else {
        debug_println("No message sender set");
    }
}

void set_message_sender(void(*s)(char * pkt, int len)) {
    sender = s;
}

// Message encode functions:
// These do marshalling and network byte order conversions

void enumerate_channels() {
    EnumerateChannelsPayload ecp;
    msg_send((char *)&ecp, sizeof(ecp));
}

void time_sync() {
    TimeSyncPayload tsp;

    uint64_t now = timestamp_64_now();
    tsp.current_time = hton64(&now);
     
    msg_send((char *)&tsp, sizeof(tsp));
}

void listen_request(int channel, char * callsign) {
    ListenRequestPayload lrp;
    
    lrp.channel = (ChannelIdType) htons(channel);
    strncpy(lrp.callsign, callsign, 16);

    msg_send((char *)&lrp, sizeof(lrp));
}

void unlisten(int channel, unsigned short key) {
    UnlistenPayload up;
    
    up.channel = (ChannelIdType) htons(channel);
    up.carrier_key = (CarrierKeyType) htons(key);

    msg_send((char *)&up, sizeof(up));
} 

void key_value(char * key, char * value) {
    KeyValuePayload kvp;
    
    strncpy(kvp.key, key, KEY_SIZE);
    strncpy(kvp.value, value, VALUE_SIZE);

    msg_send((char *)&kvp, sizeof(kvp));
}

void carrier_event(int channel, unsigned short key, uint64_t start_time, CarrierBitEvent (*bitEvents)[MAX_BIT_EVENTS]) {
    CarrierEventPayload cep;

    cep.channel = htons(channel);
    cep.carrier_key = htons(key);
    cep.start_timestamp = hton64(&start_time);

    for (int i=0; i<MAX_BIT_EVENTS; i++) {
        cep.bit_events[i].bit_event = (*bitEvents)[i].bit_event;
        cep.bit_events[i].time_offset = htonl((*bitEvents)[i].time_offset);
    } 

    msg_send((char *)&cep, sizeof(cep));
}

PayloadHandler handlers[MAX_VERBS];

void set_handler(unsigned char verb, void (*handler)(void *payload)) {
    handlers[verb-ZERO_VERB] = handler;
}

PayloadHandler get_handler(unsigned char verb){
   return handlers[verb-ZERO_VERB]; 
}

void decode_message(uint8_t * message, int length) {
    char verb = *message;
    void *payload = nullptr;

    debug_printf("got verb 0x%x, length %d\n", verb, length);
    debug_data((char*)message, length);

    switch (verb)
    {
        case LIST_CHANNELS: {
            ListChannelsPayload *lcp = (ListChannelsPayload*)message;
            for (int i=0; i<MAX_CHANNELS_PER_MESSAGE;i++) {
                lcp->channels[i] = ntohs(lcp->channels[i]);
            }
            payload = lcp;
        } break;
        
        case TIME_SYNC_RESPONSE: {
            TimeSyncResponsePayload *tsrp = (TimeSyncResponsePayload*)message;
            tsrp->given_time = ntoh64(&tsrp->given_time);
            tsrp->server_rx_time = ntoh64(&tsrp->server_rx_time);
            tsrp->server_tx_time = ntoh64(&tsrp->server_tx_time);
            payload = tsrp;
        } break;

        case LISTEN_CONFIRM: {
            ListenConfirmPayload *lcop = (ListenConfirmPayload*)message;
            lcop->channel = ntohs(lcop->channel);
            lcop->carrier_key = ntohs(lcop->carrier_key);
            payload = lcop;
        } break;


        default: {
        } break;
    }

    if (payload != nullptr) {
        get_handler(verb)(payload);
    }
}