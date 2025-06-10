package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

// estrutura do cliente
type AppendArgs struct {
	Name string
	Val  int
}

type GetArgs struct {
	Name string
	Pos  int
}

type RemoveArgs struct {
	Name string
}

type SizeArgs struct {
	Name string
}

type ListExistsArgs struct {
	Name string
}

type GetAllArgs struct{
	Name string
}

type GetListsNamesArgs struct{}

type RPCClient struct {
	client *rpc.Client
}

func NewRPCClient(address string) (*RPCClient, error) {
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}
	
	return &RPCClient{client: client}, nil
}

func (c *RPCClient) Close() {
	c.client.Close()
}

// metodos do cliente
func (c *RPCClient) Append(name string, val int) error {
	var reply bool
	err := c.client.Call("RemoteListsService.Append", AppendArgs{Name: name, Val: val}, &reply)
	if err != nil {
		return err
	}
	
	if !reply {
		return fmt.Errorf("append failed")
	}
	
	fmt.Printf(" Successfully appended %d to list '%s'\n", val, name)
	return nil
}

func (c *RPCClient) Get(name string, pos int) error{
	var reply int
	err := c.client.Call("RemoteListsService.Get",GetArgs{Name: name, Pos: pos}, &reply)
	if err != nil{
		return err
	}

	fmt.Printf(" Value at position %d in list '%s': %d\n", pos, name, reply)
	return nil
}

func (c *RPCClient) Remove(name string) error{
	var reply int
	err :=c.client.Call("RemoteListsService.Remove",RemoveArgs{Name: name}, &reply)
	if err != nil{
		return err
	}

	fmt.Printf("Removed value %d from list '%s'\n", reply, name)
	return nil
}

func (c *RPCClient) Size(name string) error {
	var reply int
	err := c.client.Call("RemoteListsService.Size", SizeArgs{Name: name}, &reply)
	if err != nil {
		return err
	}
	
	fmt.Printf(" Size of list '%s': %d\n", name, reply)
	return nil
}

func (c *RPCClient) ListExists(name string) error {
	var reply bool
	err := c.client.Call("RemoteListsService.ListExists", ListExistsArgs{Name: name}, &reply)
	if err != nil {
		return err
	}
	fmt.Printf(" List '%s' exists: %t\n", name, reply)
	return nil
}

func (c *RPCClient) GetListsNames() error {
	var reply []string
	err := c.client.Call("RemoteListsService.GetListsNames", GetListsNamesArgs{}, &reply)
	if err != nil {
		return err
	}
	
	fmt.Println(" Available lists:")
	if len(reply) == 0 {
		fmt.Println("  (no lists found)")
	} else {
		for _, name := range reply {
			fmt.Printf("  - %s\n", name)
		}
	}
	return nil
}

func (c *RPCClient) GetAll(name string) error{
	var reply []int

	err := c.client.Call("RemoteListsService.GetAll", GetAllArgs{Name: name}, &reply)
	if err != nil {
		return err
	}
	fmt.Println("Elements of list:")
	if len(reply) == 0 {
		fmt.Println(" list is empty")
	} else {
		fmt.Printf("List '%s' contents: %v\n", name, reply)
	}
	return nil
}

// Interface
func printHelp() {
	fmt.Println("\n Available commands:")
	fmt.Println("  append <list_name> <value>     - Add value to list")
	fmt.Println("  get <list_name> <position>     - Get value at position")
	fmt.Println("  remove <list_name>             - Remove last value")
	fmt.Println("  size <list_name>               - Get list size")
	fmt.Println("  exists <list_name>             - Check if list exists")
	fmt.Println("  lists                          - Show all list names")
	fmt.Println("  elements <list_name>           - Show elements from list")
	fmt.Println("  help                           - Show this help")
	fmt.Println("  quit                           - Exit client")
	fmt.Println()
}

func main(){

	address := "localhost:8080"
	if len(os.Args) > 1 {
		address = os.Args[1]
	}

	fmt.Printf(" Connecting to RPC server at %s...\n", address)
	client, err := NewRPCClient(address)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()

	fmt.Println(" Connected successfully!")
	printHelp()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\nrpc> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		comand := strings.ToLower(parts[0])

		switch comand {
			case "q","quit","exit":
				fmt.Println(" Goodbye")
				return

			case "help","h":
				printHelp()
			
			case "append":
				if len(parts) != 3 {
					fmt.Println(" Usage: append <list_name> <value>")
					continue
				}
				val, err := strconv.Atoi(parts[2])
				if err != nil {
					fmt.Printf(" Invalid value: %s\n", parts[2])
					continue
				}
				if err := client.Append(parts[1], val); err != nil {
					fmt.Printf(" Error: %v\n", err)
				}

			case "get":
				if len(parts) != 3 {
					fmt.Println(" Usage: get <list_name> <position>")
					continue
				}
				pos, err := strconv.Atoi(parts[2])
				if err != nil {
					fmt.Printf(" Invalid position: %s\n", parts[2])
					continue
				}
				if err := client.Get(parts[1], pos); err != nil {
					fmt.Printf(" Error: %v\n", err)
				}
			
			case "remove":
				if len(parts) != 2 {
					fmt.Println(" Usage: remove <list_name>")
					continue
				}
				if err := client.Remove(parts[1]); err != nil {
					fmt.Printf(" Error: %v\n", err)
				}	
			
			case "size":
				if len(parts) != 2 {
					fmt.Println(" Usage: size <list_name>")
					continue
				}
				if err := client.Size(parts[1]); err != nil {
					fmt.Printf(" Error: %v\n", err)
				}
			
			case "exists":
				if len(parts) != 2 {
					fmt.Println(" Usage: exists <list_name>")
					continue
				}
				if err := client.ListExists(parts[1]); err != nil {
					fmt.Printf(" Error: %v\n", err)
				}
			
			case "lists":
				if err := client.GetListsNames(); err != nil {
					fmt.Printf(" Error: %v\n", err)
				}
			
			case "elements":
				if len(parts) != 2 {
					fmt.Println(" Usage: elements <list_name>")
					continue
				}
				if err := client.GetAll(parts[1]); err != nil {
					fmt.Printf(" Error: %v\n", err)
				}
			
			default:
				fmt.Printf(" Unknown command: %s\n", comand)
				fmt.Println("Type 'help' for available commands")
		}
	}
}