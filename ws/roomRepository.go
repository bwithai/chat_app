package ws

import (
	"chatapp/user"
	"context"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) saveMessage(ctx context.Context, message *DbMessage) (*DbMessage, error) {
	result := r.db.WithContext(ctx).Create(message)
	if result.Error != nil {
		return nil, result.Error
	}
	return message, nil
}

func (r *repository) finedMessagesByRoomID(roomId string) ([]DbMessage, error) {
	var messages []DbMessage
	result := r.db.Where("room_id = ?", roomId).Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}

	return messages, nil
}

func (r *repository) GetUserById(userId string) (*user.User, error) {
	var u user.User
	result := r.db.First(&u, userId)
	if result.Error != nil {
		return nil, result.Error
	}
	//if result.Error != nil {
	//	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
	//		return nil, fmt.Errorf("user not found")
	//	}
	//	return nil, result.Error
	//}
	return &u, nil
}

type Repository interface {
	saveMessage(ctx context.Context, message *DbMessage) (*DbMessage, error)
	finedMessagesByRoomID(roomId string) ([]DbMessage, error)
	GetUserById(userId string) (*user.User, error)
	//CreateUser(ctx context.Context, user *User) (*User, error)
	//GetUserByEmail(ctx context.Context, email string) (*User, error)
}
