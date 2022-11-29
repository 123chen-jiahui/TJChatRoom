package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 实体集

type User struct {
	Account  string
	Passwd   string
	NickName string
	Friends  []Friend // 所加好友
	Group    []string // 加入的群聊
}

type Friend struct {
	Friend string `json:"friend"`
}

type Group struct {
	Id      primitive.ObjectID `bson:"_id"`
	Members []Member           // 群成员
	Name    string             // 群名称
	Owner   string             // 群主
	//CreateTime time.Time // 群创建时间
}

type Member struct {
	Member string
}

type Message struct {
	Time        int64  // 信息发送时间
	Group       string // 信息所属群聊（若是私信，此项为空）
	From        string // 信息发送者
	To          string // 信息接收者
	Read        bool   // 信息接收者是否已经阅读过该消息
	ContentType int    // 信息内容类型（文本or文件）
	Content     string // 信息内容
}
