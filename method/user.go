package method

import (
	"fmt"
	"github.com/db"
	"github.com/dto"
	"github.com/entity"
	"github.com/tool"
	"go.mongodb.org/mongo-driver/mongo"
)

// MapUser 将UserForCreationDto映射为User
func MapUser(from dto.UserForCreationDto) (to entity.User) {
	to = entity.User{
		Account:  from.Account,
		Passwd:   from.Passwd,
		NickName: from.NickName,
		Friends:  nil,
		Group:    nil,
		Avatar:   "default.png",
	}
	initFriend := entity.Friend{Friend: to.Account}
	to.Friends = append(to.Friends, initFriend)
	return to
}

// MapUserToUserInfoDto 将User映射为UserInfoDto
func MapUserToUserInfoDto(from entity.User) (to dto.UserInfoDto) {
	to = dto.UserInfoDto{
		Account:  from.Account,
		Passwd:   from.Passwd,
		NickName: from.NickName,
		Avatar:   from.Avatar,
	}
	return
}

func UserExist(account string) (bool, error) {
	_, err := db.FindUserByAccount(account)
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	return true, err
}

func FindUser(account string) dto.UserInfoDto {
	user, _ := db.FindUserByAccount(account)
	return MapUserToUserInfoDto(user)
}

func UpdateUser(account, avatar, nickName, password string) {
	if avatar != "" {
		db.UpdateUserAvatar(account, avatar)
	}
	if nickName != "" {
		db.UpdateUserNickName(account, nickName)
	}
	if password != "" {
		db.UpdateUserPassword(account, password)
	}
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

func GetFriends(account string) []dto.UserInfoDto {
	users := db.GetFriends(account)
	var usersInfoDto []dto.UserInfoDto
	for _, u := range users {
		usersInfoDto = append(usersInfoDto, MapUserToUserInfoDto(u))
	}
	return usersInfoDto
}

func DeleteFriend(account, friend string) bool {
	return db.PullFriend(account, friend) != 0
}

func FriendExist(account, friend string) bool {
	return db.FriendExist(account, friend)
}
