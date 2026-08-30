package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)
func testreturn(conn net.Conn , k map[string]string){
	responseBody := "aaji mera kadddu hai ye  "
	// println(responseBody)
	bodyLength :=  len(responseBody)
	fmt.Fprintf(conn, "HTTP/1.1 500 OK\r\n"+
		"Content-Type: text/plain; charset=utf-8\r\n"+
		"Content-Length: %d\r\n"+
		"\r\n"+
		"%s", bodyLength, responseBody)
}
func testreturn1(conn net.Conn , id map[string]string){
	var responseBody string
	fmt.Sprintf(responseBody , "%q" , id["id"])
	println("so the name is " ,id["id"] )
	// println(responseBody)
	bodyLength :=  len(responseBody)
	fmt.Fprintf(conn, "HTTP/1.1 500 OK\r\n"+
		"Content-Type: text/plain; charset=utf-8\r\n"+
		"Content-Length: %d\r\n"+
		"\r\n"+
		"%s", bodyLength, responseBody)
}
func main(){
	r := &router{}
	r.insert("GET" , []string{"home"}, testreturn)
	r.insert("GET" , []string{"aniket", ":id"}, testreturn1)

	var wg sync.WaitGroup
	listener , err := net.Listen("tcp" , ":8081")
	if err != nil {
		log.Fatalf("there was a err, %q" , err)
	}
	for{
		conn , err := listener.Accept()
		if err != nil {
			log.Fatalf("there was a err, %q" , err)
			break
		}
		wg.Add(1)
		go function1(conn , &wg , r)
	}
	wg.Wait()
}

func function1(conn net.Conn ,wg *sync.WaitGroup , r *router){
	defer wg.Done()
	defer conn.Close()
	var first bool
	reader := bufio.NewReader(conn)
	var method string 
	var path []string
	for {
		line , err := reader.ReadString('\n')
		if err != nil {
			log.Printf("failed to connect: %v", err)
			break // logs the error without exiting
		}
		if line == "\r\n" || line == "\n"{
			break
		}
		if !first{
			method , path = seprate(line)
			first = !first
		}else{
			words := strings.TrimSpace(line)
			println(words)
		}
		
	}

	HandlerFunc , id :=  r.search(method , path)
	if HandlerFunc != nil {
		newfunc := auth(HandlerFunc)
		new2func := recovery(newfunc)
		new2func(conn , id)
	}else{
		fmt.Fprintf(conn, "HTTP/1.1 404 Not Found\r\nContent-Length: 9\r\n\r\nNot Found")
	}
	
}

func recovery(next HandlerFunc) HandlerFunc{
	return func(conn net.Conn , parmas map[string]string){

		defer func(){
			if r := recover(); r != nil{
				println("recover from panic ")
				return
			}
		}()
		next(conn , parmas)
	}
}
func auth(next HandlerFunc) HandlerFunc {
	return func(conn net.Conn , parmas map[string]string) {
		for range 1000{
			print("")
		}

		next(conn , parmas)
	}
}



func (r *router) search(method string , path []string)  (HandlerFunc , map[string]string) {
	if r.root == nil {
		return nil ,nil
	}
	parma:= make(map[string]string)
	currentNode := r.root
	for i , w := range path {
		if childNode, exist := currentNode.children[w] ; exist{
			currentNode = childNode
		}else if currentNode.parmachild != nil {
			parma[currentNode.parmachild.data] = w
			currentNode = currentNode.parmachild
		}else if currentNode.whildchild != nil {
			array := w[1:] + "/" + strings.Join(path[i+1:] , "/")
			parma[currentNode.whildchild.data] = array
			currentNode = currentNode.whildchild
		}else{
			return nil , nil
		}
	}
	if currentNode.handler == nil{
		return nil,nil
	}
	return  currentNode.handler[method] , parma
	
}

func seprate(str string) (string , []string){
	slices := strings.Split(str," ")
	method := slices[0]
	slice := strings.Split(slices[1] , "/")
	slice = slice[1:]

	return method , slice
}

type HandlerFunc func(net.Conn , map[string]string)

type node struct{
	data string
	children map[string]*node
	parmachild *node
	whildchild *node
	handler map[string]HandlerFunc
}
type router struct{
	root *node
}
func (r *router) insert(method string, path []string ,procidure HandlerFunc){
	
	if r.root ==nil{
		r.root = &node{}
	}
	currentNode := r.root
	for _ , w := range path{
		if strings.HasPrefix(w , ":"){
			if currentNode.parmachild == nil {
				currentNode.parmachild= &node{data : w[1:]}
			}
			currentNode = currentNode.parmachild
		}else if strings.HasPrefix(w , "*"){
			if currentNode.whildchild == nil {
				currentNode.whildchild = &node{data : w[1:]}
				break
			}
			currentNode = currentNode.whildchild

		}else{
			if currentNode.children == nil {
				currentNode.children = make(map[string]*node)
			}
			if _ , exist := currentNode.children[w] ; !exist{
				currentNode.children[w] = &node{data : w}
			}
			currentNode = currentNode.children[w]
		}
	}
	if currentNode.handler == nil {
		currentNode.handler = make(map[string]HandlerFunc)
	}
	currentNode.handler[method] = procidure

}
 


