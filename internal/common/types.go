package common

import (
	"time"
)

// User struct
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	LastSeen  time.Time `json:"last_seen"`
}


// Chat struct
type Chat struct {
	ID        string    `json:"id"`
	Users     []User    `json:"users"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Message struct
type Message struct {
	ID        string    `json:"id"`
	ChatID    string    `json:"chat_id"`
	Author    User      `json:"author"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}


type UserApiService interface {
	CreateUser(user User) error
	GetUser(id string) (User, error)
	UpdateUser(id string, user User) error
	DeleteUser(id string) error
}

type ChatApiService interface {
	CreateChat(chat Chat) error
	GetChat(id string) (Chat, error)
	UpdateChat(id string, chat Chat) error
	DeleteChat(id string) error
}

type MessageApiService interface {
	CreateMessage(message Message) error
	GetMessages(id string) ([]Message, error)
	UpdateMessage(id string, message Message) error
	DeleteMessage(id string) error
}