package main

import (
	"net"
	"fmt"
	)

type Proxy struct{
	from *net.Conn
	to *net.Conn
}

func (p Proxy)Forward(msg chan string){
	fmt.Print("got here")
	close(msg)
}


func main(){
	
	fmt.Print("started\n")
	msg1 := make(chan string)
	msg2 := make(chan string)
	
	ln, err := net.Listen("tcp", ":5000")
	if err != nil {
		// handle error
	}
	
	conn1, err := ln.Accept()
	if err != nil {
		// handle error
	}
	
	conn2, err := net.Dial("tcp", ":5001")
	if err != nil {
		// handle error
	}
	
	cToS := Proxy{&conn2, &conn1}
	sToC := Proxy{&conn1, &conn2}
	
	
	go cToS.Forward(msg1)
	go sToC.Forward(msg2)
}
