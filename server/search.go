package server

import (
	"context"
	"slices"

	search "github.com/Kirisakiii/kuroko/proto"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *SearchEngine) Search(ctx context.Context, req *search.SearchRequest) (*search.SearchResponse, error) {
	var keywords []string
	searchContent := req.GetQuery()

	keywords = append(keywords, s.splitter.CutForSearch(searchContent, true)...)

	type match struct {
		id    int64
		score int64
	}

	matches := make(map[int64]match)

	for _, keyword := range keywords {
		filter := bson.M{
			"keyword": keyword,
		}
		r := s.collection.FindOne(ctx, filter)
		if r.Err() != nil {
			continue
		}

		var result struct {
			PostIds []int64 `bson:"post_ids"`
		}
		err := r.Decode(&result)
		if err != nil {
			continue
		}

		for _, id := range result.PostIds {
			if _, ok := matches[id]; !ok {
				matches[id] = match{
					id:    id,
					score: 0,
				}
			}
			matches[id] = match{
				id:    id,
				score: matches[id].score + 1,
			}
		}
	}

	// sort matches by score
	var matchesSlice []match
	for _, m := range matches {
		matchesSlice = append(matchesSlice, m)
	}
	slices.SortFunc(matchesSlice, func(i, j match) int {
		if i.score > j.score {
			return -1
		}
		if i.score < j.score {
			return 1
		}
		return 0
	})

	var result []int64
	for _, m := range matchesSlice {
		result = append(result, m.id)
	}

	return &search.SearchResponse{
		Ids: result,
	}, nil
}
