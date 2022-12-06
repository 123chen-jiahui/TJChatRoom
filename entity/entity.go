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
	Avatar   string   // 头像
}

type Friend struct {
	Friend string `json:"friend"`
}

type Group struct {
	Id      primitive.ObjectID `bson:"_id"`
	Members []Member           // 群成员
	Name    string             // 群名称
	Owner   string             // 群主
	// CreateTime time.Time // 群创建时间
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
	Content     string // 信息内容，对于文件，表示文件名，确切说是上传时的文件名
	RemoteName  string // 仅用于信息内容为文件，表示存在oss中的文件名，为了避免重名
	// 获取文件资源时，用的时RemoteName，而用户下载文件的时候，用的还是Content
}
