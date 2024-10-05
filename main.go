package main

import (
	"flag"
	"fmt"
	"net"
)

type Config struct {
	Interval   int
	ListenPort int
	RemoteHost string
}

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

	// print config
	fmt.Printf("---\nconfig:\n- Interval %v\n- ListenPort %v\n- RemoteHost %s\n---\n", config.Interval, config.ListenPort, config.RemoteHost)

	RunChecks(&config)
}

func RunChecks(cfg *Config) (bool, error) {
	d := net.Dialer{Timeout: 3}
	conn, err := d.Dial("tcp", cfg.RemoteHost)
	if err != nil {
		return false, err
	}
	conn.Close()
	return true, nil
}
