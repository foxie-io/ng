package dtos

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type GetUserRequest struct {
	ID int `json:"id" param:"id" path:"id" validate:"required,gt=0"`
}

type DeleteUserRequest struct {
	ID int `json:"id" param:"id" path:"id" validate:"required,gt=0"`
}

type ListUsersRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}
