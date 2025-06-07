package remotelist

import (
    "errors"
    "sync"
)

type List struct {
	mu   sync.RWMutex
	data []int
}

func NewList() *List {
	return &List{
		data: make([]int, 0),
	}
}

func (l *List) Append(val int) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	l.data = append(l.data, val)
	return nil
}

func (l *List) Remove() (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.data) == 0 {
		return 0, errors.New("empty list")
	}
	lastItem := l.data[len(l.data)-1]
	l.data = l.data[:len(l.data)-1]
	
	return lastItem, nil
}

func (l *List) Get(pos int) (int, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	if pos < 0 || pos >= len(l.data) {
		return 0, errors.New("position out of range")
	}
	
	return l.data[pos], nil
}

func (l *List) Size() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	return len(l.data)
}

func (l *List) IsEmpty() bool {
    l.mu.RLock()
    defer l.mu.RUnlock()
    
    return len(l.data) == 0
}

