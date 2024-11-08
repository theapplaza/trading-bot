package core

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

var errorChannel = make(chan error, 10)

func HandleErrors() {
    go func() {
        // Listen for shutdown signals in a separate goroutine.
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
        <-sigChan // Wait for an OS signal
        log.Println("Received shutdown signal. Closing error channel.")
        close(errorChannel) // Closing the channel triggers all listeners to stop
    }()

    go func() {
        for err := range errorChannel {
            if err != nil {
                log.Println("Error received:", err)
            }
        }
    }()
}

func GetErrorChannel() chan error {
    return errorChannel
}