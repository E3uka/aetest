package main

import (
	"aetest"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	httpAddr = flag.String("http", ":3000", "http listen address")
)

func run() error {
	// Parse input flags.
	flag.Parse()
	ctx := context.Background()
	errChan := make(chan error)

	item_store, discount, order_store := aetest.NewStore()

	// Create a new service that will handle the order API's requests.
	service := aetest.New(item_store, discount, order_store)
	router := aetest.NewOrdersRouter(service)

	// Ignoring logging, TLS and timeouts for simplicity.
	server := &http.Server{
		Addr:    *httpAddr, // read from input flag
		Handler: router,
	}

	// Start server in goroutine with graceful termination on fatal error.
	go func() {
		if err := server.ListenAndServe(); err != nil {
			errChan <- server.Shutdown(ctx)
		}
	}()

	// Listen for user termination signal (CTRL+C), ungracefully terminate upon
	// receiving this signal.
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
		errChan <- fmt.Errorf("%s", <-ch)
	}()

	return <-errChan
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
