package main

import (
	"fmt"
	"log"
	"net/http"

	"ShareNetwork/backend"
	"ShareNetwork/handler"
	"ShareNetwork/util"
)

func main() {
	fmt.Println("started-service")

	config, err := util.LoadApplicationConfig("conf", "deploy.yml")
	if err != nil {
		panic(err)
	}

	backend.InitElasticSearchBackend(config.ElasticsearchConfig)
	backend.InitGCSBackend(config.GCSConfig)

	log.Fatal(http.ListenAndServe(":8080", handler.InitRouter(config.TokenConfig)))
}
