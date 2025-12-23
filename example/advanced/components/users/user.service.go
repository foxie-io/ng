package users

import (
	"context"
	"errors"
	"example/advanced/components/users/dtos"
	"example/advanced/dal"
	"example/advanced/models"
	"sync"

	. "example/advanced/dal/option"

	"github.com/foxie-io/gormqs"
	nghttp "github.com/foxie-io/ng/http"
	"gorm.io/gorm"
)

type UserService struct {
	mu      sync.Mutex
	userDao *dal.UserDao
}

func NewUserService(userDao *dal.UserDao) *UserService {
	return &UserService{

		userDao: userDao,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req dtos.CreateUserRequest) (*dtos.CreateUserResponse, error) {
	record := &models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := s.userDao.CreateOne(ctx, record); err != nil {
		return nil, err
	}

	return &dtos.CreateUserResponse{
		ID:    record.ID,
		Name:  record.Name,
		Email: record.Email,
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, id int) (*dtos.GetUserResponse, error) {
	var (
		record = new(dtos.GetUserResponse)
	)

	if err := s.userDao.GetOneTo(ctx, record, USERS.ID.Eq(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nghttp.NewErrNotFound().Update(nghttp.Meta("entity", "User"))
		}
		return nil, err
	}

	return record, nil
}

func (s *UserService) GetAllUsers(ctx context.Context, dto *dtos.ListUsersRequest) (*gormqs.ListResulter[dtos.GetUserResponse], error) {

	var (
		limit  = dto.PageSize
		offset = (dto.Page - 1) * dto.PageSize
	)

	resulter := gormqs.NewListResulter[dtos.GetUserResponse](dto.Selects)
	if err := s.userDao.GetListTo(ctx, resulter, gormqs.LimitAndOffset(limit, offset)); err != nil {
		return nil, err
	}

	return resulter, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int, req *dtos.UpdateUserRequest) (*dtos.UpdateUserResponse, error) {
	if _, err := s.userDao.GetOne(ctx, USERS.ID.Eq(id), USERS.Select(USERS.ID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nghttp.NewErrNotFound().Update(nghttp.Meta("entity", "User"))
		}
		return nil, err
	}

	updatedUser := &models.User{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	}

	if _, err := s.userDao.Update(ctx, updatedUser, gormqs.WhereID(id)); err != nil {
		return nil, err
	}

	return &dtos.UpdateUserResponse{
		ID:    updatedUser.ID,
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) (*dtos.DeleteUserResponse, error) {
	if _, err := s.userDao.GetOne(ctx, USERS.ID.Eq(id), USERS.Select(USERS.ID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nghttp.NewErrNotFound().Update(nghttp.Meta("entity", "User"))
		}
		return nil, err
	}

	if _, err := s.userDao.Delete(ctx, gormqs.WhereID(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nghttp.NewErrNotFound().Update(nghttp.Meta("entity", "User"))
		}
	}

	return &dtos.DeleteUserResponse{
		Success: true,
	}, nil
}
