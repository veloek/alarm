package daemon

import (
	"context"
	"fmt"
	"os"
)

const (
	// PORT is the TCP port the gRPC server is listening on.
	PORT = "127.0.0.1:52543"
)

// Start is a blocking function that will create a database connection,
// run the scheduler in the background and start the gRPC server.
func Start() {
	// Create database connection.
	repo, err := dbConnect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
		return
	}

	updates := make(chan []*Alarm)
	server := newServer(repo, updates)

	// Run scheduler in the background.
	sched := newScheduler(updates, onAlarm(server))
	go sched.start()

	server.start()
}

func onAlarm(s *server) func(a *Alarm) {
	return func(a *Alarm) {
		if a.Recurrence == Alarm_NO_RECURRENCE {
			_, err := s.RemoveAlarm(context.Background(), &RemoveAlarmRequest{AlarmId: a.Id})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error cleaning up one-time alarm: %v\n", err)
			}
		}

		if err := playSound(); err != nil {
			fmt.Fprintf(os.Stderr, "Error playing sound: %v\n", err)
		}

	}
}
