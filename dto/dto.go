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
