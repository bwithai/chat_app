package main

import (
	"chatapp/db"
	"chatapp/router"
	"chatapp/user"
	"chatapp/ws"
	"log"
)

func main() {

	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("could not initialize database connection: %s", err)
	}

	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	userHandler := user.NewHandler(userSvc)

	//roomRep := chat_room.NewRepository(dbConn.GetDB())
	//roomSvc := chat_room.NewService(roomRep)
	//roomHandler := chat_room.NewHandler(roomSvc)

	hub := ws.NewHub()
	r := ws.NewRepository(dbConn.GetDB())
	wsHandler := ws.NewHandler(hub, r)
	go hub.Run()

	router.InitRouter(userHandler, wsHandler)
	router.Start("8000")
}
