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
	internal.InitLogger("test")

	args := arguments.Server{
		ListenPort: 9999,
		Stop:       make(chan os.Signal, 1),
		Store:      store2.NewStoreInmemoryService(""),
	}

	done := make(chan bool, 1)

	go func() {
		time.Sleep(100 * time.Millisecond)
		args.Stop <- os.Kill
	}()
	StartServer(&args)
	done <- true

	select {
	case <-done:
		return
	case <-time.After(1 * time.Second):
		t.Fail()
	}
}
