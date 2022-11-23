package tool

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type MongoConfig struct {
	User   string `json:"mongoUser"`
	Passwd string `json:"mongoPasswd"`
	Host   string `json:"mongoHost"`
	Port   string `json:"mongoPort"`
	DbName string `json:"mongoDbName"`
}

var MConfig MongoConfig
var MongoUrl string

func init() {
	file, err := os.Open("./config/config.json")
	if err != nil {
		log.Fatalf("读取配置文件出错：%v\n", err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&MConfig)
	if err != nil {
		log.Fatalf("反序列化配置文件：%v\n", err)
	}

	// "mongodb://root:123456@118.31.108.144:27017"
	MongoUrl = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
		MConfig.User,
		MConfig.Passwd,
		MConfig.Host,
		MConfig.Port,
		MConfig.DbName,
	)
}
