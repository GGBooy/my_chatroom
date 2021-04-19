package main

import (
	"fmt"
	"github.com/gogf/greuse"
	"io"
	"log"
	"net"
	"time"
)

var serverAddr = "122.9.77.149:10086"



func main(){
	conn, err := greuse.Dial("tcp","0.0.0.0:0", serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	linkStart(conn)
}

func linkStart(conn net.Conn)  {
	buffer := make([]byte, 2048)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Print(err)
		return
	}
	remoteAddr := string(buffer[:n])
	println("the other: " + remoteAddr)
	var a int
	_, _ = fmt.Scanf("%d", &a)
	if a==1 {
		_, _ = DialTimeout("tcp", conn.LocalAddr().String(), remoteAddr, time.Duration(10)*time.Millisecond)
		//go tryConn(conn, remoteAddr)
		ls, err := greuse.Listen("tcp", conn.LocalAddr().String())
		if err != nil {
			log.Print(err)
			return
		}
		co, _ := ls.Accept()
		go Recv(co)
		Send(co)
	} else if a==2 {
		for i := 0; i < 10; i++ {
			co, err := greuse.Dial("tcp", conn.LocalAddr().String(), remoteAddr)
			if err != nil {
				time.Sleep(time.Duration(1)*time.Second)
				continue
			}
			go Recv(co)
			Send(co)
		}
	}
}

// Added by myself
func DialTimeout(network, laddr, raddr string, timeout time.Duration) (net.Conn, error) {
	nla, err := greuse.ResolveAddr(network, laddr)
	if err != nil {
		return nil, err
	}
	d := net.Dialer{
		Timeout: timeout,
		Control:   greuse.Control,
		LocalAddr: nla,
	}
	return d.Dial(network, raddr)
}


func Send(c net.Conn) {
	for {
		_, _ = c.Write([]byte("hello," + c.RemoteAddr().String()))
		time.Sleep(time.Duration(1)*time.Second)
	}
}

func Recv(c net.Conn) {
	buffer := make([]byte, 2048)
	for {
		n, err := c.Read(buffer)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		fmt.Println(string(buffer[:n]))
		time.Sleep(time.Duration(1)*time.Second)
	}
}
