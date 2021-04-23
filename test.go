package main

import (
	"context"
	"fmt"
	"github.com/gogf/greuse"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"net"
	"net/http"
	"time"
)

func main() {
	localAddr := "127.0.0.1:12312"
	nla, _ := greuse.ResolveAddr("tcp", localAddr)
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
		Subprotocols: []string{"echo"},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < 10; i++ {
		co, _, e := websocket.Dial(ctx, "ws://127.0.0.1:46709", &MyDialOpt)
		if e == nil {
			defer co.Close(websocket.StatusInternalError, "内部错误！")
			go recv(co, ctx)
			send(co, ctx)
		}
	}
}

func send(c *websocket.Conn, ctx context.Context) {
	for {
		err := wsjson.Write(ctx, c, "Hello WebSocket Server")
		if err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)
	}
}

func recv(c *websocket.Conn, ctx context.Context) {
	for {
		var v interface{}
		err := wsjson.Read(ctx, c, &v)
		if err != nil {
			panic(err)
		}
		fmt.Printf("接收到服务端响应：%v\n", v)
	}
}
