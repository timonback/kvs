package server

import (
	"context"
	"fmt"
	"github.com/timonback/keyvaluestore/internal"
	"github.com/timonback/keyvaluestore/internal/arguments"
	context2 "github.com/timonback/keyvaluestore/internal/server/context"
	"github.com/timonback/keyvaluestore/internal/server/filter"
	"github.com/timonback/keyvaluestore/internal/server/handler"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func StartServer(arguments *arguments.Server) {
	internal.Logger.Println("Server is starting...")

	router := http.NewServeMux()
	router.Handle("/healthz", handler.Healthz())
	router.Handle("/hello", handler.Index())
	router.Handle(context2.HandlerPathStore, handler.Store(arguments.Store))
	router.Handle(context2.HandlerPathInternalId, handler.InternalId())
	router.Handle("/ui/", http.StripPrefix("/", http.FileServer(http.Dir("static"))))

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	server := &http.Server{
		Addr:         ":" + arguments.ListenPort,
		Handler:      filter.Tracing(nextRequestID)(filter.Logging(internal.Logger)(router)),
		ErrorLog:     internal.Logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, os.Kill)
	go func() {
		<-arguments.Stop
		internal.Logger.Println("Received shutdown command...")
		quit <- os.Interrupt
		close(arguments.Stop)
	}()

	go func() {
		<-quit
		internal.Logger.Println("Server is shutting down...")
		handler.SetHealthy(handler.NON_HEALTHY)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			internal.Logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	StartServerDiscovery(arguments)

	internal.Logger.Println("Server is ready to handle requests at", arguments.ListenPort)
	handler.SetHealthy(handler.HEALTHY)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		internal.Logger.Fatalf("Could not listen on %s: %v\n", arguments.ListenPort, err)
	}

	<-done
	internal.Logger.Println("Server stopped")
}
