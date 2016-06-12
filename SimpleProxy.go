//Created by Adrian Veliz
//A simple proxy to forwards tcp traffic between localhost ports 5000 and 5001
package main

import (
	"net"
	"fmt"
	"io"
	"strconv"
	)

type Proxy struct{
	from net.Conn
	to net.Conn
	dir string
}

func (p Proxy) Forward (c chan string){
	var b = make([]byte, 1024)
	for {
		num, err := p.from.Read(b)
		
		if num > 0 {//read sucessful
			
			num2, error2 := p.to.Write(b[0:num])
			if num2 > 0 {//assumes everything was written, if the numbers don't match something went wrong
				fmt.Println(p.dir + " " + strconv.Itoa(num) + " bytes rec'd, " + strconv.Itoa(num2) + " bytes sent")
			} else { //handle write error
				c <- p.doError(error2)
				return
			}
		} else { //handle read error
			c <- p.doError(err)
			return
		}
	}
}

func (p Proxy) doError(err error) (msg string){
	//p.from.Close()
	p.to.Close()
	switch err{
		case io.EOF:
			return "EOF, closing connection" + p.dir + "\n"
		default:
			return "closing connection " + p.dir + "\n"
	}
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
	fmt.Println("Connected. \nSetting up forwarding connection.")
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
