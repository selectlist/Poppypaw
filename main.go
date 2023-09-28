package main

import (
	"github.com/getsentry/sentry-go"
	"poppypaw/state"
	"time"
)

func main() {
	sentry.Init(sentry.ClientOptions{
		Dsn:              "https://0af9ed73773d95692388f88fb8b68777@trace.select-list.xyz/5",
		TracesSampleRate: 1.0,
	})

	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	sentry.CaptureMessage("It works!")

	state.Setup()

	Message("reee") // cum
}
