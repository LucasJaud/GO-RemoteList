package remotelist

import (
    "errors"
    "fmt"
    "sync"
)

type List struct {
	mu   sync.Mutex
	data []int
	size uint32
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
	l.size++
	return nil
}

func (l *List) Remove() (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.data) > 0 {
		lastItem := l.data[len(l.data)-1]
    	l.data = l.data[:len(l.data)-1]
    	l.size--
	} else {
		return 0, errors.New("empty list")
	}
	return lastItem, nil
}

func (l *List) Get(pos int) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	if pos < 0 || pos >= len(l.data) {
		return 0, errors.New("position out of range")
	}
	
	return l.data[pos], nil
}

func (l *List) Size() uint32 {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	return l.size
}

// func (l *List) Data() []int {
// 	l.mu.Lock()
// 	defer l.mu.Unlock()
	
// 	dataCopy := make([]int, len(l.data))
// 	copy(dataCopy, l.data)
// 	return dataCopy
// }