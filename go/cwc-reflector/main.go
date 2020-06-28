/*
Copyright (C) 2019 Graeme Sutherland, Nodestone Limited


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
package main

import (
	"context"
	"flag"
	"github.com/G0WCZ/cwc/bitoip"
	"github.com/golang/glog"
	"net"
	"os"
)

func main() {
	address := flag.String("address", "", "-address=host:port")
	configFile := flag.String("config", "config/cwc-reflector.txt", "-config <filename>")

	flag.Parse()

	// read Config file and defaults
	config := ReadConfig(*configFile)

	if len(*address) > 0 {
		config.CWCAddress = *address
	}

	ReflectorServer(context.TODO(), config)
}

func ReflectorServer(ctx context.Context, config *ReflectorConfig) {

	glog.Info(DisplayVersion())

	serverAddress, err := net.ResolveUDPAddr("udp", config.CWCAddress)
	if err != nil {
		glog.Fatalf("Can't use address %s: %s", config.CWCAddress, err)
		os.Exit(1)
	}

	glog.Infof("Starting reflector on %s", config.CWCAddress)

	messages := make(chan bitoip.RxMSG)

	go bitoip.UDPRx(ctx, serverAddress, messages)

	go UpdateChannelActivity(ctx)

	go APIServer(ctx, &channels, config)

	go Supervisor(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case m := <-messages:
			Handler(serverAddress, m)
		}
	}
}
