package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
)

func main1(){
	var wg sync.WaitGroup
	listner , err := net.Listen("tcp" , ":8080")
	if err != nil {
		fmt.Printf("there was a err %q" , err )
		log.Fatalf("there was a bug with the files and all ")
	}

	for{
		bytes1 , err := listner.Accept()
		if err != nil {
			fmt.Printf("there was a err %q" , err )
			// log.Fatalf("there was a bug with the files and all ")
			break
		}
		wg.Add(1)
		go read(bytes1 , &wg)
	}
	wg.Wait()
}

func read(conn net.Conn , wg *sync.WaitGroup) {
	defer wg.Done()
	var result strings.Builder
	defer conn.Close()
	var first bool
	// read :+ bufio.Reader(got) -> this was my fault 
	reader  := bufio.NewReader(conn)
		for {
			line , err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("there was a error while reading the rq")
			}
			if strings.TrimSpace(line) == ""{
				break 
			}
			if !first{
				println(line)
				splitlink(line)
				first = !first
			}
			// fmt.Println(line)
			result.WriteString(line)
		}
	// fmt.Println("Done reading headers!")
	responseBody := result.String()
	bodyLength :=  len(responseBody)
	fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\n"+
		"Content-Type: text/plain; charset=utf-8\r\n"+
		"Content-Length: %d\r\n"+
		"\r\n"+
		"%s", bodyLength, responseBody)
}

// type node struct{
// 	childers map[string]*node
// 	paramchile *node
// 	whildcard *node 
// 	handlers map[string] http.HandlerFunc
// }
// type Router struct{
// 	root *node
// }
func (rt *Router) splitlink(line string){
	words := strings.Split(line, " ")
	var requesttype string
	switch words[0]{
		case "GET":
			requesttype = "GET"
		case "PUT":
			requesttype = "PUT"
		case "DEL":
			requesttype = "DEL"
	}
	path := words[1]
	paths := strings.Split(path, "/")
	paths = paths[1:]

	curentNode := 

	// for i, w := range paths{
	// 	println(i,w)
	// } 

}