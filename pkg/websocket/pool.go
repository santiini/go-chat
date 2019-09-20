package websocket

import "fmt"

// 首先定义一个 Pool 结构体，它将包含我们进行并发通信所需的所有 channels，以及一个客户端 map。
type Pool struct {
	Register   chan *Client     // Register:  当新客户端连接时，Register channel 将向此池中的所有客户端发送 New User Joined...
	Unregister chan *Client     // 注销用户，在客户端断开连接时通知池
	Clients    map[*Client]bool // 客户端的布尔值映射。可以使用布尔值来判断客户端活动/非活动
	Broadcast  chan Message     // 一个 channel，当它传递消息时，将遍历池中的所有客户端并通过套接字发送消息。
}

// 返回 Pool 实例的指针
func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

/*
	Start 方法：
	1. 我们需要确保应用程序中只有一个点能够写入 WebSocket 连接，否则将面临并发写入问题。
	2. 定义了 Start() 方法，该方法将一直监听传递给 Pool channels 的内容，
		然后，如果它收到发送给其中一个 channel 的内容，它将采取相应的行动。
*/
func (pool *Pool) Start() {
	for {
		select {
		// client 连接时
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				fmt.Println(client)
				client.Conn.WriteJSON(Message{
					Type: 1,
					Body: "New User Joined ....",
				})
			}
			break
		// Unregister:  注销用户，在客户端断开连接时通知池
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(Message{
					Type: 1,
					Body: "User Disconnected...",
				})
			}
			break
		// 接受 broadcast 的 message
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
