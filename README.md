# CWC - The CW Commuter

## Status

**10 June 2019:**
Ongoing development of hardware and software in progress. There is now alpha support for a softeare implementation of an iambic keyer using
two GPIO pins.

Coming up next: 
* web-based (and therefore phone-based) configuration for the cwc-station on raspberry pi
* easy software updating
* morse encode and decode
* easy robot support

## The Idea

The idea:  A little box that you can plug a key and headphones into.  It Wifi connects to your phone 
hotspot.  It has a channel selector.   Dial up a channel and tx/rx CW on that channel. That's it.    

This is an internet transceiver for CW that you can take with you.  It aims to be more like a radio than a computer. 

## What is in the box
A Raspberry Pi with WiFi.  A few components for audio out and key in.
There may be a channel knob one day and a signal LED.

## Get Started
It is a bit basic for now, but is simple enough and it works:
1. Raspberry Pi:  Get a PI zero W or a PI 3.  Get it on wifi.  
1. Get latest CWC station binary here: https://github.com/G0WCZ/cwc/releases
1. See the wiki for hardware setup: https://github.com/G0WCZ/cwc/wiki/Pi-Zero---Pi-hardware-setup
log into your Raspberry PI.
1. Run a command where you put that binary:
```
./cwc-station -de <your-call> -sidetone <freq>

# example commands:
# get started with you call as Q1AAA and a 500hz sidetone
./cwc-station -de Q1AAA -sidetone 500

# add verbose log outout
./cwc-station -de Q1AAA -sidetone 500 -logtostderr -v 2

# change PI pins:
./cwc-station -de Q1AAA -sidetone <freq>

# get help on commands
./cwc-station -h
```

This is the simple way to get started with Raspberry PI using GPIO pins.
You can also use a USB serial port on desktops etc.  See the wiki for more details.


## Communications

There's a protocol based on UDP packets that sends on and off events.
So if you use the key, you are sending on and off events in UDP packets.

At the receiving end there's something that turns packetised on-offs back into contact closures or a tone in your ears.  

UDP is lossy, so it is more radio-like in that sense.    You might lose some packets,  some QSB shrug.

# Broadcast or Reflector
There are two basic modes.  Your CWC station can broadcast on the local network, or talk to a reflector.

In broadcast mode, UDP multicast is used on the local network.  This is a simplified mode for co-located CW training
or similar.

In reflector mode, the station connects to a central reflector that reflects traffic to other connected stations.

See bitoip.md for the on-the-wire protocol details.

# Terms of Use
See TERMS_OF_USE.md for details about how you callsign, activity and transmissions are used
by the cwc components.

# Implementations

* early release, still in development: Raspberry Pi GPIO / or Mac & Linux * maybe windows with serial port

# Pi Zero default setup

See the wiki for details of the basic hardware and configuration.

# Developing
You'll need a version of go that understands modules, so version 1.11 or later is required.
If you want to build distributablebinaries, look at the targets in the go/Makefile.

# Who did this
Concept by Grae G0WCZ, Andy M0VVA and The Online Radio Club MX0ONL

Go implementation (for RPi and others) by Grae G0WCZ and Andy M0VVA
