package db

import (
	"context"
	"fmt"
	"github.com/entity"
	"github.com/tool"
	"go.mongodb.org/mongo-driver/bson"
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
