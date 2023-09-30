package main

import (
	//"context"
	"fmt"
	"os"
	"os/signal"
	"poppypaw/state"
	"syscall"

	"github.com/getsentry/sentry-go"
	"net/http"

	"poppypaw/server"
	//novu "github.com/novuhq/go-novu/lib"
)

func main() {
	// Sentry
	sentry.Init(sentry.ClientOptions{
		Dsn:              "https://0af9ed73773d95692388f88fb8b68777@trace.select-list.xyz/5",
		TracesSampleRate: 1.0,
	})

	state.Setup()

	/* Novu (send basic notification)
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
	}*/

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
