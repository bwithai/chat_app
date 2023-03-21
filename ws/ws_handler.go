package ws

import (
	"chatapp/auth_jwt"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	Repository
	hub *Hub
}

func NewHandler(h *Hub, r Repository) *Handler {
	return &Handler{
		hub:        h,
		Repository: r,
	}
}

type CreateRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateRoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req CreateRoomReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	resp := &CreateRoomRes{
		ID:   req.ID,
		Name: req.Name,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) checkRoomId(roomId string) bool {
	for _, r := range h.hub.Rooms {
		if r.ID == roomId {
			return true
		}
	}
	// If no room is found.
	return false
}

func (h *Handler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	roomID := mux.Vars(r)["roomId"]
	clientID := mux.Vars(r)["userID"]
	tokenStr := r.Header.Get("Token")
	if tokenStr == "" {
		err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Add JWT Token at http header request")))
		defer conn.Close()
		if err != nil {
			log.Println(err)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//clientID := r.URL.Query().Get("userId")

	status, _ := auth_jwt.VaalidateJWT(tokenStr, clientID)

	if status != "authorized" {
		err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Un Authorized user_id: %v", clientID)))
		defer conn.Close()
		if err != nil {
			log.Println(err)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if h.checkRoomId(roomID) != true {
		err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Room Id: %v is not created", roomID)))
		defer conn.Close()
		if err != nil {
			log.Println(err)
			return
		}
		return
	}

	user, err3 := h.Repository.GetUserById(clientID)
	if err3 != nil {
		err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("User Id: %v not Registered", clientID)))
		defer conn.Close()
		if err != nil {
			log.Println(err)
			return
		}
		return
	}

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       strconv.Itoa(int(user.ID)),
		RoomID:   roomID,
		Username: user.Username,
	}

	m := &Message{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		UserID:   strconv.Itoa(int(user.ID)),
		Username: user.Username,
	}

	h.hub.Register <- cl
	h.hub.Broadcast <- m

	go cl.writeMessage()
	msgChan := make(chan *Message)
	go cl.readMessage(h.hub, r.Context(), msgChan)
	message := <-msgChan

	//message := cl.readMessage(h.hub, r.Context())

	msg := &DbMessage{
		Content:  message.Content,
		RoomID:   message.RoomID,
		UserID:   message.UserID,
		Username: message.Username,
	}

	_, err2 := h.Repository.saveMessage(r.Context(), msg)
	if err2 != nil {
		return
	}
}

type RoomRes struct {
	ID    string      `json:"id"`
	Name  string      `json:"name"`
	Users []ClientRes `json:"users"`
}

func (h *Handler) GetRoom(w http.ResponseWriter, r *http.Request) {
	roomId := mux.Vars(r)["roomId"]
	for _, r := range h.hub.Rooms {
		if r.ID == roomId {
			clients := h.GetClientsForRoom(r.ID)
			room := RoomRes{
				ID:    r.ID,
				Name:  r.Name,
				Users: clients,
			}
			json.NewEncoder(w).Encode(room)
			return
		}
	}

	// If no room is found, return a 404 error.
	w.WriteHeader(http.StatusNotFound)
}

func (h *Handler) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms := make([]RoomRes, 0)

	for _, r := range h.hub.Rooms {
		clients := h.GetClientsForRoom(r.ID)
		rooms = append(rooms, RoomRes{
			ID:    r.ID,
			Name:  r.Name,
			Users: clients,
		})
	}

	json.NewEncoder(w).Encode(rooms)
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *Handler) GetClientsForRoom(roomID string) []ClientRes {
	var clients []ClientRes

	if _, ok := h.hub.Rooms[roomID]; !ok {
		return clients
	}

	for _, c := range h.hub.Rooms[roomID].Clients {
		clients = append(clients, ClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	return clients
}

type CreateMassageRes struct {
	Text    string
	UserID  string
	RoomID  string
	CreatAt time.Time
}

func (h *Handler) GetRoomMessages(w http.ResponseWriter, r *http.Request) {
	roomId := mux.Vars(r)["roomId"]

	messages, err := h.Repository.finedMessagesByRoomID(roomId)
	if err != nil {
		log.Printf("error: %v", err)
	}

	var res []*CreateMassageRes
	for _, m := range messages {
		res = append(res, &CreateMassageRes{
			Text:    m.Content,
			RoomID:  m.RoomID,
			UserID:  m.UserID,
			CreatAt: m.CreatedAt,
		})
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}
