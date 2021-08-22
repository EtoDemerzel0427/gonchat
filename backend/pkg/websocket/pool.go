package websocket

import "fmt"

type Pool struct {
	Register   chan *Client       // Send msg to al other clients to register
	Unregister chan *Client       // Unregister a user and notify the pool when exit
	Clients    map[*Client]bool   // active status
	Broadcast  chan Message       // pass a message through all clients in the chatroom
}

func NewPool() *Pool {
	return &Pool{
		Register: make(chan *Client),
		Unregister: make(chan *Client),
		Clients: make(map[*Client]bool),
		Broadcast: make(chan Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of connection pool: ", len(pool.Clients))
			for client := range pool.Clients {
				fmt.Println(client)
				err := client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
				if err != nil {
					return 
				}
			}
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
			}
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
				}
			}


		}
	}
}