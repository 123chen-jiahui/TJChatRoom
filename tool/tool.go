package tool

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"os"
)

type MongoConfig struct {
	User            string `json:"mongoUser"`
	Passwd          string `json:"mongoPasswd"`
	Host            string `json:"mongoHost"`
	Port            string `json:"mongoPort"`
	DbName          string `json:"mongoDbName"`
	SecretStr       string `json:"secretStr"`
	EndPoint        string `json:"endPoint"`
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	BucketName      string `json:"bucketName"`
}

var MConfig MongoConfig
var OssClient *oss.Client
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
	MongoUrl = fmt.Sprintf("mongodb://%s:%s@%s:%s",
		MConfig.User,
		MConfig.Passwd,
		MConfig.Host,
		MConfig.Port,
	)

	OssClient, err = oss.New(
		MConfig.EndPoint,
		MConfig.AccessKeyId,
		MConfig.AccessKeySecret)
	if err != nil {
		panic(err)
	}
}
