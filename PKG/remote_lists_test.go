package remotelist

import (
	"sync"
	"testing"
)

func TestRemoteListsRaceCondition(t *testing.T) {
	rl := NewRemoteLists()
	const listName = "concurrent"

	var wg sync.WaitGroup
	const numGoroutines = 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			var ok bool
			err := rl.Append(listName, val, &ok)
			if err != nil || !ok {
				t.Errorf("failed to append value %d: %v", val, err)
			}
		}(i)
	}

	wg.Wait()

	var size int
	err := rl.Size(listName, &size)
	if err != nil {
		t.Errorf("failed to get size: %v", err)
	}
	if size != numGoroutines {
		t.Errorf("expected size %d, got %d", numGoroutines, size)
	}
}