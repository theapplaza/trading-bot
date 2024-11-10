package core

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// HandleErrorsAndShutdown listens for errors or signals, then initiates shutdown.
func HandleErrorsAndShutdown(cancel context.CancelFunc) {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan // block and wait for an OS signal
	log.Printf("Received signal: %v. Initiating shutdown...", sig)

	cancel() //signals all goroutines to stop

	// Wait for all streamers to finish, so that we don't have hanging goroutines/resources running in the background
	//even though our process has exited
	streamersGroup.Wait()
	log.Println("Program has exited gracefully.")
}
