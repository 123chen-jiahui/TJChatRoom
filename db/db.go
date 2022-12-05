package db

import (
	"context"
	"fmt"
	"github.com/entity"
	"github.com/tool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var DB *mongo.Database

func initDB() *mongo.Client {
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI(tool.MongoUrl)
	// 连接到MongoDB
	mongoClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic("无法连接到mongoDB" + err.Error())
	}
	// 检查连接
	err = mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		panic("无法连接到mongoDB" + err.Error())
	}
	return mongoClient
}

func init() {
	DB = initDB().Database(tool.MConfig.DbName)
}

func InsertUser(user entity.User) error {
	table := DB.Collection("User")
	_, err := table.InsertOne(context.TODO(), user)
	return err
}

func UpdateUserAvatar(account, avatarUrl string) {
	table := DB.Collection("User")
	table.UpdateOne(
		context.TODO(),
		bson.M{"account": account},
		bson.M{"$set": bson.M{"avatar": avatarUrl}})
}

func UpdateUserNickName(account, nickName string) {
	table := DB.Collection("User")
	table.UpdateOne(
		context.TODO(),
		bson.M{"account": account},
		bson.M{"$set": bson.M{"nickname": nickName}})
}

func UpdateUserPassword(account, password string) {
	table := DB.Collection("User")
	table.UpdateOne(
		context.TODO(),
		bson.M{"account": account},
		bson.M{"$set": bson.M{"password": password}})
}

func PushFriend(account, friend string) {
	table := DB.Collection("User")
	filter := bson.M{"account": account}

	f := new(entity.Friend)
	f.Friend = friend
	update := bson.M{"$push": bson.M{"friends": f}}
	_, err := table.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("添加成功")
}

// PullFriend 删除account的好友friend
// 返回值为删除的数目
func PullFriend(account, friend string) int64 {
	table := DB.Collection("User")
	filter := bson.M{"account": account}
	update := bson.M{"$pull": bson.M{"friends": bson.M{"friend": friend}}}

	one, err := table.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println(err)
		return 0
	}
	return one.ModifiedCount
}

func FriendExist(account, friend string) bool {
	table := DB.Collection("User")
	filter := bson.M{"friends": bson.M{"friend": friend}, "account": account}
	err := table.FindOne(context.TODO(), filter).Err()
	return !(err == mongo.ErrNoDocuments)
}

func GetFriends(account string) []entity.User {
	var me entity.User
	table := DB.Collection("User")
	table.FindOne(context.TODO(), bson.M{"account": account}).Decode(&me)
	friends := me.Friends

	var users []entity.User
	var tmp entity.User
	for _, f := range friends {
		table.FindOne(context.TODO(), bson.M{"account": f.Friend}).Decode(&tmp)
		users = append(users, tmp)
	}
	fmt.Println(users)
	return users
}

func FindUserByAccount(account string) (entity.User, error) {
	var user entity.User
	table := DB.Collection("User")
	table.Find(context.TODO(), bson.M{"account": account})
	err := table.FindOne(context.TODO(), bson.M{"account": account}).Decode(&user)
	fmt.Println("hello, world", user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("未找到记录")
		}
	}
	return user, err
}

//
// Group
//

func InsertGroup(group entity.Group) {
	group.Id = primitive.NewObjectID()
	table := DB.Collection("Group")
	_, _ = table.InsertOne(context.TODO(), group)
	fmt.Println(group.Id.String())
}

func RemoveGroup(groupId string) {
	objId, _ := primitive.ObjectIDFromHex(groupId)
	table := DB.Collection("Group")
	table.DeleteOne(context.TODO(), bson.M{"_id": objId})
}

func GetGroups(account string) []entity.Group {
	table := DB.Collection("Group")
	filter := bson.M{"members": bson.M{"member": account}}
	c, _ := table.Find(context.TODO(), filter)
	var groups []entity.Group
	_ = c.All(context.TODO(), &groups)
	fmt.Println(groups)
	return groups
}

func GetGroupById(groupId string) entity.Group {
	var group entity.Group
	objId, _ := primitive.ObjectIDFromHex(groupId)
	table := DB.Collection("Group")
	_ = table.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&group)
	return group
}

func AddMember(groupId string, members []string) {
	var group entity.Group
	table := DB.Collection("Group")
	objId, err := primitive.ObjectIDFromHex(groupId)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(objId)
	filter := bson.M{"_id": objId}
	err = table.FindOne(context.TODO(), filter).Decode(&group)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 防止重复添加，使用map以减小时间开销
	e := make(map[string]bool)
	for _, m := range group.Members {
		e[m.Member] = true
	}
	for _, m := range members {
		if _, ok := e[m]; ok {
			continue
		}
		member := new(entity.Member)
		member.Member = m
		update := bson.M{"$push": bson.M{"members": member}}
		_, err = table.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			fmt.Println("添加成员失败", err)
			return
		}
		e[m] = true
	}
}

