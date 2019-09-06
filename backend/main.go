package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	WS "./websocket"
)


func serveWs(pool *WS.Pool, w http.ResponseWriter, r *http.Request){
	fmt.Println("Websocket Endpoint Hit")
	conn,err := WS.Upgrade(w,r)
	if err != nil{
		fmt.Fprint(w, "%++v\n",err)
	}

	client := &WS.Client{
		Conn: conn,
		Pool: pool,
	}
	
	pool.Register <- client
	client.Read()
}

var router = mux.NewRouter()

func main(){
	fmt.Println("Server starting on port: 8080")

	pool := WS.NewPool()
	go pool.Start()
	router.HandleFunc("/ws",func(w http.ResponseWriter, r *http.Request){
		serveWs(pool,w,r)
	})

	http.Handle("/",router)
	http.ListenAndServe(":8080",nil)
}