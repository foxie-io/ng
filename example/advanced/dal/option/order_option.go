package dalo

import (
	"github.com/foxie-io/gormqs"
	"gorm.io/gorm"
)

type (
	orderColumn struct {
		Column
	}

	orderEntity struct {
		ID       orderColumn
		UserID   orderColumn
		Quantity orderColumn
		Product  orderColumn
	}
)

func newOrderColumn(name string) orderColumn {
	return orderColumn{
		Column: Column{name},
	}
}

var ORDERS = orderEntity{
	ID:       newOrderColumn("id"),
	UserID:   newOrderColumn("user_id"),
	Quantity: newOrderColumn("quantity"),
	Product:  newOrderColumn("product"),
}

func (s orderEntity) Select(cols ...orderColumn) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		columns := make([]string, len(cols))
		for i, col := range cols {
			columns[i] = gormqs.WithTable(string(col.String()), db)
		}
		return db.Select(columns)
	}
}
