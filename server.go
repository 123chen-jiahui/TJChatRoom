package main

import (
	"encoding/json"
	"fmt"
	"github.com/dto"
	"github.com/gorilla/websocket"
	"github.com/method"
	"github.com/socket"
	"github.com/tool"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/ws", wsEndpoint)           // 建立websocket连接
	http.HandleFunc("/register", register)       // 注册
	http.HandleFunc("/login", login)             // 登录
	http.HandleFunc("/friends", handleFriends)   // 添加好友
	http.HandleFunc("/groups", handleGroups)     // 群聊
	http.HandleFunc("/messages", handleMessages) // 发送信息
	http.ListenAndServe(":8888", nil)
}

// 创建群聊请求格式
/*
{
	name: string
	owner: string
	members: []string
}
*/
func handleGroups(writer http.ResponseWriter, request *http.Request) {
	cros(&writer)
	token := request.Header.Get("Authorization")
	if len(token) < 8 {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	account := tool.ParseToken(token[7:]) // 获取token
	if account == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	switch request.Method {
	case http.MethodPost:
		var groupForCreationDto dto.GroupForCreationDto
		decoder := json.NewDecoder(request.Body)
		_ = decoder.Decode(&groupForCreationDto)
		group := groupForCreationDto.MapToGroup()
		// 懒得检查了，直接添加
		method.AddGroup(group)
	case http.MethodGet: // 获取所有群聊信息（包括群聊_id)，用于用户登录时使用
		groups := method.GetGroups(account)
		res, _ := json.Marshal(groups)
		writer.Write(res)
	case http.MethodPut: // 群聊中添加成员
		var groupForUpdateDto dto.GroupForUpdateDto
		decoder := json.NewDecoder(request.Body)
		_ = decoder.Decode(&groupForUpdateDto)
		groupId := groupForUpdateDto.Id
		list := groupForUpdateDto.List
		if groupForUpdateDto.Method == "add" {
			method.AddMemberToGroup(groupId, list)
		} else if groupForUpdateDto.Method == "delete" {
			fmt.Println("nothing")
		} else {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodOptions:
		writer.WriteHeader(http.StatusOK)
	}
}

// 信息格式
/*
{
	sender: string
	receiver: string
	isGroup: bool
	groupId: xxx
	contentType: xxx
	content: string
	sendTime: time
}
*/
func handleMessages(writer http.ResponseWriter, request *http.Request) {
	cros(&writer)
}

// 解决跨域问题
func cros(writer *http.ResponseWriter) {
	(*writer).Header().Set("Access-Control-Allow-Origin", "*")
	(*writer).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
	(*writer).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func reader(conn *websocket.Conn) {
	fmt.Println("jjjjjjjjjjjjjjjjjjj")
	// 将客户端信息放入liveClients中
	_, p, _ := conn.ReadMessage()
	chHb := socket.AddToLiveClients(string(p), conn)
	// 开一个goroutine发送心跳检测
	//go socket.HeartBeat(string(p), conn, chHb)

	for {
		// read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		if string(p) == "pong" {
			chHb <- true
		}

		fmt.Println("收到", string(p))

		//if err := conn.WriteMessage(messageType, []byte("Hello from GoLang!")); err != nil {
		//	log.Println(err)
		//	return
		//}
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
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
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
		method.AddFriend(friend, account)
		w.WriteHeader(http.StatusNoContent)
	case http.MethodDelete:
		// 不能删除自己
		if account == friend {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "删除失败！不能删除自己")
			return
		}
		if method.DeleteFriend(account, friend) && method.DeleteFriend(friend, account) {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "删除失败！")
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
