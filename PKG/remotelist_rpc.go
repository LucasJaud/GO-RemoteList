package remotelist

import (
    "errors"
    "sync"
)

type RemoteLists struct {
    mu    sync.RWMutex
    lists map[string]*List
}

func NewRemoteLists() *RemoteLists {
	return &RemoteLists{
		lists: make(map[string]*List),
	}
}

func (rl *RemoteLists) Append(name string,val int, reply *bool) error {
	if name == ""{
		*reply = false
        return errors.New("invalid name")
    }
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	list, ok := rl.lists[name]
	if !ok {
		list = NewList()
		rl.lists[name] = list
	}

    err := list.Append(val)
	if err != nil {
		*reply =false
		return err
	}

    *reply = true
	return nil
}

func (rl *RemoteLists) Get(name string,pos int, reply *int) error {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	if rl.lists[name] == nil{
		*reply = 0
		return errors.New("list dont exist ")
	}

	val, err := rl.lists[name].Get(pos)
	if err != nil{
		*reply =0
		return  err
	}
	*reply = val
	return  nil
}

func (rl *RemoteLists) Remove(name string,reply *int) error{
	rl.mu.Lock()
	defer rl.mu.Unlock()

	list , ok := rl.lists[name]
	if !ok {
		*reply = 0
		return errors.New("list not found")
	}
	val ,err := list.Remove()
	if err != nil{
		*reply = 0
		return err
	}
	*reply = val 
	return nil
}

func (rl *RemoteLists) Size(name string, reply *int) error {
	if name == ""{
		*reply = 0
		return errors.New("invalid name")
	}
	rl.mu.RLock()
    defer rl.mu.RUnlock()

	list, ok :=rl.lists[name]
	if !ok{
		*reply = 0
		return errors.New("list not found")
	}

	*reply = list.Size()
	return nil
}

func (rl *RemoteLists) ListExists(name string,reply *bool) error{
	if name == "" {
        *reply = false
		return errors.New("invalid name")
    }
	rl.mu.RLock()
    defer rl.mu.RUnlock()

	_, ok := rl.lists[name]
	if !ok{
		*reply = false
		return errors.New("list not found")
	}
	*reply = true
	return nil
}




