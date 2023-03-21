package user

import (
	"context"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Update("user_id", u.ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &u, nil
}
