package main

import (
	"context"
	"fmt"
	"os"

	"github.com/veloek/alarm/daemon"
)

func listAlarms() {
	c := newClient()
	resp, err := c.GetAlarms(context.Background(), &daemon.GetAlarmsRequest{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	printHeader()
	for _, a := range resp.Alarms {
		printAlarm(a)
	}
}

func printHeader() {
	fmt.Println("ID  \tTIME       \tRECURRENCE")
}

func printAlarm(a *daemon.Alarm) {
	var format string
	if a.Time.Format == daemon.Time_AM {
		format = "AM"
	}
	if a.Time.Format == daemon.Time_PM {
		format = "PM"
	}
	fmt.Printf("%-4d\t%02d:%02d:%02d%3s\t%s\n",
		a.Id,
		a.Time.Hours,
		a.Time.Minutes,
		a.Time.Seconds,
		format,
		a.Recurrence)
}
