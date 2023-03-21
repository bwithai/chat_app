package router

import (
	"chatapp/auth_jwt"
	"chatapp/user"
	"chatapp/ws"
	"fmt"
	"log"

	"github.com/gorilla/mux"
	"net/http"
)

var r = mux.NewRouter()

func InitRouter(userHandler *user.Handler, wsHandler *ws.Handler) {

	r.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "Server is Up and running...")
	})

	r.HandleFunc("/api/register", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/api/login", userHandler.Login).Methods("POST")

	r.Handle("/api/users/{userID}/profile/logout", auth_jwt.ValidateJWT(http.HandlerFunc(userHandler.Logout))).Methods("GET")

	r.Handle("/api/users/{userID}/profile/createRoom", auth_jwt.ValidateJWT(http.HandlerFunc(wsHandler.CreateRoom))).Methods("POST")
	r.Handle("/api/users/{userID}/profile/chat/rooms", auth_jwt.ValidateJWT(http.HandlerFunc(wsHandler.GetRooms))).Methods("GET")
	r.Handle("/api/users/{userID}/profile/chat/rooms/{roomId}", auth_jwt.ValidateJWT(http.HandlerFunc(wsHandler.GetRoom))).Methods("GET")

	r.HandleFunc("/api/users/{userID}/profile/chat/room/{roomId}/messages", wsHandler.JoinRoom)

	r.Handle("/api/users/{userID}/profile/chat/rooms/{roomId}/messages", auth_jwt.ValidateJWT(http.HandlerFunc(wsHandler.GetRoomMessages))).Methods("GET")

}

func Start(addr string) {
	log.Println("Server listening on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
