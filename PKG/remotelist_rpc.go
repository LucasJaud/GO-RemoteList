package remotelist

import (
    "errors"
    "fmt"
    "sync"
)

type RemoteLists struct {
    mu    sync.Mutex
    lists map[string][]int
    limit uint32  // opcional: limite total de elementos por lista, por exemplo
}