package models

import (
	"errors"
	"time"
)

var ErrNoRecordedSnippet = errors.New("models: подходящей записи не найдено")
var ErrSnippetExpired = errors.New("models: запись просрочена")

type Snippet struct {
	ID        int       `gorm:"primarykey"`
	Title     string    `gorm:"type:varchar(100)"`
	Content   string    `gorm:"type:text"`
	Created   time.Time `gorm:"type:datetime"`
	Expires   time.Time `gorm:"type:datetime"`
	CreatedBy int
}
type User struct {
	UserID int `gorm:"primarykey"`
	ChatID int64
}
