package main

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)
func main() {

	if err := run(); err != nil {
		log.Fatal(err)
	}
}


func run() error {
	cancelCtx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(cancelCtx)
	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello world\n")
	})

	g.Go(func() error {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
			return err
		}
		return nil
	})
	g.Go(func() error {
		<-ctx.Done() // wait for stop signal
		return  srv.Close()
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				cancel()
			}
		}
	})
	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}


