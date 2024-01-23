package types

import (
	"github.com/google/uuid"
)

type Task struct {
	Id        int64  `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type AppResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type UserRegistration struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6,max=30"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6,max=30"`
}

type User struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type TokenPayload struct {
	UserId uuid.UUID `json:"user_id"`
}

type SignInPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
