package dalo

import (
	"fmt"

	"github.com/foxie-io/gormqs"
	"gorm.io/gorm"
)

type Column struct {
	Value string
}

func (col Column) String() string {
	return col.Value
}

func (col Column) Asc() string {
	return fmt.Sprintf("%s ASC", col.Value)
}

func (col Column) Desc() string {
	return fmt.Sprintf("%s DESC", col.Value)
}

func (col Column) Eq(value any) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		column := gormqs.WithTable(string(col.Value), db)
		return db.Where(fmt.Sprintf("%s = ?", column), value)
	}
}

func (col Column) Gt(value any) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		column := gormqs.WithTable(string(col.Value), db)
		return db.Where(fmt.Sprintf("%s > ?", column), value)
	}
}

func (col Column) Gte(value any) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		column := gormqs.WithTable(string(col.Value), db)
		return db.Where(fmt.Sprintf("%s >= ?", column), value)
	}
}

func (col Column) Lt(value any) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		column := gormqs.WithTable(string(col.Value), db)
		return db.Where(fmt.Sprintf("%s < ?", column), value)
	}
}

func (col Column) Lte(value any) gormqs.Option {
	return func(db *gorm.DB) *gorm.DB {
		column := gormqs.WithTable(string(col.Value), db)
		return db.Where(fmt.Sprintf("%s <= ?", column), value)
	}
}
