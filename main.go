package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"

	search "github.com/Kirisakiii/kuroko/proto"
	"github.com/Kirisakiii/kuroko/server"
	"github.com/yanyiwu/gojieba"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
)

var (
	mongoClient *mongo.Client
	splitter    *gojieba.Jieba
)

type Config struct {
	MongoURL string `json:"mongo_url"`
}

func init() {
	var config Config
	// read config file
	configData, err := os.Open("./config.json")
	if err != nil {
		fmt.Println("failed to open config file: ", err)
		panic(err)
	}
	defer configData.Close()

	// decode config file
	decoder := json.NewDecoder(configData)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("failed to decode config file: ", err)
		panic(err)
	}

	mongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MongoURL))
	if err != nil {
		fmt.Println("failed to connect to mongo: ")
		panic(err)
	}
	err = mongoClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		fmt.Println("failed to ping mongo: ")
		panic(err)
	}

	splitter = gojieba.NewJieba()
}

func main() {
	listen, err := net.Listen("tcp", "localhost:5016")
	if err != nil {
		fmt.Println("failed to listen: ", err)
	}
	fmt.Println("server listening on localhost:5016")

	s := grpc.NewServer()

	collection := mongoClient.Database("kuroko").Collection("indexes")

	search.RegisterSearchEngineServer(s, server.NewSearchEngine(collection, splitter))
	if err := s.Serve(listen); err != nil {
		fmt.Println("failed to serve: ", err)
	}
}
