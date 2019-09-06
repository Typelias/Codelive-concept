package websocket

import(
	"log"
	"net/http"
	"fmt"

	"github.com/gorilla/websocket"
)

//STRUCTS
type Pool struct{
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

type Client struct{
	ID string
	Conn *websocket.Conn
	Pool *Pool
}

type Message struct{
	Type int `json:"type"`
	Body string `json:"body"`
}

//WS

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool{return true},
}

func Upgrade(w http.ResponseWriter, r *http.Request)(*websocket.Conn, error){
	conn,err := upgrader.Upgrade(w,r,nil)
	if err != nil{
		log.Println(err)
		return conn,err
	}

	return conn,err
}

//POOL
func NewPool() *Pool{
	return &Pool{
		Register: make(chan *Client),
		Unregister: make(chan *Client),
		Clients: make(map[*Client]bool),
		Broadcast: make(chan Message),
	}
}

func (pool *Pool) Start(){
	for {
		select{
		case client := <- pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of connection Pool",len(pool.Clients))
			for client := range pool.Clients{
				fmt.Println(client)
			}
			break
		case client := <- pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of connection Pool",len(pool.Clients))
			break
		case message := <- pool.Broadcast:
			for client := range pool.Clients{
				if err := client.Conn.WriteJSON(message); err !=nil{
					fmt.Println(err)
					return
				}
			}
		}
	}
}

//CLIENT
func (c *Client) Read(){
	defer func(){
		c.Pool.Unregister <- c
		c.Conn.Close()
	 }()
	 
	 for {
		 messageType, p, err := c.Conn.ReadMessage()
		 if err != nil{
			 log.Println(err)
			 return
		 }
		 message := Message{Type: messageType, Body: string(p)}
		 c.Pool.Broadcast <- message
		 fmt.Printf("Message Recived: %+v\n",message)
	 }
}