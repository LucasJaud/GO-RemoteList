package main

import (
	"net/rpc"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"
)

type AppendArgs struct {
	Name string
	Val  int
}

type SizeArgs struct {
	Name string
}

func startServer(t *testing.T) {
	cmd := exec.Command("go", "run", "server/server,go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		t.Fatalf("failed to start server: %v", err)
	}

	// Aguarda o servidor iniciar
	time.Sleep(2 * time.Second)

	t.Cleanup(func() {
		cmd.Process.Kill()
	})
}

func TestRemoteList_ConcurrentAccess(t *testing.T) {
	startServer(t)

	client, err := rpc.DialHTTP("tcp", "localhost:8080")
	if err != nil {
		t.Fatalf("failed to connect to RPC server: %v", err)
	}
	defer client.Close()

	var wg sync.WaitGroup
	name := "testList"

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := AppendArgs{Name: name, Val: i}
			var reply bool
			if err := client.Call("RemoteListsService.Append", args, &reply); err != nil || !reply {
				t.Errorf("failed to append %d: %v", i, err)
			}
		}(i)
	}

	wg.Wait()

	var size int
	err = client.Call("RemoteListsService.Size", SizeArgs{Name: name}, &size)
	if err != nil {
		t.Fatalf("failed to get size: %v", err)
	}
	if size != 100 {
		t.Errorf("expected size 100, got %d", size)
	}
}
