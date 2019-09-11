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
	fmt.Println("Server starting on port: 80")

	fs := http.FileServer(http.Dir("../frontend/build/static/css"))
	http.Handle("/static/css/", http.StripPrefix("/static/css/", fs))

	fs2 := http.FileServer(http.Dir("../frontend/build/static/js"))
	http.Handle("/static/js/", http.StripPrefix("/static/js/", fs2))

	pool := WS.NewPool()
	go pool.Start()
	router.HandleFunc("/ws",func(w http.ResponseWriter, r *http.Request){
		serveWs(pool,w,r)
	})

	router.HandleFunc("/",func(w http.ResponseWriter, r *http.Request){
		http.ServeFile(w,r, "../frontend/build/index.html")
	})

	http.Handle("/",router)
	http.ListenAndServe(":80",nil)
}
