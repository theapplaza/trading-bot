package core

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Global channels
var (
	errorChannel = make(chan error, 10)
)

// HandleErrorsAndShutdown listens for errors or signals, then initiates shutdown.
func HandleErrorsAndShutdown() {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Printf("Received signal: %v. Initiating shutdown...", sig)

	// Wait for all streamers to finish
	streamersGroup.Wait()
	log.Println("Program has exited gracefully.")
}

func GetErrorChannel() chan error {
	return errorChannel
}
