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

func FindUserByAccount(account string) (entity.User, error) {
	var user entity.User
	table := DB.Collection("User")
	err := table.FindOne(context.TODO(), bson.M{"account": account}).Decode(&user)
	fmt.Println("hello, world", user)
	if err != nil {
		fmt.Println("holy shit")
	}
	return user, err
}

// FindOnd 查询
func FindOnd(collection string) {
	var result bson.M
	table := DB.Collection(collection)
	err := table.FindOne(context.TODO(), bson.M{"name": "CJH"}).Decode(&result) //
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("%v\n", result)

	fmt.Println(result["age"])
	//v, err := encoder.Encode(result, encoder.SortMapKeys)
	//if err != nil {
	//	fmt.Printf("%v\n", err)
	//}
	//fmt.Println(string(v))
}
