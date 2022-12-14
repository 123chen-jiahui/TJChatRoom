package method

import (
	"fmt"
	"github.com/db"
	"github.com/dto"
	"github.com/entity"
	"sort"
)

func ExtendMessages(messageForCreation dto.MessageForCreation) []entity.Message {
	var group entity.Group
	var message entity.Message
	var messages []entity.Message
	if messageForCreation.Group != "" {
		group = db.GetGroupById(messageForCreation.Group)
		for _, member := range group.Members {
			// 我发送给我自己，也要存库。必定是已读的，这点会在入库时体现
			message = messageForCreation.MapToMessage(member.Member)
			messages = append(messages, message)
			// db.AddMessage(message)
		}
	} else {
		message = messageForCreation.MapToMessage(messageForCreation.To)
		messages = append(messages, message)
		// db.AddMessage(messageForCreation.MapToMessage(messageForCreation.To))
	}
	return messages
}

func AddMessage(msg entity.Message, read bool) {
	msg.Read = read
	db.AddMessage(msg)
}

// GetAllMessages 获取n条历史记录+所有未读消息
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
				RemoteName:  h.RemoteName,
				Flag:        flag,
			})
		}
	}
	// 增加所有群聊的历史记录
	groups := GetGroups(account)
	for _, group := range groups {
		histories := db.GetLatestHistoriesOfGroup(account, group.Id.Hex(), 10)
		for _, h := range histories {
			var flag int
			if h.From == account {
				flag = 1
			}
			m[group.Id.Hex()] = append(m[group.Id.Hex()], dto.MessageDto{
				Time:        h.Time,
				Group:       h.Group,
				GroupMember: h.From,
				ContentType: h.ContentType,
				Content:     h.Content,
				RemoteName:  h.RemoteName,
				Flag:        flag,
			})
		}
	}

	// 增加未读记录
	messages := db.GetUnreadMessages(account)
	fmt.Println("未读记录", messages)
	for _, message := range messages {
		m[message.From] = append(m[message.From], dto.MessageDto{
			Time:        message.Time,
			Group:       message.Group,
			ContentType: message.ContentType,
			Content:     message.Content,
			RemoteName:  message.RemoteName,
			Flag:        2,
		})
	}

	// 增加群聊未读记录
	groupMessages := db.GetUnreadGroupMessages(account)
	for _, message := range groupMessages {
		m[message.Group] = append(m[message.Group], dto.MessageDto{
			Time:        message.Time,
			Group:       message.Group,
			GroupMember: message.From,
			ContentType: message.ContentType,
			Content:     message.Content,
			RemoteName:  message.RemoteName,
			Flag:        2,
		})
	}
	// 整理、排序
	var messagesReturns []dto.MessagesReturn
	for k, v := range m {
		sort.Sort(dto.MessageDtoSlice(v))
		isGroup := false
		if len(k) > 10 {
			isGroup = true
		}
		messagesReturns = append(messagesReturns, dto.MessagesReturn{
			From:     k,
			IsGroup:  isGroup,
			Messages: v,
		})
	}
	return messagesReturns
}

func SetMessagesRead(me, opposite, isGroup string) {
	if isGroup == "1" { // opposite 表示群聊号
		db.MessagesReadGroup(me, opposite)
	} else { // opposite 表示用户账号
		db.MessagesReadUser(me, opposite)
	}
}
