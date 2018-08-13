package daemon

import (
	"fmt"
	"os"
)

const (
	PORT = ":52543"
)

// Start is a blocking function that will setup database connection, start a
// scheduler and setup the gRPC server.
func Start() {
	// Create database connection.
	repo, err := dbConnect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
		return
	}

	// Start scheduler.
	updates := make(chan []*Alarm)
	sched := newScheduler(updates)
	go sched.start()

	startServer(repo, updates)
}
