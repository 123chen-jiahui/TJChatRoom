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
	Group       string // 无用
	GroupMember string // 只有当IsGroup为true时，这项才有用，表示群里的谁
	ContentType int
	Content     string
	Flag        int
	// Flag为0或1表示历史记录，并且0表示别人发给我的，1表示我发给别人的
	// Flag为2表示未读信息
}

type MessageDtoSlice []MessageDto

func (MDS MessageDtoSlice) Len() int           { return len(MDS) }
func (MDS MessageDtoSlice) Swap(i, j int)      { MDS[i], MDS[j] = MDS[j], MDS[i] }
func (MDS MessageDtoSlice) Less(i, j int) bool { return MDS[i].Time < MDS[j].Time }

type MessagesReturn struct {
	From     string // 聊天对象，如果IsGroup为true，表示GroupId
	IsGroup  bool   // 聊天对象是否是群聊
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
