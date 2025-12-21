package users

import (
	"example/advance/components/users/dtos"
	"example/advance/models"
	"sync"

	nghttp "github.com/foxie-io/ng/http"
)

type UserService struct {
	mu       sync.Mutex
	users    map[int]*models.User
	userList []*models.User
}

func NewUserService() *UserService {
	return &UserService{
		users:    make(map[int]*models.User),
		userList: []*models.User{},
	}
}

func (s *UserService) CreateUser(req dtos.CreateUserRequest) (*dtos.CreateUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := len(s.users) + 1
	user := &models.User{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	}
	s.users[id] = user
	s.userList = append(s.userList, user)
	return &dtos.CreateUserResponse{
		ID:    id,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *UserService) GetUser(id int) (*dtos.GetUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[id]
	if !exists {
		return nil, nghttp.NewErrNotFound()
	}
	return &dtos.GetUserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *UserService) GetAllUsers(dto *dtos.ListUsersRequest) *dtos.GetAllUsersResponse {
	s.mu.Lock()
	defer s.mu.Unlock()

	start := (dto.Page - 1) * dto.PageSize
	end := start + dto.PageSize

	if start > len(s.userList) {
		return &dtos.GetAllUsersResponse{Users: []dtos.GetUserResponse{}}
	}

	if end > len(s.userList) {
		end = len(s.userList)
	}

	var users []dtos.GetUserResponse
	for _, user := range s.userList[start:end] {
		users = append(users, dtos.GetUserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}
	return &dtos.GetAllUsersResponse{Users: users}
}

func (s *UserService) UpdateUser(id int, req *dtos.UpdateUserRequest) (*dtos.UpdateUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[id]; !exists {
		return nil, nghttp.NewErrNotFound()
	}

	updatedUser := &models.User{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	}
	s.users[id] = updatedUser

	for i, user := range s.userList {
		if user.ID == id {
			s.userList[i] = updatedUser
			break
		}
	}

	return &dtos.UpdateUserResponse{
		ID:    updatedUser.ID,
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
	}, nil
}

func (s *UserService) DeleteUser(params *dtos.DeleteUserRequest) *dtos.DeleteUserResponse {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[params.ID]; !exists {
		return &dtos.DeleteUserResponse{Success: false}
	}
	delete(s.users, params.ID)

	for i, user := range s.userList {
		if user.ID == params.ID {
			s.userList = append(s.userList[:i], s.userList[i+1:]...)
			break
		}
	}
	return &dtos.DeleteUserResponse{Success: true}
}
