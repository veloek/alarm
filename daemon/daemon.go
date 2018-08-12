package daemon

const (
	PORT = ":52543"
)

// Start is a blocking function that starts the gRPC server.
func Start() error {
	repo, err := dbConnect()
	if err != nil {
		return err
	}

	return startServer(repo)
}
