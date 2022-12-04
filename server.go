package main

import (
	"encoding/json"
	"fmt"
	"github.com/dto"
	"github.com/entity"
	"github.com/gorilla/websocket"
	"github.com/method"
	"github.com/socket"
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
	go socket.Debug()
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
	if request.Method == http.MethodOptions {
		writer.WriteHeader(http.StatusOK)
		return
	}
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
			writer.WriteHeader(http.StatusNoContent)
		} else if groupForUpdateDto.Method == "delete" {
			method.DeleteMembersFromGroup(groupId, list)
			writer.WriteHeader(http.StatusNoContent)
		} else {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		query := request.URL.Query()
		groupId := query["id"][0]
		if ok := method.DeleteGroup(groupId, account); !ok {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("解散群聊失败"))
		} else {
			writer.WriteHeader(http.StatusNoContent)
		}
	}
}

// 该函数还要改
func notice(msg entity.Message) {
	account := msg.To
	clientPtr := socket.FindClient(account)
	if clientPtr == nil {
		log.Println("对方不在线")
		method.AddMessage(msg, false)
		return
	}
	client := *clientPtr
	res, _ := json.Marshal(msg)
	err := client.Conn.WriteMessage(websocket.BinaryMessage, res)
	if err != nil {
		log.Println("发送失败", err)
		method.AddMessage(msg, false)
		return
	}
	timeStart := time.Now()
	for {
		select {
		case <-client.ChMsg:
			fmt.Println("消息已读")
			method.AddMessage(msg, true)
			return
		default:
			if interval := time.Since(timeStart).Milliseconds(); interval > 500 {
				method.AddMessage(msg, false)
				return
			}
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func handleMessages(writer http.ResponseWriter, request *http.Request) {
	cros(&writer)
	if request.Method == http.MethodOptions {
		writer.WriteHeader(http.StatusOK)
		return
	}
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
		var messageForCreationDto dto.MessageForCreation
		decoder := json.NewDecoder(request.Body)
		decoder.Decode(&messageForCreationDto)
		// if messageForCreationDto.Group != "" { // 若为群聊信息，直接加入
		// 	method.AddGroupMessage(messageForCreationDto.MapToGroupMessage())
		// }
		// 不应该立刻存到数据库中，如果用户处于实时聊天的状态，这样会增加更新数据库的次数
		// 并且对于信息的已读或未读状态记录与实际情况不符。这种情况下，信息一定是已读的，
		// 但是存到数据库中的信息确实未读
		// 解决方法：先通知目标用户，并且程序等待一段时间，在这段时间内，
		// 如果用户对该信息给出了正反馈，那么将信息的状态置为已读然后入库；如果用户无响应，
		// 或是给出了负反馈，那么就将信息的状态置为未读然后入库
		messages := method.ExtendMessages(messageForCreationDto) // 扩展信息
		fmt.Println("messages is ", messages)
		// 通知信息接收者
		for _, msg := range messages {
			go notice(msg) // notice可能不会立刻返回，所以开一个go routine
		}
	case http.MethodGet: // 返回所有未读数据，同时，需要返回最近10条聊天记录
		// TODO
		// 该函数有缺陷，因为目前只考虑了对象使用户的情况
		// 而没有考虑对象是群聊的情况
		// 思路：在请求中加入参数表示是否对象是否是群聊，其余逻辑类似
		messages := method.GetAllMessages(account)
		res, _ := json.Marshal(messages)
		writer.Write(res)
	case http.MethodPut: // 将消息设为已读
		query := request.URL.Query()
		opposite := query["opposite"][0]
		isGroup := query["group"][0]
		method.SetMessagesRead(account, opposite, isGroup)
	}
}

// 解决跨域问题
func cros(writer *http.ResponseWriter) {
	(*writer).Header().Set("Access-Control-Allow-Origin", "*")
	(*writer).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
	(*writer).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func reader(conn *websocket.Conn) {
	// 将客户端信息放入liveClients中
	// _, p, _ := conn.ReadMessage()
	// chHb := socket.AddToLiveClients(string(p), conn)
	// 开一个goroutine发送心跳检测
	// go socket.HeartBeat(string(p), conn, chHb)
	go socket.Debug()
	for {
		// read in a message
		_, _, err := conn.ReadMessage()
		if err != nil {
			if err == websocket.ErrCloseSent {
				fmt.Println("寄")
			}
			log.Println(err)
			return
		}
		// print out that message for clarity
		// if string(p) == "pong" {
		// 	chHb <- true
		// }

		// fmt.Println("收到", string(p))
		err = conn.WriteMessage(websocket.TextMessage, []byte("ping"))
		if err != nil {
			log.Println("发送寄了")
		}

		// if err := conn.WriteMessage(messageType, []byte("Hello from GoLang!")); err != nil {
		// 	log.Println(err)
		// 	return
		// }
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("新的websocket连接")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	_, p, err := ws.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}
	account := string(p)
	client := socket.AddClient(account, ws)
	// reader(ws)
	go client.Reader()
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
