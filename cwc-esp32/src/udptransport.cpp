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
#include <AsyncUDP.h>
#include <WiFi.h>
#include "debug.h"
#include "config.h"
#include "messages.h"


AsyncUDP udp;
IPAddress reflector_IP;
int reflector_port;

void udp_sender(char * pkt, int length) {
    udp.writeTo((uint8_t *)pkt, length, reflector_IP, reflector_port);
}

void udp_transport_setup() {

    WiFi.hostByName(get_config("ReflectorHost").c_str(), reflector_IP);
    debug_print("reflector host: " + reflector_IP.toString() + "\n");

    reflector_port = get_config("ReflectorPort").toInt(); 
    int local_port = get_config("LocalPort").toInt();

    debug_print("Listening UDP port " + get_config("LocalPort") + "\n");
    
   udp.listen(local_port);

    set_message_sender(&udp_sender);

    udp.onPacket([](AsyncUDPPacket p) {
        decode_message(p.data(), p.length()); 
    });
}

void udp_transport_run() {
}
