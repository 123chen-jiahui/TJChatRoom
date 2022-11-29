package dto

import (
	"github.com/entity"
)

type UserForCreationDto struct {
	Account  string `json:"account"`
	Passwd   string `json:"passwd"`
	NickName string `json:"nickName"`
}

type UserInfoDto struct {
	Account  string `json:"account"`
	Passwd   string `json:"passwd"`
	NickName string `json:"nickName"`
}

type FriendDto struct {
	Account string `json:"account"`
}

type GroupForCreationDto struct {
	Name    string   `json:"name"`
	Owner   string   `json:"owner"`
	Members []string `json:"members"`
}

type GroupForUpdateDto struct {
	Method string   `json:"method"`
	Id     string   `json:"id"`
	List   []string `json:"list"`
}

type MessageForCreation struct {
	Time        int64  `json:"time"`
	Group       string `json:"group"`       // 信息所属群聊（若是私信，此项为空）
	From        string `json:"from"`        // 信息发送者
	To          string `json:"to"`          // 信息接收者
	Read        bool   `json:"read"`        // 信息接收者是否已经阅读过该消息
	ContentType int    `json:"contentType"` // 信息内容类型（文本or文件）
	Content     string `json:"content"`     // 信息内容
}

// MessageDto 返回
type MessageDto struct {
	Time        int64
	Group       string
	ContentType int
	Content     string
}

type MessagesReturn struct {
	From     string
	Messages []MessageDto
}

func (g GroupForCreationDto) MapToGroup() entity.Group {
	var group = entity.Group{
		Name:  g.Name,
		Owner: g.Owner,
	}
	var members []entity.Member
	for _, ele := range g.Members {
		members = append(members, entity.Member{Member: ele})
	}
	group.Members = members
	return group
}

func (m MessageForCreation) MapToMessage(to string) entity.Message {
	return entity.Message{
		Time:        m.Time,
		Group:       m.Group,
		From:        m.From,
		To:          to,
		Read:        false,
		ContentType: m.ContentType,
		Content:     m.Content,
	}
}
