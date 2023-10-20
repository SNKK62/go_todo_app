package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context, l net.Listener) error {
	s := &http.Server{
		// don't use Addr field because of using net.Listener arg
		// Addr: ":18080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	eg, ctx := errgroup.WithContext(ctx)

	//launch HTTP server in another Goroutine
	eg.Go(func() error {
		// replace ListernAndServe with Serve
		if err := s.Serve(l); err != nil &&
			//http.ErrServerClosed is not an error, but is signal for terminating http.Server.Shutdown() normally
			err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	// waiting for end signal from channel
	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}
	//waiting for end of another Goroutine with Go method.
	return eg.Wait()
}

func main() {
	if len(os.Args) != 2 {
		log.Printf("need port number\n")
		os.Exit(1)
	}
	p := os.Args[1]
	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Fatalf("failed to listen port %s: %v", p, err)
	}
	if err := run(context.Background(), l); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}
