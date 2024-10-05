package main

import (
	"flag"
	"fmt"
)

func main() {

	// defaults
	defaultInterval := 60
	defaultRemote := "http://localhost:8761"
	defaultPort := 8761

	interval := flag.Int("interval", defaultInterval, "seconds between checks")
	listenPort := flag.Int("port", defaultPort, "local port to listen")
	remoteHost := flag.String("remote", defaultRemote, "remote host to check")

	flag.Parse()

	// temp default before reading from a config file
	config := Config{
		Interval:   *interval,
		ListenPort: *listenPort,
		RemoteHost: *remoteHost,
	}

	fmt.Printf("config: %s\n", config)

}

type Config struct {
	Interval   int
	ListenPort int
	RemoteHost string
}
