package main

import (
	"context"
	"fmt"
	"os"

	"github.com/veloek/alarm/daemon"
)

func removeAlarm(id int) {
	c := newClient()
	_, err := c.RemoveAlarm(context.Background(), &daemon.RemoveAlarmRequest{
		AlarmId: int32(id)})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
}
