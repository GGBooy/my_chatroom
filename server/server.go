package main

import (
	"log"
	"net"
)

func main(){
	listener, err := net.Listen("tcp", "0.0.0.0:10086")
	if err != nil {
		log.Fatal(err)
	}

	conn1, err := listener.Accept()
	if err != nil {
		log.Print(err)
	}
	println("client1")

	//buffer := make([]byte, 2048)
	//n1, err := conn1.Read(buffer)
	//if err != nil {
	//	log.Print(err)
	//	return
	//}
	//remoteAddr1 := string(buffer[:n1])


	conn2, err := listener.Accept()
	if err != nil {
		log.Print(err)
	}
	println("client2")

	//n2, err := conn2.Read(buffer)
	//if err != nil {
	//	log.Print(err)
	//	return
	//}
	//remoteAddr2 := string(buffer[:n2])

	handleConn(conn1, conn2)
	listener.Accept()
}

func handleConn(c1, c2 net.Conn) {
	defer c1.Close()
	defer c2.Close()
	c1.Write([]byte(c2.RemoteAddr().String()))
	println("send to client1: " + c2.RemoteAddr().String())
	c2.Write([]byte(c1.RemoteAddr().String()))
	println("send to client2: " + c1.RemoteAddr().String())

}
