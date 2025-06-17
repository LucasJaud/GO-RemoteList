package remotelist

import (
    "errors"
)

type List struct {
	// mu   sync.RWMutex
	Data []int
}

func NewList() *List {
	return &List{
		Data: make([]int, 0),
	}
}

func (l *List) Append(val int) error {
	
	l.Data = append(l.Data, val)
	return nil
}

func (l *List) Remove() (int, error) {
	if len(l.Data) == 0 {
		return 0, errors.New("empty list")
	}
	lastItem := l.Data[len(l.Data)-1]
	l.Data = l.Data[:len(l.Data)-1]
	
	return lastItem, nil
}

func (l *List) Get(pos int) (int, error) {
	
	if pos < 0 || pos >= len(l.Data) {
		return 0, errors.New("position out of range")
	}
	
	return l.Data[pos], nil
}

func (l *List) Size() int {
	return len(l.Data)
}

func (l *List) IsEmpty() bool {
    return len(l.Data) == 0
}

