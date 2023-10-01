package main

import (
	"fmt"
	"os"
	"os/signal"
	"poppypaw/state"
	"syscall"

	"net/http"

	"github.com/getsentry/sentry-go"

	"poppypaw/server"
)

func main() {
	// Sentry
	sentry.Init(sentry.ClientOptions{
		Dsn:              "https://0af9ed73773d95692388f88fb8b68777@trace.select-list.xyz/5",
		TracesSampleRate: 1.0,
	})

	state.Setup()

	// Server
	go func() {
		mux := server.Initialize()
		fmt.Println("Now listening on http://localhost:" + state.Config.Port + "/")
		http.ListenAndServe(":"+state.Config.Port, mux)
	}()

	// Terminate
	done := make(chan struct{})

	go func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		s := <-sc

		fmt.Println("Recieved", s, "signal. Closing...")
		state.Close()
		fmt.Println("Closed open state connections. Exiting...")

		close(done)
	}()

	<-done
	fmt.Println("Recieved done channel close. Exiting...")
}
