package main

import (
	"encoding/json"
	"fmt"
	"github.com/db"
	"github.com/entity"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	db.FindOnd("user")
	http.HandleFunc("/ws", wsEndpoint)
	http.HandleFunc("/register", register) // 注册
	// 登录涉及到token，暂且不写
	//http.HandleFunc("/friends", handleFriends) // 添加好友
	http.ListenAndServe(":8888", nil)
}

// 解决跨域问题
func cros(writer *http.ResponseWriter) {
	(*writer).Header().Set("Access-Control-Allow-Origin", "*")
	(*writer).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
	(*writer).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func reader(conn *websocket.Conn) {
	time.Sleep(time.Second * 3)
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, []byte("Hello from GoLang!")); err != nil {
			log.Println(err)
			return
		}
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	reader(ws)
}

func register(w http.ResponseWriter, r *http.Request) {
	cros(&w)
	switch r.Method {
	case http.MethodPost:
		var user entity.User
		decoder := json.NewDecoder(r.Body)
		_ = decoder.Decode(&user)
		fmt.Println(user)
		err := db.InsertUser(user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("注册失败！错误原因：" + err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
	case http.MethodOptions: // 对于post请求，浏览器首先会发option请求，如果服务器响应完全符合请求要求，浏览器则会发送真正的post请求。
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
