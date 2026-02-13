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
	OTP        string     `json:"-"`
	OTPExpiry  int64      `json:"-"`
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
	Status    string         `json:"status"`
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

// 4. Struct untuk Delete Task
type DeleteTask struct {
	IDs []uint `json:"ids"`
}

// 5. Struct untuk Update Task
type UpdateTask struct {
	Title     string   `json:"title"`
	ShortDesc string   `json:"short_desc"`
	LongDesc  string   `json:"long_desc"`
	Priority  string   `json:"priority"`
	Status    string   `json:"status" gorm:"default:'todo'"`
	Time      string   `json:"time"`
	Date      string   `json:"date"`
	Tags      []string `json:"tags"`
}

// 6. Struct untuk Update Batch Status
type UpdateBatchStatus struct {
	IDs    []uint `json:"ids"`
	Status string `json:"status"`
}

// 7. Struct untuk Forgot Password
type ForgotPassword struct {
	Email string `json:"email"`
}

// 8. Struct untuk Reset Password
type ResetPassword struct {
	Email    string `json:"email"`
	OTP      string `json:"otp"`
	Password string `json:"password"`
}

// 9. Struct untuk Change Password
type ChangePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
