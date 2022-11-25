package main

import (
	"encoding/json"
	"fmt"
	"github.com/dto"
	"github.com/gorilla/websocket"
	"github.com/method"
	"github.com/tool"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/ws", wsEndpoint)
	http.HandleFunc("/register", register)     // 注册
	http.HandleFunc("/login", login)           // 登录
	http.HandleFunc("/friends", handleFriends) // 添加好友
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
		var userForCreationDto dto.UserForCreationDto
		decoder := json.NewDecoder(r.Body)
		_ = decoder.Decode(&userForCreationDto)
		// 判断是否已经注册
		exist, err := method.UserExist(userForCreationDto.Account)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "注册失败！错误原因：%v", err)
			return
		} else if exist {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "注册失败！该用户已存在")
			return
		}
		// dot->entity
		user := method.MapUser(userForCreationDto)
		err = method.AddUser(user)
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

func login(w http.ResponseWriter, r *http.Request) {
	cros(&w)
	switch r.Method {
	case http.MethodGet:
		query := r.URL.Query()
		account := query["account"][0]
		passwd := query["passwd"][0]
		token, err := method.CheckLogin(account, passwd)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "登录错误：%v", err)
			return
		} else {
			if token == "" {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "账号或密码错误")
				return
			} else {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, token)
			}
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleFriends(w http.ResponseWriter, r *http.Request) {
	cros(&w)

	// 处理输入数据
	var friendDto dto.FriendDto
	decoder := json.NewDecoder(r.Body)
	_ = decoder.Decode(&friendDto)
	token := r.Header.Get("Authorization")[7:] // 获取token
	account := tool.ParseToken(token)          // 获取token
	if account == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	friend := friendDto.Account

	switch r.Method {
	case http.MethodGet:
		friends := method.GetFriends(account)
		res, _ := json.Marshal(friends)
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	case http.MethodPost:
		// 判断被添加者的账号是否存在
		exist, err := method.UserExist(friend)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "添加失败！错误原因：%v", err)
			return
		} else if !exist {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "添加失败！不存在该用户")
			return
		}
		// 判断是否已经存在该好友
		if method.FriendExist(account, friend) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "添加失败！好友已存在")
			return
		}
		method.AddFriend(account, friend)
		w.WriteHeader(http.StatusNoContent)
	case http.MethodDelete:
		// 不能删除自己
		if account == friend {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "删除失败！不能删除自己")
			return
		}
		if method.DeleteFriend(account, friend) {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "删除失败！")
		}
	case http.MethodOptions: // 对于post请求，浏览器首先会发option请求，如果服务器响应完全符合请求要求，浏览器则会发送真正的post请求。
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
