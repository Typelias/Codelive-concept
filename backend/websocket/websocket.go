package websocket

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//"os"
	"os/exec"
	//"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
)

//STRUCTS
type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

var latest Message = Message{Type: 1, Body: ""}

//WS

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool { return true },
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return conn, err
	}

	return conn, err
}

//POOL
func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of connection Pool", len(pool.Clients))
			client.Conn.WriteJSON(latest)
			for client := range pool.Clients {
				fmt.Println(client)
			}
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of connection Pool", len(pool.Clients))
			break
		case message := <-pool.Broadcast:
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}

//CLIENT
func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Message recieved")
		body := strings.Split(string(p), ";:")
		fmt.Println(body)
		if body[0] == "1" {
			message := Message{Type: messageType, Body: body[1]}
			latest = message
			c.Pool.Broadcast <- message
			fmt.Printf("Message Recived: %+v\n", message)
		} else if body[0] == "2" {
			message := Message{Type: messageType, Body: body[1]}
			fmt.Printf("Message Recived: %+v\n", message)
			compile(c, body[1])
		}

	}
}

func compile(c *Client, code string) {
	//os.Mkdir("."+string(filepath.Separator)+"TEST", 0777)
	if err := ioutil.WriteFile("main.cpp", []byte(code), 0777); err != nil {
		panic(err)
	}
	/*make := "main : main.cpp\n\t /usr/bin/g++ -o main main.cpp\n run : main\n\t ./main"
	if err := ioutil.WriteFile("makefile", []byte(make), 0777); err != nil {
		panic(err)
	}

	out, err := exec.Command("make run").Output()
	if err != nil {
		fmt.Println(err)
	}*/

	out, err := exec.Command("g++", "-o main main.cpp").Output()
	if err != nil {
		fmt.Println(err)
	}

	output := string(out[:])
	fmt.Println(output)

	/*out, err = exec.Command("./main").Output()
	if err != nil {
		fmt.Println(err)
	}*/

	//os.RemoveAll("main.cpp")
	//os.RemoveAll("main")
	//os.RemoveAll("makefile")

	output = string(out[:])
	fmt.Println(output)

	//fmt.Println("Compiling code", code)
	message := Message{Type: 2, Body: output}
	c.Pool.Broadcast <- message

}
