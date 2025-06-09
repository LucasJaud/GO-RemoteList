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

func NewRemoteListsService() *RemoteListsService {
	return &RemoteListsService{
		lists: remotelist.NewRemoteLists(),
	}
}

// funções das Listas Remotas

func (s *RemoteListsService) Append(args AppendArgs, reply *bool) error{
	fmt.Printf("[List] Append: name=%s, val=%d\n",args.Name, args.Val)
	return s.lists.Append(args.Name, args.Val, reply)
}

func (s *RemoteListsService) Get(args GetArgs, reply *int) error{
	fmt.Printf("[List] Get :name=%s, pos=%d\n",args.Name,args.Pos)
	return s.lists.Get(args.Name, args.Pos, reply)
}

func (s *RemoteListsService) Remove(args RemoveArgs,reply *int) error{
	fmt.Printf("[List] Remove: name=%s\n", args.Name)
	return s.lists.Remove(args.Name, reply)
}

func (s *RemoteListsService) Size(args SizeArgs, reply *int) error {
	fmt.Printf("[List] Size: name=%s\n",args.Name)
	return s.lists.Size(args.Name, reply)
}

func (s *RemoteListsService) ListExists(args ListExistsArgs, reply *bool) error {
    fmt.Printf("[List] ListExists: name=%s\n", args.Name)
    return s.lists.ListExists(args.Name, reply)
}

// metodos do server

var (
    serverStartTime = time.Now()
    totalRequests int = 0
)

type StatsService struct {
    *RemoteListsService
}

func (s *StatsService) incrementRequests() {
    totalRequests++
}

func startServer(port string) error{
	service := NewRemoteListsService()

	rpc.Register(service)
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp",":"+port)
	if err != nil{
		return fmt.Errorf("error creating the listener: %v", err)
	}

	fmt.Printf("initalized Server on port: %s",port)
	fmt.Printf("Endpoint: http://localhost:%s",port)
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