package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/veloek/alarm/daemon"
)

const (
	runAsDaemon = "startdaemon"
)

func printUsageAndExit() {
	fmt.Println("usage: alarm [-s <hh:mm[:ss]>[AM|PM] [-r <hourly|daily|weekly>] | -l | -rm <id>]")
	os.Exit(1)
}

func main() {
	// Start background process if receiving special argument.
	if len(os.Args) == 2 && os.Args[1] == runAsDaemon {
		err := daemon.Start()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not start daemon: %v\n", err)
		}
		os.Exit(1)
	}

	list := flag.Bool("l", false, "list active alarms")
	remove := flag.Int("rm", 0, "remove alarm by id")
	set := flag.String("s", "", "set alarm")
	recurrence := flag.String("r", "", "recurrence")
	help := flag.Bool("h", false, "print help")
	flag.Parse()

	if *help || (!*list && *remove == 0 && *set == "") {
		printUsageAndExit()
	}

	if !daemonRunning() {
		fmt.Println("Daemon not running. Starting...")
		startDaemon()
		time.Sleep(time.Second) // Wait for background process to start.
	}

	if *list {
		listAlarms()
	} else if *remove > 0 {
		removeAlarm(*remove)
	} else if *set != "" {
		setAlarm(*set, *recurrence)
	}
}

func daemonRunning() bool {
	conn, _ := net.Dial("tcp", daemon.PORT)
	if conn != nil {
		conn.Close()
		return true
	}
	return false
}

func startDaemon() {
	cmd := exec.Command(os.Args[0], runAsDaemon)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
}
