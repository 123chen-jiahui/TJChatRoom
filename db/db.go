package db

import (
	"context"
	"fmt"
	"github.com/entity"
	"github.com/tool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		panic("无法连接到mongoDB")
	}
	return mongoClient
}

func init() {
	DB = initDB().Database(tool.MConfig.DbName)
}

func InsertUser(user entity.User) error {
	table := DB.Collection("user")
	_, err := table.InsertOne(context.TODO(), user)
	return err
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
