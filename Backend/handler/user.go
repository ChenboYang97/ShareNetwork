package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"ShareNetwork/model"
	"ShareNetwork/service"

	jwt "github.com/form3tech-oss/jwt-go" // jwt可以作为这个lib名字
)

func signinHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one signin request")
	w.Header().Set("Content.Type", "text/plain")

	if r.Method == "OPTIONS" {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var user model.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Failed to read json data from request", http.StatusBadRequest)
		return
	}

	exist, err := service.CheckUser(user.Username, user.Password)
	if err != nil {
		http.Error(w, "Failed to read user from Elasticsearch", http.StatusInternalServerError)
		return
	}
	if !exist {
		http.Error(w, "User doesn't exist or wrong password", http.StatusUnauthorized)
		return
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tokenString))
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one signup request")

	decoder := json.NewDecoder(r.Body)
	var user model.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Failed to read json data from request", http.StatusBadRequest)
		return
	}

	// 这里其实放到前端判断更好
	if user.Username == "" || user.Password == "" || regexp.MustCompile(`^[a-z0-9]$`).MatchString(user.Username) {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		fmt.Printf("Invalid username or password\n")
		return
	}

	success, err := service.AddUser(&user)
	if err != nil {
		http.Error(w, "Failed to save user to backend", http.StatusInternalServerError)
		return
	}
	if !success {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	fmt.Println("User is added succesfully")
}
