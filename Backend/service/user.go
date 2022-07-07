package service

import (
	"ShareNetwork/backend"
	"ShareNetwork/constants"
	"ShareNetwork/model"
	"fmt"

	"github.com/olivere/elastic/v7"
)

func CheckUser(username, password string) (bool, error) {

	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("username", username))
	query.Must(elastic.NewTermQuery("password", password))

	searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
	if err != nil {
		return false, err
	}

	if searchResult.TotalHits() > 0 {
		fmt.Printf("Login as %s\n", username)
		return true, nil
	}

	return false, nil
}

// 将用户添加到数据库中
// 如果添加失败
// 1. 用户问题，返回false
// 2. 后台问题，返回error
func AddUser(user *model.User) (bool, error) {
	query := elastic.NewTermQuery("username", user.Username)
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
	if err != nil {
		return false, err
	}

	// 如果该用户名已存在，返回false
	if searchResult.TotalHits() > 0 {
		return false, nil
	}

	// 将该用户存入到elastic search数据库中
	err = backend.ESBackend.SaveToES(user, constants.USER_INDEX, user.Username)
	if err != nil {
		return false, err
	}
	fmt.Printf("User is added: %s\n", user.Username)
	return true, nil
}
