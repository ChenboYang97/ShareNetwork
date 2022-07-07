package handler

import (
	"net/http"

	"ShareNetwork/util"

	jwtmiddleware "github.com/auth0/go-jwt-middleware" // 中间件
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var mySigningKey []byte

func InitRouter(config *util.TokenInfo) http.Handler {
	mySigningKey = []byte(config.Secret)

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	router := mux.NewRouter()

	// these two methods need to have middleware
	router.Handle("/upload", jwtMiddleware.Handler(http.HandlerFunc(uploadHandler))).Methods("POST")
	router.Handle("/search", jwtMiddleware.Handler(http.HandlerFunc(searchHandler))).Methods("GET")
	router.Handle("/post/{id}", jwtMiddleware.Handler(http.HandlerFunc(deleteHandler))).Methods("DELETE")

	router.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST")
	router.Handle("/signin", http.HandlerFunc(signinHandler)).Methods("POST")

	// 表明一下跨域访问的情况下，哪些操作被允许
	originsOk := handlers.AllowedOrigins([]string{"*"})                             // *表示支持所有访问地址
	headersOk := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}) // 支持authorization，contenttype等header
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "DELETE"})         //支持get,post,delete等操作

	return handlers.CORS(originsOk, headersOk, methodsOk)(router)
}
