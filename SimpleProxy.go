//Created by Adrian Veliz
//A simple proxy to forwards tcp traffic between localhost ports 5000 and 5001
package main

import (
	"net"
	"fmt"
	)

type Proxy struct{
	from net.Conn
	to net.Conn
	dir string
}

func (p Proxy) Forward (c chan string){
	var b = make([]byte, 1024)
	for {
		num1, error1 := p.from.Read(b)
		if error1 != nil || num1 == 0 {
			p.closeAll()
			c <- "closed " + p.dir
			return
		}

		num2, error2 := p.to.Write(b[0:num1])
		if error2 != nil || num2 == 0 {
			p.closeAll()
			c <- "closed " + p.dir
			return
		}	
	}
}

func (p Proxy) closeAll(){
	p.from.Close()
	p.to.Close()
}

func main(){
	fmt.Println("Starting up SimpleProxy. Listening on port 5000, forwarding to 5001.")
	fmt.Println("Listening for connection.")
	ln, err := net.Listen("tcp", "127.0.0.1:5000")
	if err != nil {
		fmt.Print(err)
		return //error
	}
	
	conn1, err1 := ln.Accept()
	if err1 != nil {
		fmt.Print(err1)
		return //error
	}
	fmt.Print("Connected. Setting up forwarding connection.")
	conn2, err2 := net.Dial("tcp", "127.0.0.1:5001")
	if err2 != nil {
		fmt.Print(err2)
		return //error
	}
	fmt.Println("Connected.")
	
	cToS := Proxy{conn2, conn1, " --> "}
	sToC := Proxy{conn1, conn2, " <-- "}
	
	c := make(chan string)
	go cToS.Forward(c)
	go sToC.Forward(c)

	fmt.Println("SimplProxy now forwarding.")

	done1, done2 :=  <-c, <-c
	fmt.Println(done1 + done2)
}
