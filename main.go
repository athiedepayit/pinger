package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"
)

type Config struct {
	Interval   int
	ListenPort int
	RemoteHost string
	ErrCommand string
	RecCommand string
}

func (cfg Config) Print() {
	// Yes I KNOW multiline strings exist, but they're UGLY
	fmt.Printf("---\n")
	fmt.Printf("- Interval: %v\n", cfg.Interval)
	fmt.Printf("- ListenPort: %v\n", cfg.ListenPort)
	fmt.Printf("- RemoteHost: %v\n", cfg.RemoteHost)
	fmt.Printf("- ErrCommand: %v\n", cfg.ErrCommand)
	fmt.Printf("- RecCommand: %v\n", cfg.RecCommand)
	fmt.Printf("---\n")
}

func main() {

	// defaults
	defaultInterval := 60
	defaultRemote := "http://localhost:8761"
	defaultPort := 8761

	interval := flag.Int("interval", defaultInterval, "seconds between checks")
	listenPort := flag.Int("port", defaultPort, "local port to listen")
	remoteHost := flag.String("remote", defaultRemote, "remote host to check")

	errCommand := flag.String("errorcmd", "", "command to run when check errors")
	recCommand := flag.String("recoverycmd", "", "command to run when check recovers")

	flag.Parse()

	// temp default before reading from a config file
	config := Config{
		Interval:   *interval,
		ListenPort: *listenPort,
		RemoteHost: *remoteHost,
		RecCommand: *recCommand,
		ErrCommand: *errCommand,
	}

	config.Print()

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

func ExecCommand(command string) error {
	fmt.Printf("sending command '%s'\n", command)
	execCommand := exec.Command("sh", "-c", command)
	err := execCommand.Run()
	if err != nil {
		return err
	}
	return nil
}

func CheckLoop(cfg *Config, checkResult *bool) {
	for {
		// Sleeping first, to make sure that you have time to get both
		// sides running. Don't want to bring one side up, then
		// immediately have it send both an error and a recovery alert
		// before you have enough time to start the other side.
		time.Sleep(time.Second * time.Duration(cfg.Interval))
		localResult, err := RunChecks(cfg)
		if err != nil {
			fmt.Printf("Error running check: %s\n", err)
		}
		fmt.Printf("checking...\nresult is %t\n", localResult)

		// if the local result is true, and previously it was false, send a recovery command
		if !*checkResult && localResult {
			if cfg.RecCommand != "" {
				ExecCommand(cfg.RecCommand)
			}
		}
		// and the opposite, if the command was true, and now it's false, send an error command
		if *checkResult && !localResult {
			if cfg.ErrCommand != "" {
				ExecCommand(cfg.ErrCommand)
			}
		}

		*checkResult = localResult
	}
}

func WebServer(config *Config, checkResult *bool) {

	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		if *checkResult {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, "OK\n")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "Error\n")
		}
	}

	promHandler := func(w http.ResponseWriter, r *http.Request) {
		if *checkResult {
			io.WriteString(w, "pinger_health 1\n")
		} else {
			io.WriteString(w, "pinger_health 0\n")
		}
	}

	http.HandleFunc("/health", handleFunc)
	http.HandleFunc("/metrics", promHandler)

	portStr := fmt.Sprintf(":%v", config.ListenPort)
	err := http.ListenAndServe(fmt.Sprintf("%v", portStr), nil)
	if err != nil {
		fmt.Printf("error serving on port: %s\n", err)
	}
	fmt.Println("Server exiting!")
}
