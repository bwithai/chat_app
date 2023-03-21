package ws

import (
	"gorm.io/gorm"
	"time"

	"github.com/gorilla/websocket"
)

type Service struct {
	Client
	repository
	timeout time.Duration
}

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       string `json:"id"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

type Message struct {
	messageID int64  `gorm:"primaryKey;autoIncrement"`
	Content   string `json:"content"`
	RoomID    string `json:"roomId"`
	UserID    string `json:"userId"`
	Username  string `json:"username"`
}

type DbMessage struct {
	gorm.Model
	messageID int64  `gorm:"primaryKey;autoIncrement"`
	Content   string `json:"content"`
	RoomID    string `json:"roomId"`
	UserID    string `json:"userId"`
	Username  string `json:"username"`
}

//func (s *Service) SaveMessage() {
//
//}

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		c.Conn.WriteJSON(message)
	}
}
