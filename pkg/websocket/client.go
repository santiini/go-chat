package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
)

/*
	Client 客户端
*/

type Client struct {
	ID   string          // 特定连接的唯一可识别字符串
	Conn *websocket.Conn // 指向 websocket.Conn 的指针
	Pool *Pool           // 指向 Pool 的指针
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

// 持续监听 websocket.conn 中的消息，通过 Pool Broadcast channel 进行广播通知到池中的每个客户端。
func (c *Client) Read() {
	defer func() {
		// 关闭 websocket.conn
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		message := Message{
			Type: messageType,
			Body: string(p),
		}

		// websocket.conn 接受消息后，使用 Pool.Broadcast 通知所有相关 chan 最新 message
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}
