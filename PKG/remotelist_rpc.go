package remotelist

import (
    "errors"
    "fmt"
    "sync"
)

type RemoteLists struct {
    mu    sync.Mutex
    lists map[string]*List
}

func NewRemoteLists() *RemoteLists {
	return &RemoteLists{
		lists: make(map[string]*List),
	}
}

func (rl *RemoteLists) Append(name string,val int, reply *bool) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	if name == "" || name == nil {
        return errors.New("Invalid name")
    }
	
	list, ok := rl.lists[name]
	if !ok {
		list = NewList()
		rl.lists[name] = list
	}

    rl.lists[name].Append(val)
    *reply = true
	return nil
}



// func (rl *RemoteLists) Append(name string, val int,reply *bool) error{
// 	rl.mu.Lock()
// 	defer rl.mu.Unlock()

// 	if _, ok := rl.lists[name]; !ok {
//         rl.lists[name] = []int{}
//     }
    
//     rl.lists[name] = append(rl.lists[name], val)
//     fmt.Println(rl.lists[name])
//     rl.size++
//     *reply = true
//     return nil 
// }

// func (rl *RemoteLists) Get(name string, pos int,reply *bool) (int, error){
//     rl.mu.Lock()
//     defer rl.mu.Unlock()

//     list, ok := rl.lists[name]
//     if !ok {
//         *reply = false
//         return 0, errors.New("list not found")
//     }

    
//     if pos < 0 || pos >= len(list) {
//         *reply = false
//         return 0, errors.New("position doesn't exist in list")
//     }

//     fmt.Println(list)
//     *reply = true
//     return rl.lists[name][pos], nil
// }

// func (rl *RemoteLists) Remove(name string,reply *int) (int, error){
//     rl.mu.Lock()
//     defer rl.mu.Unlock()
//     list, ok := rl.lists[name]
//     if (!ok || len(list) == 0){
//         return 0, errors.New("empty list or not found")
//     } 
    
//     rl.lists[name] = list[:len(list)-1]
//     return *reply, nil   
// }

// func (rl *RemoteLists) Size()

// func (rl *RemoteLists) GetListsNames() []string{
//     rl.mu.Lock()
//     defer rl.mu.Unlock()

//     names := make([]string, 0, len(rl.lists))
//     for name := range rl.lists {
//         names = append(names, name)
//     }
//     return names
// }

// // retorna o Remote list com o map iniciado
// func NewRemoteLists() *RemoteLists {
//     return &RemoteLists{
//         lists: make(map[string][]int),
//     }
// }