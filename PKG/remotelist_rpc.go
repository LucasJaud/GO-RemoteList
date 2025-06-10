package remotelist

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type RemoteLists struct {
    mu    sync.RWMutex
    lists map[string]*List
}

func NewRemoteLists() *RemoteLists {
	rl := &RemoteLists{
		lists: make(map[string]*List),
	}

	if err := rl.LoadLatestSnapshot(); err != nil {
		fmt.Printf("  Failed Loading Snapshot: %v\n", err)
	} else {
		fmt.Println(" Snapshot loaded sucessfully")
	}

	return rl
}

// func ensureSnapshotDirExists(dir string) error {
// 	return os.MkdirAll(dir, os.ModePerm)
// }

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
	if name == "" {
    *reply = 0 
    return errors.New("invalid name")
}
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	list, ok := rl.lists[name]
	if !ok{
		*reply = 0
		return errors.New("list not found")
	}

	val, err := list.Get(pos)
	if err != nil{
		*reply =0
		return  err
	}
	*reply = val
	return  nil
}

func (rl *RemoteLists) Remove(name string,reply *int) error{
	if name == "" {
    *reply = 0
    return errors.New("invalid name")
}
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
	*reply = ok
	return nil
}

func (rl *RemoteLists) GetListsNames(reply *[]string) error {
    rl.mu.RLock()
    defer rl.mu.RUnlock()
    
    names := make([]string, 0, len(rl.lists))
    for name := range rl.lists {
        names = append(names, name)
    }
    
    *reply = names
    return nil
}

func (rl *RemoteLists) SaveToFile() error {
	basePath := "data/snapshots"

	err := os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		return err
	}

	timestamp := time.Now().Unix()
	fileName := fmt.Sprintf("snapshot-%d.gob", timestamp)
	fullPath := filepath.Join(basePath, fileName)

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(rl.lists)
	if err != nil {
		return err
	}
	fmt.Printf(" [SAVE] Snapshot salvo em %s\n", fullPath)
	return nil
}

func (rl *RemoteLists) LoadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var loaded map[string]*List
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&loaded)
	if err != nil {
		return err
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.lists = loaded
	return nil
}

func (rl *RemoteLists) LoadLatestSnapshot() error {
	basePath := "data/snapshots"
	files, err := os.ReadDir(basePath)
	if err != nil {
		return err
	}

	var snapshots []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".gob" {
			snapshots = append(snapshots, file.Name())
		}
	}

	if len(snapshots) == 0 {
		return fmt.Errorf(" Snapshot not found")
	}

	sort.Strings(snapshots)
	latest := snapshots[len(snapshots)-1]
	fullPath := filepath.Join(basePath, latest)

	return rl.LoadFromFile(fullPath)
}

func (rl *RemoteLists) BeforeShutdown() {
	err := rl.SaveToFile()
	if err != nil {
		log.Printf("[ERROR] Failed to save snapshot: %v\n", err)
	} else {
		log.Println("[SUCCESS] Snapshot saved")
	}
}