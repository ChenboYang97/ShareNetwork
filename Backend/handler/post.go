package handler // controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"ShareNetwork/model"
	"ShareNetwork/service"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
)

var (
	mediaTypes = map[string]string{
		".jpeg": "image",
		".jpg":  "image",
		".gif":  "image",
		".png":  "image",
		".mov":  "video",
		".mp4":  "video",
		".avi":  "video",
		".flv":  "video",
		".wmv":  "video",
	}
)

// func uploadHandler(w http.ResponseWriter, r *http.Request) {
// 	// Parse from body of request to get a json object.
// 	fmt.Println("Received one post request")
// 	decoder := json.NewDecoder(r.Body)
// 	var p model.Post
// 	if err := decoder.Decode(&p); err != nil {
// 		panic(err) // 过于极端，因为直接重启服务器了
// 	}

// 	fmt.Fprintf(w, "Post received: %s\n", p.Message)
// }
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one upload request")

	// 从token中找到现在操作的username，之前输入username不合理
	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"]

	p := model.Post{
		Id:      uuid.New(),
		User:    username.(string),
		Message: r.FormValue("message"),
	}

	file, header, err := r.FormFile("media_file")
	if err != nil {
		http.Error(w, "Media file is not available", http.StatusBadRequest)
		fmt.Printf("Media file is not available %v\n", err)
		return
	}

	suffix := filepath.Ext(header.Filename)
	if t, ok := mediaTypes[suffix]; ok { // map 有对应的ok就是true，otherwise是false
		p.Type = t
	} else {
		p.Type = "unknown"
	}

	err = service.SavePost(&p, file)
	if err != nil {
		http.Error(w, "Failed to save post to backend", http.StatusInternalServerError)
		fmt.Printf("Failed to save post to backend %v\n", err)
		return
	}

	fmt.Println("Post is saved successfully.")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one search request")

	user := r.URL.Query().Get("user")
	keywords := r.URL.Query().Get("keywords")

	var posts []model.Post
	var err error
	if user != "" {
		posts, err = service.SeachPostsByUser(user)
	} else {
		posts, err = service.SeachPostsByKeywords(keywords)
	}
	if err != nil {
		http.Error(w, "Failed to read data from Elasticsearch", http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Failed to get json data from search result", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one request for delete")

	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"].(string)
	id := mux.Vars(r)["id"]

	if err := service.DeletePost(id, username); err != nil {
		http.Error(w, "Failed to delete post from backend", http.StatusInternalServerError)
		fmt.Printf("Failed to delete post from backend %v\n", err)
		return
	}
	fmt.Println("Post is deleted successfully")
}
