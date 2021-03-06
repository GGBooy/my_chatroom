package main

import (
	"context"
	"fmt"
	"github.com/gogf/greuse"
	"log"
	"net"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

var serverAddr = "122.9.77.149:10086"
var localAddr = ""
var remoteAddr = ""

func main() {
	conn, err := greuse.Dial("tcp", "0.0.0.0:0", serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr = conn.LocalAddr().String()
	println("local address: " + localAddr)
	linkStart(conn)
}

func linkStart(conn net.Conn) {
	// read the other address from the server
	buffer := make([]byte, 2048)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Print(err)
		return
	}
	remoteAddr = string(buffer[:n])
	println("the other: " + remoteAddr)

	// create another connection at the same port

	var a int
	_, _ = fmt.Scanf("%d", &a)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if a == 1 {
		_, err = DialTimeout("tcp", localAddr, remoteAddr, time.Duration(30)*time.Millisecond)
		//go tryConn(conn, remoteAddr)
		//ls, err := greuse.Listen("tcp", localAddr)
		//if err != nil {
		//	log.Print(err)
		//	return
		//}
		//co, _ := ls.Accept()
		//go Recv(co, ctx)
		//Send(co, ctx)

		http.HandleFunc("/ws", WebSocketHandleFunc)
		log.Fatal(MyListenAndServe(localAddr, nil))

	} else if a == 2 {
		nla, err := greuse.ResolveAddr("tcp", localAddr)
		if err != nil {
			log.Print(err)
		}
		var MyDialOpt = websocket.DialOptions{
			HTTPClient: &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyFromEnvironment,
					DialContext: (&net.Dialer{
						Timeout:   30 * time.Second,
						KeepAlive: 30 * time.Second,
						// set reuse
						Control:   greuse.Control,
						LocalAddr: nla,
					}).DialContext,
					ForceAttemptHTTP2:     true,
					MaxIdleConns:          100,
					IdleConnTimeout:       90 * time.Second,
					TLSHandshakeTimeout:   10 * time.Second,
					ExpectContinueTimeout: 1 * time.Second,
				},
			},
		}

		for i := 0; i < 10; i++ {
			co, _, e := websocket.Dial(ctx, "ws://"+remoteAddr+"/ws", &MyDialOpt)
			if e == nil {
				defer co.Close(websocket.StatusInternalError, "???????????????")
				go Recv(co, ctx)
				Send(co, ctx)
			}
		}
	}
}

func WebSocketHandleFunc(w http.ResponseWriter, req *http.Request) {
	// Accept ?????????????????? WebSocket ?????????????????????????????? WebSocket???
	// ?????? Origin ?????????????????????Accept ????????????????????????????????? InsecureSkipVerify ?????????????????????????????? AcceptOptions ????????????
	// ?????????????????????????????????????????????????????????????????????????????????Accept ??????????????????????????????
	conn, err := websocket.Accept(w, req, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		log.Println("websocket accept error:", err)
		return
	}
	go Recv(conn, req.Context())
	Send(conn, req.Context())
}

type myHandler struct{}

func (mh myHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	WebSocketHandleFunc(w, req)
}

func MyListenAndServe(addr string, handler http.Handler) error {
	server := &http.Server{Addr: addr, Handler: handler}
	ls, _ := greuse.Listen("tcp", localAddr)
	return server.Serve(ls)
}

// Added by myself
func DialTimeout(network, laddr, raddr string, timeout time.Duration) (net.Conn, error) {
	nla, err := greuse.ResolveAddr(network, laddr)
	if err != nil {
		return nil, err
	}
	d := net.Dialer{
		Timeout:   timeout,
		Control:   greuse.Control,
		LocalAddr: nla,
	}
	return d.Dial(network, raddr)
}

func Send(c *websocket.Conn, ctx context.Context) {
	for {
		err := wsjson.Write(ctx, c, "Hello WebSocket Server")
		if err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)
	}
}

func Recv(c *websocket.Conn, ctx context.Context) {
	for {
		var v interface{}
		err := wsjson.Read(ctx, c, &v)
		if err != nil {
			panic(err)
		}
		fmt.Printf("???????????????????????????%v\n", v)
	}
}
