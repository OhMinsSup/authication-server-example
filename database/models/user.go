package models

import (
	"github.com/OhMinsSup/lafu-server/lib"
	"time"
)

// User 유저 모델
type User struct {
	ID          string `gorm:"primary_key;uuid"`
	Email       string `gorm:"size:255;unique_index"`
	Username    string `gorm:"size:255;unique_index"`
	Thumbnail   string `gorm:"size:255"`
	Password    string `gorm:"size:255"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func (user *User) Serialize() lib.JSON {
	return lib.JSON{
		"id": user.ID,
		"username": user.Username,
		"email": user.Email,
		"thumbnail": user.Thumbnail,
	}
}

func (user *User) TokenData(tokenID string) lib.JSON {
	return lib.JSON{
		"id": user.ID,
		"username": user.Username,
		"email": user.Email,
		"thumbnail": user.Thumbnail,
		"token": tokenID,
	}
}
