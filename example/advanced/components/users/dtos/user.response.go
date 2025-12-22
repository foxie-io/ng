package dtos

import (
	"github.com/foxie-io/gormqs"
)

type CreateUserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (g CreateUserResponse) TableName() string {
	return "users"
}

type GetUserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (g GetUserResponse) TableName() string {
	return "users"
}

type GetAllUsersResponse struct {
	gormqs.ListResulter[GetUserResponse]
}

type UpdateUserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type DeleteUserResponse struct {
	Success bool `json:"success"`
}
