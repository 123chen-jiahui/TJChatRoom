package method

import (
	"fmt"
	"github.com/db"
	"github.com/dto"
	"github.com/entity"
	"sort"
)

func AddMessages(messageForCreation dto.MessageForCreation) []entity.Message {
	var group entity.Group
	var message entity.Message
	var messages []entity.Message
	if messageForCreation.Group != "" {
		group = db.GetGroupById(messageForCreation.Group)
		// fmt.Println(group)
		for _, member := range group.Members {
			if member.Member == messageForCreation.From {
				continue
			}
			message = messageForCreation.MapToMessage(member.Member)
			messages = append(messages, message)
			db.AddMessage(message)
			// db.AddMessage(messageForCreation.MapToMessage(member.Member))
		}
	} else {
		message = messageForCreation.MapToMessage(messageForCreation.To)
		messages = append(messages, message)
		db.AddMessage(messageForCreation.MapToMessage(messageForCreation.To))
	}
	return messages
}

func GetAllMessages(account string) []dto.MessagesReturn {
	m := make(map[string][]dto.MessageDto)
	user, _ := db.FindUserByAccount(account)
	friends := user.Friends
	// 增加所有好友的历史记录
	for _, f := range friends {
		friend := f.Friend
		if friend == account {
			continue
		}
		histories := db.GetLatestHistory(account, friend, 10)
		fmt.Println(friend, "history", histories)
		for _, h := range histories {
			var flag int
			if h.From == account {
				flag = 1
			}
			m[friend] = append(m[friend], dto.MessageDto{
				Time:        h.Time,
				Group:       h.Group,
				ContentType: h.ContentType,
				Content:     h.Content,
				Flag:        flag,
			})
		}
	}
	fmt.Println("half way", m)
	// 增加未读记录
	messages := db.GetUnreadMessages(account)
	fmt.Println("未读记录", messages)
	for _, message := range messages {
		m[message.From] = append(m[message.From], dto.MessageDto{
			Time:        message.Time,
			Group:       message.Group,
			ContentType: message.ContentType,
			Content:     message.Content,
			Flag:        2,
		})
	}
	// 整理、排序
	var messagesReturns []dto.MessagesReturn
	for k, v := range m {
		sort.Sort(dto.MessageDtoSlice(v))
		messagesReturns = append(messagesReturns, dto.MessagesReturn{
			From:     k,
			Messages: v,
		})
	}
	return messagesReturns
}

func GetLatestHistory(a, b string, num int64) []entity.Message {
	return db.GetLatestHistory(a, b, num)
}

func SetMessagesRead(me, opposite, isGroup string) {
	if isGroup == "1" { // opposite 表示群聊号
		db.MessagesReadGroup(me, opposite)
	} else { // opposite 表示用户账号
		db.MessagesReadUser(me, opposite)
	}
}
