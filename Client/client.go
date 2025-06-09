package main

import (
    "fmt"
    "log"
    "net/rpc"
    "time"
)

// estrutura do cliente

type AppendArgs struct{
	Name string
	Val int
}

type GetArgs struct{
	Name string
	Pos int
}

type RemoveArgs struct{
	Name string
}

type SizeArgs struct{
	Name string
}

type ListExistsArgs struct {
    Name string
}

type RemoteListsClient struct {
    client *rpc.Client
}

func NewRemoteListsClient(serverAddr string) (*RemoteListsClient, error) {
    client, err := rpc.DialHTTP("tcp", serverAddr)
    if err != nil {
        return nil, fmt.Errorf("erro ao conectar: %v", err)
    }
    
    return &RemoteListsClient{client: client}, nil
}

func (c *RemoteListsClient) Close() error {
    return c.client.Close()
}

// metodos do cliente 

func (c *RemoteListsClient) Append(name string, val int) (bool, error) {
    args := AppendArgs{Name: name, Val: val}
    var reply bool
    
    err := c.client.Call("RemoteListsService.Append", args, &reply)
    return reply, err
}

func (c *RemoteListsClient) Get(name string, pos int) (int, error) {
    args := GetArgs{Name: name, Pos: pos}
    var reply int
    
    err := c.client.Call("RemoteListsService.Get", args, &reply)
    return reply, err
}

func (c *RemoteListsClient) Remove(name string) (int, error) {
    args := RemoveArgs{Name: name}
    var reply int
    
    err := c.client.Call("RemoteListsService.Remove", args, &reply)
    return reply, err
}

func (c *RemoteListsClient) Size(name string) (int, error) {
    args := SizeArgs{Name: name}
    var reply int
    
    err := c.client.Call("RemoteListsService.Size", args, &reply)
    return reply, err
}

func (c *RemoteListsClient) ListExists(name string) (bool, error) {
    args := ListExistsArgs{Name: name}
    var reply bool
    
    err := c.client.Call("RemoteListsService.ListExists", args, &reply)
    return reply, err
}


// testes

func testRemoteLists(client *RemoteListsClient) {
    fmt.Println(" === INICIANDO TESTES DO CLIENTE RPC ===")
    
    listName := "test_list"
    
    // Teste 1: Verificar se lista existe (deve ser false)
    fmt.Println("\n1. Testando ListExists (lista inexistente):")
    exists, err := client.ListExists(listName)
    if err != nil {
        log.Printf(" Erro: %v", err)
    } else {
        fmt.Printf("   Lista '%s' existe: %t \n", listName, exists)
    }

	   // Teste 2: Adicionar elementos
    fmt.Println("\n2. Testando Append:")
    values := []int{10, 20, 30, 40, 50}
    for i, val := range values {
        success, err := client.Append(listName, val)
        if err != nil {
            log.Printf(" Erro ao adicionar %d: %v", val, err)
        } else {
            fmt.Printf("   [%d] Adicionado %d: %t \n", i, val, success)
        }
    }

	// Teste 3: Verificar tamanho
    fmt.Println("\n3. Testando Size:")
    size, err := client.Size(listName)
    if err != nil {
        log.Printf(" Erro: %v", err)
    } else {
        fmt.Printf("   Tamanho da lista: %d \n", size)
    }
    
    // Teste 4: Verificar se lista existe agora (deve ser true)
    fmt.Println("\n4. Testando ListExists (lista existente):")
    exists, err = client.ListExists(listName)
    if err != nil {
        log.Printf(" Erro: %v", err)
    } else {
        fmt.Printf("   Lista '%s' existe: %t \n", listName, exists)
    }

	// Teste 5: Obter elementos
    fmt.Println("\n5. Testando Get:")
    for i := 0; i < len(values); i++ {
        val, err := client.Get(listName, i)
        if err != nil {
            log.Printf(" Erro ao obter posição %d: %v", i, err)
        } else {
            fmt.Printf("   [%d] = %d \n", i, val)
        }
    }
    
    // Teste 6: Remover elementos
    fmt.Println("\n6. Testando Remove:")
    for i := 0; i < 3; i++ {
        val, err := client.Remove(listName)
        if err != nil {
            log.Printf(" Erro ao remover: %v", err)
        } else {
            fmt.Printf("   Removido: %d \n", val)
        }
    }
    
    // Teste 7: Verificar tamanho final
    fmt.Println("\n7. Testando Size após remoções:")
    size, err = client.Size(listName)
    if err != nil {
        log.Printf(" Erro: %v", err)
    } else {
        fmt.Printf("   Tamanho final: %d \n", size)
    }
    
    
	 fmt.Println("\n === TESTES CONCLUÍDOS ===")
}

// Teste de concorrência
func testConcurrency(client *RemoteListsClient) {
    fmt.Println("\n === TESTE DE CONCORRÊNCIA ===")
    
    listName := "concurrent_list"
    numGoroutines := 10
    numOperations := 5
    
    done := make(chan bool, numGoroutines)
    
    for i := 0; i < numGoroutines; i++ {
        go func(id int) {
            for j := 0; j < numOperations; j++ {
                val := id*100 + j
                client.Append(listName, val)
            }
            done <- true
        }(i)
    }
    
    for i := 0; i < numGoroutines; i++ {
        <-done
    }
    
    size, _ := client.Size(listName)
    fmt.Printf("   Operações concorrentes: %d goroutines × %d ops = %d total\n", 
               numGoroutines, numOperations, numGoroutines*numOperations)
    fmt.Printf("   Tamanho final da lista: %d \n", size)
}



func main() {
    serverAddr := "localhost:8080"
    fmt.Printf("🔌 Conectando ao servidor RPC em %s...\n", serverAddr)
    
    client, err := NewRemoteListsClient(serverAddr)
    if err != nil {
        log.Fatalf(" Erro ao conectar: %v", err)
    }
    defer client.Close()
    
    fmt.Println(" Conectado com sucesso!")
    
    time.Sleep(500 * time.Millisecond)
    
    testRemoteLists(client)
    testConcurrency(client)
    
    fmt.Println("\n Cliente finalizado!")
}