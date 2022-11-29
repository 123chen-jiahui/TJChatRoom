package method

import (
	"github.com/db"
	"github.com/dto"
	"github.com/entity"
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
	messages := db.GetMessages(account)
	m := make(map[string][]dto.MessageDto)
	for _, message := range messages {
		m[message.From] = append(m[message.From], dto.MessageDto{
			Time:        message.Time,
			Group:       message.Group,
			ContentType: message.ContentType,
			Content:     message.Content,
		})
	}
	var messagesReturns []dto.MessagesReturn
	for k, v := range m {
		messagesReturns = append(messagesReturns, dto.MessagesReturn{
			From:     k,
			Messages: v,
		})
	}
	return messagesReturns
}
