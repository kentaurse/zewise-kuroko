package server

import (
	"context"

	search "github.com/Kirisakiii/kuroko/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *SearchEngine) CreatePostIndex(ctx context.Context, req *search.CreatePostIndexRequest) (*search.CreatePostIndexResponse, error) {
	var keywords []string

	title := req.GetTitle()
	content := req.GetContent()

	keywords = append(keywords, s.splitter.CutForSearch(title, true)...)
	keywords = append(keywords, s.splitter.CutForSearch(content, true)...)

	for _, keyword := range keywords {
		filter := bson.M{
			"keyword": keyword,
		}
		update := bson.M{
			"$addToSet": bson.M{
				"post_ids": req.GetId(),
			},
		}
		_, err := s.collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		if err != nil {
			return &search.CreatePostIndexResponse{
				Code: 1,
			}, err
		}
	}

	return &search.CreatePostIndexResponse{
		Code: 0,
	}, nil
}
