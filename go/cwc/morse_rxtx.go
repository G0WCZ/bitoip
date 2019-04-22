package cwc

import (
	"context"
	"log"
	"sort"
	"sync"
	"time"
)
import "../bitoip"

/**
 * Morse hardware receiver and sender
 *
 * Takes incoming morse (a bit going high and low) and turns it into
 * CarrierBitEvents to send.
 *
 * Based on a regular tick that samples the input and builds a buffer
 */

const Ms = int64(1e6)
const Us = int64(1000)
const DefaultTickTime = time.Duration(5 * Ms)
const MaxSendTimespan = time.Duration(4000 * Ms)
const BreakinTime = time.Duration(300 * Ms)
const MaxEvents = 100

var TickTime = DefaultTickTime
var SendWait = MaxSendTimespan

var LastBit bool = false

type Event struct {
	startTime time.Time
	bitEvent bitoip.BitEvent
}

var events = make([]Event, 0, MaxEvents)

var ticker *time.Ticker
var done = make(chan bool)

func SetTickTime(tt time.Duration) {
	TickTime = tt
}

func SetSendWait(sw time.Duration) {
	SendWait = sw
}

func RunMorseRx(ctx context.Context, morseIO IO, toSend chan bitoip.CarrierEventPayload) {
	LastBit = false // make sure turned off to begin -- the default state
	ticker = time.NewTicker(TickTime)

	Startup(morseIO)

	for {
		select {
		case <- done:
			ticker.Stop()
			return

		case t := <-ticker.C:
			Sample(t, toSend, morseIO)
		}
	}
}

func Stop(morseIO IO) {
	done <- true
	LastBit = false
	morseIO.Close()
}

func Startup(morseIO IO) {
	err := morseIO.Open()
	if err != nil {
		log.Fatalf("Can't access Morse hardware: %s", err)
	}
}

// Sample input
// TODO should have some sort of back-off if not used recently for power saving
func Sample(t time.Time, toSend chan bitoip.CarrierEventPayload, morseIO IO) {

	TransmitToHardware(t, morseIO)

	rxBit := morseIO.Bit()
	if rxBit != LastBit {
		// change so record it
		LastBit = rxBit

		var bit uint8 = 0

		if rxBit {
			bit = 1
		}
		events = append(events, Event{t, bitoip.BitEvent(bit) })
		if  len(events) >= MaxEvents {
			events = Flush(events, toSend)
			return
		}
	}
	if len(events)> 0 && ((t.Sub(events[0].startTime) >= MaxSendTimespan) ||
		(t.Sub(events[len(events) - 1].startTime) >= BreakinTime)) {
		events = Flush(events, toSend)
	}
}


// Flush events into an output stream
func Flush(events []Event, toSend chan bitoip.CarrierEventPayload) []Event {
	log.Printf("Flushing events %v", events)
	if len(events) > 0 {
		toSend <- BuildPayload(events)
		events = events[:0]
	}
	return events
}

func BuildPayload(events []Event) bitoip.CarrierEventPayload {
	baseTime := events[0].startTime.UnixNano()
	cep := bitoip.CarrierEventPayload{
		0,
		0,
		baseTime,
		[bitoip.MaxBitEvents]bitoip.CarrierBitEvent{},
	}
	for i, event := range events {
		bit := event.bitEvent

		// mark last event this message
		if i == (len(events) - 1) {
			bit = bit | bitoip.LastEvent
		}

		cep.BitEvents[i] = bitoip.CarrierBitEvent{
			uint32(event.startTime.UnixNano() - baseTime),
			bit,
		}
	}
	return cep
}

/**
 * Transmitting morse out a gpio pin
 */


var TxMutex = sync.Mutex{}
var TxQueue = make([]Event, 100)

// Queue this stuff for sending... Basically add to queue
// that will be sent out based on the tick timing
func QueueForTransmit(carrierEvents bitoip.CarrierEventPayload) {
	// compose into events
	newEvents := make([]Event, 1)
	//start := time.Unix(0, carrierEvents.StartTimeStamp)
	now := time.Now()
	for _, ce := range carrierEvents.BitEvents {
		newEvents = append(newEvents, Event{
			now.Add(time.Duration(ce.TimeOffset)),
			ce.BitEvent,
		})
		if (ce.BitEvent & bitoip.LastEvent) > 0 {
			break
		}
	}
	TxMutex.Lock()

	TxQueue = append(TxQueue, newEvents...)

	sort.Slice(TxQueue, func(i, j int) bool {return TxQueue[i].startTime.Before(TxQueue[j].startTime)})

	TxMutex.Unlock()
}


func TransmitToHardware(t time.Time, morseIO IO) {
	now := time.Now()

	TxMutex.Lock()

	if len(TxQueue) > 0 && TxQueue[0].startTime.Before(now) {
		be := TxQueue[0].bitEvent
		morseIO.SetBit(!(be & bitoip.BitOn == 0))
		TxQueue = TxQueue[1:]
	}

	TxMutex.Unlock()
}