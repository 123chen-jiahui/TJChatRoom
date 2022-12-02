package socket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type LiveClient struct {
	Account string
	Conn    *websocket.Conn
	ChMsg   chan string
	ChHeal  chan bool
}

const INTERVAL = 2000 // 心跳检测时间间隔（ms）

var liveClients []LiveClient
var mu sync.Mutex

func Debug() {
	for {
		mu.Lock()
		fmt.Println(liveClients)
		mu.Unlock()
		time.Sleep(time.Second * 2)
	}
}

func AddClient(account string, conn *websocket.Conn) LiveClient {
	mu.Lock()
	defer mu.Unlock()
	// 坑点：go中for range并不会改变值，例如
	// a := []int{1, 2, 3}
	// for i, v := range a {
	// 	v = 100 * i
	// 	fmt.Println(v, "->", a[i])
	// }
	// 若要改变值：
	// a := []int{1, 2, 3}
	// for i, v := range a {
	// 	a[i] = 100 * i
	// 	fmt.Println(v, "->", a[i])
	// }
	for i, ele := range liveClients {
		if ele.Account == account {
			liveClients[i].Conn = conn
			return liveClients[i]
		}
	}
	newClient := LiveClient{
		Account: account,
		Conn:    conn,
		ChMsg:   make(chan string),
		ChHeal:  make(chan bool),
	}
	liveClients = append(liveClients, newClient)
	return newClient
}

func FindClient(account string) *LiveClient {
	mu.Lock()
	defer mu.Unlock()
	for i, ele := range liveClients {
		if account == ele.Account {
			return &liveClients[i]
			// return ele.conn
		}
	}
	return nil
}

func RemoveClient(account string) {
	mu.Lock()
	defer mu.Unlock()
	for i, ele := range liveClients {
		if ele.Account == account {
			liveClients = append(liveClients[:i], liveClients[i+1:]...)
			return
		}
	}
}

// Reader
// 客户端给服务器的反馈信息都通过这个函数，并通过信道通知其他函数
// 如果客户端断开连接或者更新了连接，那么err!=nil，该函数返回
func (liveClient *LiveClient) Reader() {
	account := liveClient.Account
	conn := liveClient.Conn
	chMsg := liveClient.ChMsg
	// go heartBeat(account, conn, chMsg)
	for {
		fmt.Println("reader is alive!")
		// fmt.Println(conn)
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("客户端断开连接")
			RemoveClient(account)
			return
		}
		fmt.Println("收到了来自客户端的反馈：", string(p))
		chMsg <- string(p)
	}
}

func heartBeat(account string, conn *websocket.Conn, ch chan bool) {
	err := conn.WriteMessage(websocket.TextMessage, []byte("ping")) // 发送心跳检查
	if err != nil {
		fmt.Println("心跳检测发送失败！")
		RemoveClient(account)
		return
	}
	timeSent := time.Now()
	for {
		time.Sleep(time.Millisecond * 100)
		select {
		case flag := <-ch:
			if !flag {
				fmt.Println("客户端退出")
				RemoveClient(account)
				return
			}
			fmt.Println("心跳检测成功！")
			timeSent = time.Now()
			err := conn.WriteMessage(websocket.TextMessage, []byte("ping")) // 发送心跳检查
			if err != nil {
				fmt.Println("心跳检测发送失败！")
				RemoveClient(account)
				return
			}
		default:
			interval := time.Since(timeSent).Milliseconds()
			if interval > INTERVAL { // 超时
				fmt.Println("连接超时！")
				RemoveClient(account)
				return
			}
		}
	}
}
