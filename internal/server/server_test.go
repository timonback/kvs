package server

import (
	"github.com/timonback/keyvaluestore/internal"
	"github.com/timonback/keyvaluestore/internal/arguments"
	store2 "github.com/timonback/keyvaluestore/internal/store"
	"os"
	"testing"
	"time"
)

func TestStartServer(t *testing.T) {
	internal.InitLogger(false)

	store := store2.NewStoreInmemoryService()
	args := arguments.Server{
		ListenAddr: ":9999",
		Stop:       make(chan os.Signal, 1),
	}

	done := make(chan bool, 1)

	go func() {
		time.Sleep(100 * time.Millisecond)
		args.Stop <- os.Kill
	}()
	StartServer(&args, store)
	done <- true

	select {
	case <-done:
		return
	case <-time.After(1 * time.Second):
		t.Fail()
	}
}
