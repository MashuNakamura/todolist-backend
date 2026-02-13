package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// 0. Return Message
type Ret struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   int    `json:"code"`
	Data    any    `json:"data"`
}

// 1. Tabel Users
type User struct {
	gorm.Model
	Name       string     `json:"name"`
	Email      string     `json:"email" gorm:"unique"`
	Password   string     `json:"-"`
	Tasks      []Task     `json:"tasks"`
	Categories []Category `json:"categories"`
}

// 2. Tabel Tasks (Todolist)
type Task struct {
	gorm.Model
	Title     string         `json:"title"`
	ShortDesc string         `json:"short_desc"`
	LongDesc  string         `json:"long_desc"`
	Priority  string         `json:"priority"`
	Status    bool           `json:"status"`
	Time      string         `json:"time"`
	Date      string         `json:"date"`
	Tags      pq.StringArray `json:"tags" gorm:"type:text[]"`
	UserID    uint           `json:"user_id"`
}

// 3. Tabel Categories (Label Warna)
type Category struct {
	gorm.Model
	Name   string `json:"name"`
	Color  string `json:"color"`
	UserID uint   `json:"user_id"`
}
