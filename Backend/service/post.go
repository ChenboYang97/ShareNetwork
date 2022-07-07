package service

import (
	"mime/multipart"
	"reflect"

	"ShareNetwork/backend"
	"ShareNetwork/constants"
	"ShareNetwork/model"

	"github.com/olivere/elastic/v7"
)

func SeachPostsByUser(user string) ([]model.Post, error) {
	query := elastic.NewTermQuery("user", user)
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
	if err != nil {
		return nil, err
	}
	return getPostFromSearchResult(searchResult), nil
}

func SeachPostsByKeywords(keywords string) ([]model.Post, error) {
	query := elastic.NewMatchQuery("message", keywords) // vincent + rick + sun, elastic中用的是+
	query.Operator("and")
	if keywords == "" {
		query.ZeroTermsQuery("all") // if keywords are null, return all result
	}

	searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
	if err != nil {
		return nil, err
	}
	return getPostFromSearchResult(searchResult), nil
}

func getPostFromSearchResult(searchResult *elastic.SearchResult) []model.Post {
	var ptype model.Post
	var posts []model.Post

	for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {
		p := item.(model.Post)
		posts = append(posts, p)
	}
	return posts
}

func SavePost(post *model.Post, file multipart.File) error {
	medialink, err := backend.GCSBackend.SaveToGCS(file, post.Id)
	if err != nil {
		return err
	}
	post.Url = medialink

	return backend.ESBackend.SaveToES(post, constants.POST_INDEX, post.Id)
}

func DeletePost(id string, user string) error {
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("id", id))
	query.Must(elastic.NewTermQuery("user", user))

	return backend.ESBackend.DeleteFromES(query, constants.POST_INDEX)
}
