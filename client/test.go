package main

import (
	"context"
	"fmt"

	search "github.com/Kirisakiii/kuroko/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:5016", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := search.NewSearchEngineClient(conn)
	_, err = client.Search(context.Background(), &search.SearchRequest{})
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// _, err = client.CreatePostIndex(ctx, &search.CreatePostIndexRequest{
	// 	Id:      3,
	// 	Title:   "test4",
	// 	Content: "大部分的生活都乏味得不值一提，根本就没有不乏味的时候。换另一种牌子的香烟也好，搬到一个新地方去住也好，订阅别的报纸也好，坠入爱河又脱身出来也好，我们一直在以或轻浮或深沉的方式，来对抗日常生活那无法消释的乏味成分。",
	// })
	// if err != nil {
	// 	panic(err)
	// }

	result, err := client.Search(ctx, &search.SearchRequest{
		Query: "没有只要",
	})

	fmt.Println(result)

	if err != nil {
		panic(err)
	}
}
