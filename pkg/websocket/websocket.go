package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

// 我们需要定义一个 Upgrader
// 它需要定义 ReadBufferSize 和 WriteBufferSize
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// 可以用来检查连接的来源
	// 这将允许从我们的 React 服务向这里发出请求。
	// 现在，我们可以不需要检查并运行任何连接
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 将连接更新为 WebSocket 连接
func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
		return ws, err
	}
	return ws, nil
}
