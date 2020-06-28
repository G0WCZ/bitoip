
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
typedef long long time64;
typedef long int time32;
typedef char BitEvent;

// Messages
#define ENUMBERATE_CHANNELS 0x90
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

typedef struct {
   ChannelIdType channels[MAX_CHANNELS_PER_MESSAGE];
} ListChannelsPayload;

typedef struct {
	time64 given_time;
    time64 server_rx_time;
    time64 server_tx_time;
} TimeSyncResponsePayload; 

typedef struct {
    ChannelIdType channel;
    CallsignType callsign;
} ListenRequestPayload;

typedef struct {
    ChannelIdType channel;
    CarrierKeyType carrier_key;
} ListenConfirmPayload;

typedef struct {
    ChannelIdType channel;
    CarrierKeyType carrier_key;
} UnlistenPayload;

typedef struct {
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
#define MAX_NS_PER_CARRIER_EVENT = 2^32

typedef struct {
    time32 time_offset;
    BitEvent bit_event;
} CarrierBitEvent;

typedef struct {
    ChannelIdType channel;
    CarrierKeyType carrier_key;
    time64 start_timestamp;
    CarrierBitEvent bit_events[MAX_BIT_EVENTS];
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