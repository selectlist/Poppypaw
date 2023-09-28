package main

import (
	"context"
	"poppypaw/state"

	"github.com/getsentry/sentry-go"
	novu "github.com/novuhq/go-novu/lib"
)

func main() {
	// Sentry
	sentry.Init(sentry.ClientOptions{
		Dsn:              "https://0af9ed73773d95692388f88fb8b68777@trace.select-list.xyz/5",
		TracesSampleRate: 1.0,
	})

	state.Setup()

	ctx := context.Background()
	to := map[string]interface{}{
		"subscriberId": "564164277251080208",
	}

	payload := map[string]interface{}{
		"message": "This is a notification from Poppypaw!",
	}

	data := novu.ITriggerPayloadOptions{To: to, Payload: payload}
	_, err := state.Novu.EventApi.Trigger(ctx, "primary", data)
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	Message("reee") // cum
}
