package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/sachinagada/secretsanta"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	srv := secretsanta.NewServer(nil, nil)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.Handle("/santa", srv)

	// TODO: close the server
	go func() {
		// TODO: make the port configurable
		http.ListenAndServe(":3000", nil)
	}()

	<-ctx.Done()
}
