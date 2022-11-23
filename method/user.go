package method

import (
	"fmt"
	"github.com/db"
	"github.com/entity"
	"github.com/tool"
)

func CheckLogin(account, passwd string) (token string, err error) {
	token = ""
	err = nil
	user, err := db.FindUserByAccount(account)
	if err != nil {
		return
	}
	if user.Passwd == passwd {
		token, err = tool.GenerateToken(account)
		if err != nil {
			fmt.Println("寄")
		}
	}
	return
}

func AddUser(user entity.User) error {
	f := entity.Friend{
		Friend: user.Account,
	}
	user.Friends = append(user.Friends, f)
	err := db.InsertUser(user)
	return err
}

func AddFriend(account, friend string) {
	db.PushFriend(account, friend)
}
