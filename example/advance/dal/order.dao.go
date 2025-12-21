package dal

// data access layer

import (
	"context"
	"example/advance/models"

	"github.com/foxie-io/gormqs"
	"gorm.io/gorm"
)

type (
	// order data access object
	OrderDao struct {
		gormqs.Queries[models.Order, *OrderDao]
		db    *gorm.DB
		model models.Order
	}
)

func (qs *OrderDao) DBInstance(ctx context.Context) *gorm.DB {
	db := gormqs.ContextValue(ctx, qs.db)
	return db.WithContext(ctx).Model(qs.model).Table(qs.model.TableName())
}

func NewOrderDao(db *gorm.DB) *OrderDao {
	qs := &OrderDao{db: db}
	qs.Queries = gormqs.NewQueries[models.Order](qs)
	return qs
}
