// Package database	Provides a unified API interface for
// relational and non-relational databases.
//
// For relational databases, MySQL is integrated by default,
// and PostgreSQL, SQLite, Oracle, etc. will be integrated later.
// For non-relational databases, Redis and MongoDB are integrated by default,
// and will be integrated according to actual use later.
//
// A simple integration example
//
//	type (
//		MysqlDB     *sqlx.DB
//		MysqlOption func(MysqlDB)
//	)
//
//	func InitMysql(dsn string) MysqlDB {
//	  	if dsn == "" {
//	  		dsn = _default_dsn
//	  	}
//	  	ctx, cancle := context.WithTimeout(context.Background(), _conn_timeout)
//	  	defer cancle()
//	  	db, err := sqlx.ConnectContext(ctx, _driver_mysql, dsn)
//	  	if err != nil {
//			panic("init mysql database instance fail.")
//	  	}
//	  	return db
//	}
//
//	func SetupMysql(mdb MysqlDB, opts ...MysqlOption) {
//	    for _, fn := range opts {
//	    	fn(mdb)
//	    }
//	}
//
// The most regrettable aspect is that it's impossible to bind all public APIs to each individual database type.
// This is a limitation of Go, and it's also why we need to pass the database instance.
// Since some public APIs don't always require generic types, we have to simplify the design here,
// but we've still achieved a highly templated code structure.
package database

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wendisx/puzzle/pkg/clog"
)

const (
	_driver_mysql  = "mysql"
	_driver_sqlite = "sqlite3"

	_default_mysql_dsn  = "<user>:<password>@<proto>(<host>:<port>)[/<db_name>][?options]"
	_default_sqlite_dsn = "file:<db_name>[?{key=value&}]"
	_conn_sql_timeout   = 3 * time.Second
)

// InsertWithPlace return error occurred during the execution of the insert SQL with placeholder parameters.
// The primary key that returns a successful insert may be modified later.
// The database instance needs to be explicitly specified.
func InsertWithPlace(ctx context.Context, db *sqlx.DB, sqlStr string, args ...any) error {
	_, err := db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		clog.Error(err.Error())
		return err
	}
	return nil
}

// UpdateWithPlace return error occurred during the execution of the update SQL with placeholder parameters.
// The database instance needs to be explicitly specified.
func UpdateWithPlace(ctx context.Context, db *sqlx.DB, sqlStr string, args ...any) error {
	_, err := db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		clog.Error(err.Error())
		return err
	}
	return nil
}

// DeleteWithPlace return error occurred during the execution of the delete SQL with placeholdler parameters.
// The database instance needs to be explicitly specified.
func DeleteWithPlace(ctx context.Context, db *sqlx.DB, sqlStr string, args ...any) error {
	_, err := db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		clog.Error(err.Error())
		return err
	}
	return nil
}

// QueryWithPlace return the specify generic type and error occurred during the execution of the select SQL with placeholdler parameters.
// R should have the largest set of all fields that need to be retrieved and not be a pointer type.
// The database instance needs to be explicitly specified.
func QueryWithPlace[R any](ctx context.Context, db *sqlx.DB, sqlStr string, args ...any) (R, error) {
	var dest R
	err := db.GetContext(ctx, &dest, sqlStr, args...)
	if err != nil {
		clog.Error(err.Error())
	}
	return dest, err
}

// QListWithPlace return the list of specify generic type and error occurred during the execution of the select SQL with placeholdler parameters.
// R should have the largest set of all fields that need to be retrieved and not be a pointer type.
// The database instance needs to be explicitly specified.
func QListWithPlace[R any](ctx context.Context, db *sqlx.DB, sqlStr string, args ...any) ([]R, error) {
	var dest []R
	err := db.SelectContext(ctx, &dest, sqlStr, args...)
	if err != nil {
		clog.Error(err.Error())
	}
	return dest, err
}

// QPageWithPlace return the page of specify generic type and error occurred during the execution of the select SQL with placeholdler parameters.
// R should have the largest set of all fields that need to be retrieved and not be a pointer type.
// The database instance needs to be explicitly specified.
func QPageWithPlace[R any](ctx context.Context, db *sqlx.DB, sqlStr string, args ...any) (Page[R], error) {
	var page Page[R]
	list := make([]R, 0)
	list, err := QListWithPlace[R](ctx, db, sqlStr, args...)
	if err != nil {
		clog.Error(err.Error())
	}
	page.Items = list
	return page, err
}

// InsertWithName return error occurred during the execution of the insert SQL with named parameters.
// The primary key that returns a successful insert may be modified later.
// The database instance needs to be explicitly specified.
func InsertWithName(ctx context.Context, db *sqlx.DB, sqlStr string, obj any) error {
	_, err := db.NamedExecContext(ctx, sqlStr, obj)
	if err != nil {
		clog.Error(err.Error())
		return err
	}
	return nil
}

// UpdateWithName return error occurred during the execution of the udpate SQL with named parameters.
// The database instance needs to be explicitly specified.
func UpdateWithName(ctx context.Context, db *sqlx.DB, sqlStr string, obj any) error {
	_, err := db.NamedExecContext(ctx, sqlStr, obj)
	if err != nil {
		clog.Error(err.Error())
		return err
	}
	return nil
}

// DeleteWithName return error occurred during the execution of the delete SQL with named parameters.
// The database instance needs to be explicitly specified.
func DeleteWithName(ctx context.Context, db *sqlx.DB, sqlStr string, obj any) error {
	_, err := db.NamedExecContext(ctx, sqlStr, obj)
	if err != nil {
		clog.Error(err.Error())
		return err
	}
	return nil
}

// QueryWithName return the specify generic type and error occurred during the execution of the select SQL with named parameters.
// R should have the largest set of all fields that need to be retrieved and not be a pointer type.
// The database instance needs to be explicitly specified.
func QueryWithName[R any](ctx context.Context, db *sqlx.DB, sqlStr string, obj any) (R, error) {
	var dest R
	rows, err := db.NamedQueryContext(ctx, sqlStr, obj)
	if err == nil {
		for rows.Next() {
			err = rows.StructScan(&dest)
		}
	}
	if err != nil {
		clog.Error(err.Error())
	}
	return dest, err
}

// QListWithName return the list of specify generic type and error occurred during the execution of the select SQL with named parameters.
// R should have the largest set of all fields that need to be retrieved and not be a pointer type.
// The database instance needs to be explicitly specified.
func QListWithName[R any](ctx context.Context, db *sqlx.DB, sqlStr string, obj any) ([]R, error) {
	dest := make([]R, 0)
	rows, err := db.NamedQueryContext(ctx, sqlStr, obj)
	var row R
	for rows.Next() {
		err = rows.StructScan(&row)
		if err != nil {
			clog.Error(err.Error())
			continue
		}
		dest = append(dest, row)
	}
	if err != nil {
		clog.Error(err.Error())
	}
	return dest, err
}

// QPageWithName return the page of specify generic type and error occurred during the execution of the select SQL with named parameters.
// R should have the largest set of all fields that need to be retrieved and not be a pointer type.
// The database instance needs to be explicitly specified.
func QPageWithName[R any](ctx context.Context, db *sqlx.DB, sqlStr string, obj any) (Page[R], error) {
	var page Page[R]
	list := make([]R, 0)
	list, err := QListWithName[R](ctx, db, sqlStr, obj)
	if err != nil {
		clog.Error(err.Error())
	}
	page.Items = list
	return page, err
}
