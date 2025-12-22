package dtos

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func (c CreateUserRequest) TableName() string {
	return "users"
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type PathID struct {
	ID int `param:"id" validate:"required,gt=0"`
}

type ListUsersRequest struct {
	Page     int `query:"page" default:"1" validate:"gte=1"`
	PageSize int `query:"pageSize" default:"10" validate:"gte=1,lte=100"`

	// selects can be "count", "list" or "count,list"
	// select to get total count, list data or both
	Selects string `query:"selects" default:"count,list"`
}
