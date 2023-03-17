package router

import (
	"chatapp/user"
	"fmt"
	"log"

	"github.com/gorilla/mux"
	"net/http"
)

var r = mux.NewRouter()

func InitRouter(userHandler *user.Handler) {

	r.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "Server is Up and running...")
	})

	r.HandleFunc("/signup", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/logout", userHandler.Logout).Methods("GET")
}

func Start(addr string) {
	log.Println("Server listening on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
