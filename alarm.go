package main

import (
	"fmt"
	"os"

	"github.com/veloek/alarm/daemon"
	"google.golang.org/grpc"
)

func newClient() daemon.AlarmServiceClient {
	conn, err := grpc.Dial(":52543", grpc.WithInsecure())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while connecting to alarm service: %v", err)
		os.Exit(2)
	}
	return daemon.NewAlarmServiceClient(conn)
}
