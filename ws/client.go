package ws

import (
	"context"
	"gorm.io/gorm"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type service struct {
	*Client
	Repository
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

func (c *Client) readMessage(hub *Hub, cnt context.Context, msgChan chan<- *Message) {
	//ctx, cancel := context.WithTimeout(cnt, time.Duration())
	//defer cancel()
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &Message{
			Content:  string(m),
			UserID:   c.ID,
			RoomID:   c.RoomID,
			Username: c.Username,
		}

		hub.Broadcast <- msg
		msgChan <- msg
	}
}
