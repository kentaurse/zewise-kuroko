package server

import (
	search "github.com/Kirisakiii/kuroko/proto"
	"github.com/yanyiwu/gojieba"
	"go.mongodb.org/mongo-driver/mongo"
)

type SearchEngine struct {
	search.UnimplementedSearchEngineServer
	collection *mongo.Collection
	splitter   *gojieba.Jieba
}

func NewSearchEngine(collection *mongo.Collection, splitter *gojieba.Jieba) *SearchEngine {
	return &SearchEngine{
		collection: collection,
		splitter:   splitter,
	}
}
