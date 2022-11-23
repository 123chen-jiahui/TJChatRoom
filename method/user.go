package method

import (
	"github.com/db"
	"github.com/dto"
	"github.com/entity"
	"github.com/tool"
)

func MapUser(from dto.UserForCreationDto) (to entity.User) {
	to = entity.User{
		Account:  from.Account,
		Passwd:   from.Passwd,
		NickName: from.NickName,
		Friends:  nil,
		Group:    nil,
	}
	initFriend := entity.Friend{Friend: to.Account}
	to.Friends = append(to.Friends, initFriend)
	return to
}

func CheckLogin(account, passwd string) (token string, err error) {
	token = ""
	err = nil
	user, err := db.FindUserByAccount(account)
	if err != nil {
		return
	}
	if user.Passwd == passwd {
		token, err = tool.GenerateToken(account)
	}
	return
}

func AddUser(user entity.User) (err error) {
	err = db.InsertUser(user)
	return
}

func AddFriend(account, friend string) {
	db.PushFriend(account, friend)
}
