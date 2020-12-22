package server

import (
	"context"
	"fmt"
	"github.com/timonback/keyvaluestore/internal"
	"github.com/timonback/keyvaluestore/internal/cli"
	"github.com/timonback/keyvaluestore/internal/server/filter"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

var (
	healthy int32
)

func StartServer(arguments *cli.Arguments) {
	internal.Logger.Println("Server is starting...")

	router := http.NewServeMux()
	router.Handle("/healthz", healthz())
	router.Handle("/hello", index())
	router.Handle("/debug", debug())
	router.Handle("/ui/", http.StripPrefix("/", http.FileServer(http.Dir("static"))))

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	server := &http.Server{
		Addr:         arguments.ListenAddr,
		Handler:      filter.Tracing(nextRequestID)(filter.Logging(internal.Logger)(router)),
		ErrorLog:     internal.Logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		internal.Logger.Println("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			internal.Logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	internal.Logger.Println("Server is ready to handle requests at", arguments.ListenAddr)
	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		internal.Logger.Fatalf("Could not listen on %s: %v\n", arguments.ListenAddr, err)
	}

	<-done
	internal.Logger.Println("Server stopped")
}

func debug() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-INFO", r.RequestURI)
		fmt.Fprintln(w, r.RequestURI)
	})
}

func index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hello, World!")
	})
}

func healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&healthy) == 1 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})
}
