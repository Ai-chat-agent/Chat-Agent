package models

import (
	"time"
)

// ChatMessage represents a chat message in the system
type ChatMessage struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"not null;index"`
	Message   string    `json:"message" gorm:"not null"`
	Timestamp string    `json:"timestamp"`
	IsBot     bool      `json:"is_bot" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ChatMessageRequest represents the request structure for sending a message
type ChatMessageRequest struct {
	UserID  string `json:"user_id" binding:"required"`
	Message string `json:"message" binding:"required,min=1,max=1000"`
}

// ChatMessageResponse represents the response structure after sending a message
type ChatMessageResponse struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
}

// User represents a user in the system
type User struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ChatSession represents a chat session
type ChatSession struct {
	ID        string        `json:"id" gorm:"primaryKey"`
	UserID    string        `json:"user_id" gorm:"not null;index"`
	Title     string        `json:"title"`
	IsActive  bool          `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Messages  []ChatMessage `json:"messages,omitempty" gorm:"foreignKey:UserID;references:UserID"`
}

// TableName returns the table name for ChatMessage
func (ChatMessage) TableName() string {
	return "chat_messages"
}

// TableName returns the table name for User
func (User) TableName() string {
	return "users"
}

// TableName returns the table name for ChatSession
func (ChatSession) TableName() string {
	return "chat_sessions"
}

