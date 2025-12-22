package dal

// data access layer

import (
	"context"
	"example/advanced/models"

	"github.com/foxie-io/gormqs"
	"gorm.io/gorm"
)

type (
	// rename to have same prefix
	// user data access object
	UserDao struct {
		gormqs.Queries[models.User, *UserDao]
		db    *gorm.DB
		model models.User
	}
)

func (qs *UserDao) DBInstance(ctx context.Context) *gorm.DB {
	db := gormqs.ContextValue(ctx, qs.db)
	return db.WithContext(ctx).Model(qs.model).Table(qs.model.TableName())
}

func NewUserDao(db *gorm.DB) *UserDao {
	qs := &UserDao{db: db}
	qs.Queries = gormqs.NewQueries[models.User](qs)
	return qs
}
