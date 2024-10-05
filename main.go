package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
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

	var checkResult bool
	go WebServer(&config, &checkResult)
	CheckLoop(&config, &checkResult)

}

func RunChecks(cfg *Config) (bool, error) {
	fmt.Printf("checking %s\n", cfg.RemoteHost)
	resp, err := http.Get(cfg.RemoteHost)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return true, nil
}

func CheckLoop(cfg *Config, checkResult *bool) {
	for {

		localResult, err := RunChecks(cfg)
		if err != nil {
			fmt.Printf("Error running check: %s\n", err)
		}
		fmt.Printf("checking...\nresult is %t\n", localResult)
		*checkResult = localResult
		time.Sleep(time.Second * time.Duration(cfg.Interval))
	}
}

func WebServer(config *Config, checkResult *bool) {

	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		if *checkResult {
			io.WriteString(w, "OK\n")
		} else {
			io.WriteString(w, "Error\n")
		}
	}

	http.HandleFunc("/health", handleFunc)

	portStr := fmt.Sprintf(":%v", config.ListenPort)
	err := http.ListenAndServe(fmt.Sprintf("%v", portStr), nil)
	if err != nil {
		fmt.Printf("error serving on port: %s\n", err)
	}
	fmt.Println("Server exiting!")
}
