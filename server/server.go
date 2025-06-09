package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"project-RPA/PKG"
	"sync/atomic"
	"syscall"
	"time"
)

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

type RemoteListsService struct{
	lists *remotelist.RemoteLists
}

type GetListsNamesArgs struct{}

var (
	serverStartTime = time.Now()
	totalRequests   int64 = 0
)

func NewRemoteListsService() *RemoteListsService {
	return &RemoteListsService{
		lists: remotelist.NewRemoteLists(),
	}
}

func (s *RemoteListsService) incrementRequests() {
	atomic.AddInt64(&totalRequests, 1)
}

//helpers for log Error and Success
func (s *RemoteListsService) LogError(method, name string, err error){
log.Printf("[ERROR] %s falied for list '%s': %v ",method, name, err)
}

func (s *RemoteListsService) LogSuccess(method, name string, details string){
	log.Printf("[SUCCESS] %s for list '%s': %s",method, name, details)
}

// funções das Listas Remotas

func (s *RemoteListsService) Append(args AppendArgs, reply *bool) error{
	s.incrementRequests()
	fmt.Printf("[List] Append: name=%s, val=%d\n",args.Name, args.Val)

	err := s.lists.Append(args.Name, args.Val, reply)
	if err != nil {
		s.LogError("Append", args.Name, err)
		*reply = false
		return fmt.Errorf("failed to append to list '%s': %w", args.Name, err)
	}
	s.LogSuccess("Append", args.Name, fmt.Sprintf("value=%d", args.Val))
	return nil
}

func (s *RemoteListsService) Get(args GetArgs, reply *int) error{
	s.incrementRequests()
	fmt.Printf("[List] Get :name=%s, pos=%d\n",args.Name,args.Pos)

	err := s.lists.Get(args.Name, args.Pos, reply)
	if err != nil {
		s.LogError("Get", args.Name, err)
		*reply = 0
		return fmt.Errorf("failed to get from list '%s' at position %d: %w", args.Name, args.Pos, err)
	}
	s.LogSuccess("Get", args.Name, fmt.Sprintf("pos=%d, value=%d", args.Pos, *reply))
	return nil
}

func (s *RemoteListsService) Remove(args RemoveArgs,reply *int) error{
	s.incrementRequests()
	fmt.Printf("[List] Remove: name=%s\n", args.Name)
	err := s.lists.Remove(args.Name, reply)
	if err != nil {
		s.LogError("Remove", args.Name, err)
		*reply = 0
		return fmt.Errorf("failed to remove from list '%s': %w", args.Name, err)
	}
	s.LogSuccess("Remove", args.Name, fmt.Sprintf("removed_value=%d", *reply))
	return nil
}

func (s *RemoteListsService) Size(args SizeArgs, reply *int) error {
	s.incrementRequests()
	fmt.Printf("[List] Size: name=%s\n",args.Name)

	err := s.lists.Size(args.Name, reply)
	if err != nil {
		s.LogError("Size", args.Name, err)
		*reply = 0
		return fmt.Errorf("failed to get size of list '%s': %w", args.Name, err)
	}
	s.LogSuccess("Size", args.Name, fmt.Sprintf("size=%d", *reply))
	return nil
}

func (s *RemoteListsService) ListExists(args ListExistsArgs, reply *bool) error {
	s.incrementRequests()
    fmt.Printf("[List] ListExists: name=%s\n", args.Name)

	err := s.lists.ListExists(args.Name, reply)
	if err != nil {
		s.LogError("ListExists", args.Name, err)
		*reply = false
		return fmt.Errorf("failed to check if list '%s' exists: %w", args.Name, err)
	}
	s.LogSuccess("ListExists", args.Name, fmt.Sprintf("exists=%t", *reply))
    return nil
}

func (s *RemoteListsService) GetListsNames(args GetListsNamesArgs, reply *[]string) error {
	s.incrementRequests()
	fmt.Println("[List] GetListsNames")

	err := s.lists.GetListsNames(reply)
	if err != nil {
		log.Printf("[ERROR] GetListsNames failed: %v", err)
		*reply = []string{}
		return fmt.Errorf("failed to get list names: %w", err)
	}
	log.Printf("[SUCCESS] GetListsNames returned %d lists", len(*reply))
	return nil
}

// metodos do server
func startServer(port string) error{
	service := NewRemoteListsService()

	err := rpc.Register(service)
	if err != nil {
		return fmt.Errorf("error registering service: %v", err)
	}
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp",":"+port)
	if err != nil{
		return fmt.Errorf("error creating the listener: %v", err)
	}

	fmt.Printf("initalized Server on port: %s\n",port)
	fmt.Printf("Endpoint: http://localhost:%s\n",port)
	fmt.Printf("Server init at %s\n",serverStartTime.Format("15:04:05"))
	fmt.Println("List Methods:")
	fmt.Println(" - RemoteListService.Append")
	fmt.Println(" - RemoteListService.Get")
	fmt.Println(" - RemoteListService.Remove")
	fmt.Println(" - RemoteListService.Size")
	fmt.Println(" - RemoteListService.ListExists")

	return http.Serve(listener, nil)
}

func SetupServerShutdown(){
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func(){
		<-c
		fmt.Println("interrupt signal received")
		fmt.Printf("Total request of Server : %d\n", totalRequests)
		fmt.Printf("Time running: %s=\n",time.Since(serverStartTime))
		fmt.Printf("terminating server...",)
		os.Exit(0)
	}()
}


func main(){
	port := "8080"
	if p :=os.Getenv("RPC_PORT"); p !="" {
		port = p
	}

	SetupServerShutdown()

	log.Fatal(startServer(port))
}