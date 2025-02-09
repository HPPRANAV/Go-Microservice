package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"main.go/handlers"
)

func main() {
	log := log.New(os.Stdout, "product-api", log.LstdFlags)
	ph := handlers.NewProduct(log)

	sm := http.NewServeMux()
	sm.Handle("/", ph)
	sm.Handle("/{id}", ph)

	server := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {

		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	log.Println("Recieved Terminate Signal, gracefully shutting down", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(tc)
}
