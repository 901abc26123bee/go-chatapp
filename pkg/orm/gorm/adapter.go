package gormpsql

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gsm/pkg/orm"
)

// adapter defines the implementation for gorm to implement DB interface.
type adapter struct {
	db  *gorm.DB
	key string // key to encrypted sensitive data
}

// Wrap wraps a gorm db to orm DB.
func Wrap(db *gorm.DB) orm.DB {
	return &adapter{
		db: db,
	}
}

func (a *adapter) AutoMigrate(dst ...interface{}) error {
	return a.db.AutoMigrate(dst...)
}

func (a *adapter) WithContext(ctx context.Context) orm.DB {
	return &adapter{db: a.db.WithContext(ctx), key: a.key}
}
func (a *adapter) Table(name string, args ...interface{}) orm.DB {
	return &adapter{db: a.db.Table(name, args...), key: a.key}
}

func (a *adapter) Model(value interface{}) orm.DB {
	return &adapter{db: a.db.Model(value), key: a.key}
}

func (a *adapter) Create(value interface{}) orm.DB {
	return &adapter{db: a.db.Create(value), key: a.key}
}

func (a *adapter) CreateInBatches(value interface{}, batchSize int) orm.DB {
	return &adapter{db: a.db.CreateInBatches(value, batchSize), key: a.key}
}

func (a *adapter) Find(out interface{}, where ...interface{}) orm.DB {
	return &adapter{db: a.db.Find(out, where...), key: a.key}
}

func (a *adapter) Take(dest interface{}, conds ...interface{}) orm.DB {
	return &adapter{db: a.db.Take(dest, conds...), key: a.key}
}

func (a *adapter) Distinct(args ...interface{}) orm.DB {
	return &adapter{db: a.db.Distinct(args...), key: a.key}
}

func (a *adapter) Where(query interface{}, args ...interface{}) orm.DB {
	return &adapter{db: a.db.Where(query, args...), key: a.key}
}

func (a *adapter) Limit(limit int) orm.DB {
	return &adapter{db: a.db.Limit(limit), key: a.key}
}

func (a *adapter) OffSet(offset int) orm.DB {
	return &adapter{db: a.db.Offset(offset), key: a.key}
}

func (a *adapter) Order(value interface{}) orm.DB {
	return &adapter{db: a.db.Order(value), key: a.key}
}

func (a *adapter) Or(query interface{}, args ...interface{}) orm.DB {
	return &adapter{db: a.db.Or(query, args...), key: a.key}
}

func (a *adapter) Not(query interface{}, args ...interface{}) orm.DB {
	return &adapter{db: a.db.Not(query, args...), key: a.key}
}

func (a *adapter) First(out interface{}, where ...interface{}) orm.DB {
	return &adapter{db: a.db.First(out, where...), key: a.key}
}

func (a *adapter) Updates(value interface{}) orm.DB {
	return &adapter{db: a.db.Updates(value), key: a.key}
}

func (a *adapter) Joins(query string, args ...interface{}) orm.DB {
	return &adapter{db: a.db.Joins(query, args...), key: a.key}
}

func (a *adapter) Select(query interface{}, args ...interface{}) orm.DB {
	return &adapter{db: a.db.Select(query, args...), key: a.key}
}

func (a *adapter) Delete(value interface{}, where ...interface{}) orm.DB {
	return &adapter{db: a.db.Delete(value, where...), key: a.key}
}

func (a *adapter) Count(count *int64) orm.DB {
	return &adapter{db: a.db.Count(count), key: a.key}
}

func (a *adapter) ClauseReturning() orm.DB {
	return &adapter{db: a.db.Clauses(clause.Returning{}), key: a.key}
}

func (a *adapter) Exec(query string, args ...interface{}) orm.DB {
	return &adapter{db: a.db.Exec(query, args...), key: a.key}
}

func (a *adapter) Raw(query string, args ...interface{}) orm.DB {
	return &adapter{db: a.db.Raw(query, args...), key: a.key}
}

func (a *adapter) Scan(dest interface{}) orm.DB {
	return &adapter{db: a.db.Scan(dest), key: a.key}
}

func (a *adapter) DB() (*sql.DB, error) {
	return a.db.DB()
}

func (a *adapter) Error() error {
	return a.db.Error
}

func (a *adapter) RowsAffected() int64 {
	return a.db.RowsAffected
}

func (a *adapter) Transaction(fc func(tx orm.DB) error) error {
	fc2 := func(tx *gorm.DB) error {
		return fc(&adapter{db: tx, key: a.key})
	}
	return a.db.Transaction(fc2)
}

func (a *adapter) Debug() orm.DB {
	return &adapter{db: a.db.Debug(), key: a.key}
}
