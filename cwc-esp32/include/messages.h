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

#include <inttypes.h>

#define MAX_MESSAGE_SIZE 200
#define CALLSIGN_SIZE 16
#define KEY_SIZE 8
#define VALUE_SIZE 16

typedef const char MessageVerb;

typedef unsigned short ChannelIdType;
typedef unsigned short CarrierKeyType;
typedef char CallsignType[CALLSIGN_SIZE];
typedef char KeyType[KEY_SIZE];
typedef char ValueType[VALUE_SIZE];
typedef uint64_t time64;
typedef uint32_t time32;
typedef char BitEvent;

// Messages
#define ENUMERATE_CHANNELS 0x90
#define LIST_CHANNELS 0x91
#define TIME_SYNC 0x92
#define TIME_SYNC_RESPONSE 0x93
#define LISTEN_REQUEST 0x94
#define LISTEN_CONFIRM 0x95
#define UNLISTEN 0x96
#define VERSION_INFO 0x97
#define KEYVALUE 0x81
#define CARRIER_EVENT 0x82

#define MAX_CHANNELS_PER_MESSAGE (MAX_MESSAGE_SIZE - 1)/2

#define ZERO_VERB 0x81
#define MAX_VERBS 0x20

typedef void (*PayloadHandler)(void *payload);

typedef struct {
    char verb = ENUMERATE_CHANNELS;
} EnumerateChannelsPayload;

typedef struct {
    char verb = LIST_CHANNELS;
    char pad = 0x00;
    ChannelIdType channels[MAX_CHANNELS_PER_MESSAGE];
} ListChannelsPayload;

typedef struct {
    char verb = TIME_SYNC;
    char pad[3] = {0x00, 0x00, 0x00};
	time64 current_time;
} TimeSyncPayload; 

typedef struct {
    char verb = TIME_SYNC_RESPONSE;
	time64 given_time __attribute__((packed));
    time64 server_rx_time __attribute__((packed));
    time64 server_tx_time __attribute__((packed));
} TimeSyncResponsePayload; 

typedef struct {
    char verb = LISTEN_REQUEST;
    char pad = 0x00;
    ChannelIdType channel;
    CallsignType callsign;
} ListenRequestPayload;

typedef struct {
    char verb = LISTEN_CONFIRM;
    ChannelIdType channel;
    CarrierKeyType carrier_key;
} ListenConfirmPayload;

typedef struct {
    char verb = UNLISTEN;
    char pad = 0x00;
    ChannelIdType channel;
    CarrierKeyType carrier_key;
} UnlistenPayload;

typedef struct {
    char verb = KEYVALUE;
    char pad = 0x00;
    ChannelIdType channel;
    CarrierKeyType carrier_key;
    KeyType key;
    ValueType value;
} KeyValuePayload;

// Paddle keying
#define BIT_RIGHT_ON 0x10  // Right, Ring
#define BIT_RIGHT_OFF 0x08 // Right, Ring
#define BIT_LEFT_ON	0x04   // Left, Tip
#define	BIT_LEFT_OFF 0x02  // Left, Tip

// Normal keying
#define BIT_ON 0x01        // Straight On
#define BIT_OFF	0x00       // Straight Off
#define LAST_EVENT 0x80    //high bit set to indicate last one

// slightly random
#define MAX_BIT_EVENTS (MAX_MESSAGE_SIZE - 22)/5 
#define MAX_NS_PER_CARRIER_EVENT 2^32

typedef struct {
    time32 time_offset;
    BitEvent bit_event;
    char pad[3];
} CarrierBitEvent;

typedef struct {
    char verb = CARRIER_EVENT;
    char pad = 0x00;
    ChannelIdType channel;
    CarrierKeyType carrier_key;
    time64 start_timestamp;
    CarrierBitEvent bit_events[MAX_BIT_EVENTS];
    char pad2;
    time64 send_time;
} CarrierEventPayload;

// Versions
#define RT_POC 0
#define RT_ALPHA 1
#define RT_BETA 2
#define RT_RC 3
#define RT_FINAL 4  

typedef struct {
    char major;
    char minor;
    char patch;
    char rel_type;
} Version;

typedef struct {
    Version my_protocol_version;
    Version my_code_version;
    Version latest_stable_version;
} VersionInfoPayload;

void set_message_sender(void(*s)(char * pkt, int len)); 

void set_handler(unsigned char verb, void (*handler)(void *payload));

void decode_message(uint8_t * data, int length);

void time_sync();

void listen_request(int channel, char * callsign); 
