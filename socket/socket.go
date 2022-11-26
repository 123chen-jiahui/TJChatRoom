package socket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type LiveClient struct {
	account string
	conn    *websocket.Conn
	chHB    chan bool
}

const INTERVAL = 2000 // 心跳检测时间间隔（ms）

var liveClients []LiveClient
var mu sync.Mutex

// AddToLiveClients 这里需要检查是否重复登录（不过没有完成）
func AddToLiveClients(account string, conn *websocket.Conn) chan bool {
	mu.Lock()
	defer mu.Unlock()
	for _, ele := range liveClients {
		if ele.account == account {
			ele.conn = conn
			return ele.chHB
		}
	}
	newCh := make(chan bool)
	liveClients = append(liveClients, LiveClient{
		account: account,
		conn:    conn,
		chHB:    newCh,
	})
	return newCh
}

func RemoveFromLiveClients(account string) {
	mu.Lock()
	defer mu.Unlock()
	for i, ele := range liveClients {
		if ele.account == account {
			liveClients = append(liveClients[:i], liveClients[i+1:]...)
			return
		}
	}
}

func HeartBeat(account string, conn *websocket.Conn, ch chan bool) {
	err := conn.WriteMessage(websocket.TextMessage, []byte("ping")) // 发送心跳检查
	if err != nil {
		fmt.Println("心跳检测发送失败！")
		RemoveFromLiveClients(account)
		return
	}
	timeSent := time.Now()
	for {
		time.Sleep(time.Millisecond * 100)
		select {
		case <-ch:
			fmt.Println("心跳检测成功！")
			timeSent = time.Now()
			err := conn.WriteMessage(websocket.TextMessage, []byte("ping")) // 发送心跳检查
			if err != nil {
				fmt.Println("心跳检测发送失败！")
				RemoveFromLiveClients(account)
				return
			}
		default:
			interval := time.Since(timeSent).Milliseconds()
			if interval > INTERVAL { // 超时
				fmt.Println("连接超时！")
				RemoveFromLiveClients(account)
				return
			}
		}
	}
}