func PullMembers(groupId string, members []string) {
	group := GetGroupById(groupId)
	objId := group.Id
	memberNum := len(group.Members)
	var ownerKilled = false
	table := DB.Collection("Group")
	for _, m := range members {
		res, err := table.UpdateOne(context.TODO(),
			bson.M{"_id": objId},
			bson.M{"$pull": bson.M{"members": bson.M{"member": m}}})
		if res.ModifiedCount != 0 {
			memberNum -= 1
			if m == group.Owner {
				ownerKilled = true
			}
		}
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	// 考虑成员数量为0或群主退出的情况
	if memberNum == 0 {
		RemoveGroup(groupId)
	} else if ownerKilled {
		group = GetGroupById(groupId)
		table.UpdateOne(context.TODO(), bson.M{"_id": objId}, bson.M{"$set": bson.M{"owner": group.Members[0].Member}})
	}
}

func CheckAndRemoveGroup(groupId string, account string) bool {
	group := GetGroupById(groupId)
	if group.Owner != account {
		return false
	} else {
		RemoveGroup(groupId)
		return true
	}
}

//
// message
//

func AddMessage(message entity.Message) {
	fmt.Println(message)
	table := DB.Collection("Message")
	table.InsertOne(context.TODO(), message)
}

// GetUnreadMessages 获取普通聊天未读记录
func GetUnreadMessages(account string) []entity.Message {
	table := DB.Collection("Message")
	filter := bson.M{"to": account, "read": false, "group": ""}
	c, _ := table.Find(context.TODO(), filter)
	var messages []entity.Message
	_ = c.All(context.TODO(), &messages)
	return messages
}

// GetLatestHistory 获取a和b近期的聊天记录
func GetLatestHistory(me string, opposite string, num int64) []entity.Message {
	var messages []entity.Message
	table := DB.Collection("Message")
	// 过滤器：已读、且是双方之间的信息
	xx := []bson.M{
		{"from": me, "to": opposite, "group": ""},
		{"from": opposite, "to": me, "read": true, "group": ""},
	}
	filter := bson.M{"$or": xx}                           // 已读信息
	option1 := options.Find().SetLimit(num)               // 指定聊天记录数量
	option2 := options.Find().SetSort(bson.M{"time": -1}) // 最近的
	c, _ := table.Find(context.TODO(), filter, option1, option2)
	_ = c.All(context.TODO(), &messages)
	fmt.Println(messages)
	return messages
}

// MessagesReadUser 将消息设为已读（对象是用户）
func MessagesReadUser(me, opposite string) {
	table := DB.Collection("Message")
	filter := bson.M{"from": opposite, "to": me}
	update := bson.M{"$set": bson.M{"read": true}}
	_, _ = table.UpdateMany(context.TODO(), filter, update)
}

// MessagesReadGroup 将消息设为已读（对象是群聊）
func MessagesReadGroup(me, groupId string) {
	table := DB.Collection("Message")
	filter := bson.M{"group": groupId, "to": me}
	update := bson.M{"$set": bson.M{"read": true}}
	_, _ = table.UpdateMany(context.TODO(), filter, update)
}

// GetLatestHistoriesOfGroup 返回群聊的历史记录
func GetLatestHistoriesOfGroup(account, groupId string, num int64) []entity.Message {
	// var messages []entity.GroupMessage
	// table := DB.Collection("GroupMessage")
	// filter := bson.M{"group": groupId}
	// option1 := options.Find().SetLimit(num)
	// option2 := options.Find().SetSort(bson.M{"time": -1})
	// c, _ := table.Find(context.TODO(), filter, option1, option2)
	// _ = c.All(context.TODO(), &messages)
	// return messages

	var messages []entity.Message
	table := DB.Collection("Message")
	filter := bson.M{"group": groupId, "to": account, "read": true}
	option1 := options.Find().SetLimit(num)
	option2 := options.Find().SetSort(bson.M{"time": -1})
	c, _ := table.Find(context.TODO(), filter, option1, option2)
	_ = c.All(context.TODO(), &messages)
	return messages
}

func GetUnreadGroupMessages(account string) []entity.Message {
	table := DB.Collection("Message")
	filter := bson.M{"to": account, "read": false, "group": bson.M{"$ne": ""}}
	c, _ := table.Find(context.TODO(), filter)
	var messages []entity.Message
	_ = c.All(context.TODO(), &messages)
	return messages
}

// DeleteMessagesInGroup 删除某人在某群的聊天记录
// 我发的不能删，发给我的必须删
func DeleteMessagesInGroup(account, groupId string) {
	table := DB.Collection("Message")
	table.DeleteMany(context.TODO(), bson.M{"to": account, "group": groupId})
}
