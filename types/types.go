package types

import "database/sql"

type Task struct {
	Id int `json:"id"`
	Text string `json:"text"`
	Completed bool `json:"completed"`
	UserID sql.NullString `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type AppResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

type UserRegistration struct {
	Id int `json:"id" `
	Name string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=30"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=6,max=30"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type User struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}