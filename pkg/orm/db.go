package orm

import (
	"context"
	"database/sql"
)

// DB defines the interface for postgres orm.
type DB interface {
	// AutoMigrate run auto migration for given models
	AutoMigrate(dst ...interface{}) error

	// WithContext change current instance db's context to ctx
	WithContext(ctx context.Context) DB

	// Table specify the table you would like to run db operations
	Table(name string, args ...interface{}) DB

	// Model specify the model you would like to run db operations
	Model(value interface{}) DB

	// Create insert the value into database
	Create(value interface{}) DB

	// CreateInBatches insert the value in batches of batchSize
	CreateInBatches(value interface{}, batchSize int) DB

	// Find find records that match given conditions
	Find(out interface{}, where ...interface{}) DB

	// Take finds the first record returned by the database in no specified order, matching given conditions conds
	Take(dest interface{}, conds ...interface{}) DB

	// Distinct specify distinct fields that you want querying
	Distinct(args ...interface{}) DB

	// Where find records that match given conditions
	Where(query interface{}, args ...interface{}) DB

	// Limit specify the number of records to be retrieved
	Limit(limit int) DB

	// Offset specify the number of records to skip before starting to return the records
	OffSet(offset int) DB

	// Order specify order when retrieve records from database
	Order(value interface{}) DB

	// Or add OR conditions
	Or(query interface{}, args ...interface{}) DB

	// Not add NOT conditions
	Not(query interface{}, args ...interface{}) DB

	// First find first record that match given conditions, order by primary key
	First(out interface{}, where ...interface{}) DB

	// Updates update attributes with callbacks
	Updates(value interface{}) DB

	//Joins specify Joins conditions, gorm default joins type is left join
	Joins(query string, args ...interface{}) DB

	// Select select attributes with callbacks
	Select(query interface{}, args ...interface{}) DB

	// Delete delete value match given conditions, if the value has primary key, then will including the primary key as condition
	Delete(value interface{}, where ...interface{}) DB

	// Count get matched records count
	Count(count *int64) DB

	// Clauses Add returning clauses.
	ClauseReturning() DB

	// Exec executes a query. The args are for any placeholder parameters in the query.
	Exec(query string, args ...interface{}) DB

	// Raw executes a query. The args are for any placeholder parameters in the query.
	Raw(query string, args ...interface{}) DB

	// Scan scans selected value to the struct dest
	Scan(dest interface{}) DB

	// Get the *sql.DB
	DB() (*sql.DB, error)

	// Error returns current DB error.
	Error() error

	// RowsAffected returns current DB affected rows.
	RowsAffected() int64

	// Transaction start a transaction as a block, return error will rollback, otherwise to commit.
	Transaction(fc func(tx DB) error) error

	// Debug start debug mode
	Debug() DB
}
