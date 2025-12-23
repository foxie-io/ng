package dalo

import (
	"github.com/foxie-io/gormqs"
	"gorm.io/gorm"
)

type (
	userColumn struct {
		Column
	}

	userEntity struct {
		ID    userColumn
		Name  userColumn
		Email userColumn
	}
)

func newUserColumn(name string) userColumn {
	return userColumn{
		Column: Column{name},
	}
}

var USERS = userEntity{
	ID:    newUserColumn("id"),
	Name:  newUserColumn("name"),
	Email: newUserColumn("email"),
}

func (s userEntity) Select(cols ...userColumn) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		columns := make([]string, len(cols))
		for i, col := range cols {
			columns[i] = gormqs.WithTable(string(col.String()), db)
		}
		return db.Select(columns)
	}
}
