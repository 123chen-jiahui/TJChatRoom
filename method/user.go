package method

import (
	"fmt"
	"github.com/db"
	"github.com/dto"
	"github.com/entity"
	"github.com/tool"
	"go.mongodb.org/mongo-driver/mongo"
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

func UserExist(account string) (bool, error) {
	_, err := db.FindUserByAccount(account)
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	return true, err
}

func CheckLogin(account, passwd string) (token string, err error) {
	token = ""
	err = nil
	user, err := db.FindUserByAccount(account)
	fmt.Println(user)
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

func DeleteFriend(account, friend string) bool {
	return db.PullFriend(account, friend) != 0
}

func FriendExist(account, friend string) bool {
	return db.FriendExist(account, friend)
}
