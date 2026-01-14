package database

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wendisx/puzzle/pkg/clog"
)

type (
	// Simple encapsulation for sqlite database from sqlx.
	SqliteDB     *sqlx.DB
	SqliteOption func(SqliteDB)
)

func InitSqlite(dsn string) SqliteDB {
	if dsn == "" {
		dsn = _default_sqlite_dsn
	}
	ctx, cancle := context.WithTimeout(context.Background(), _conn_sql_timeout)
	defer cancle()
	db, err := sqlx.ConnectContext(ctx, _driver_sqlite, dsn)
	if err != nil {
		clog.Panic(err.Error())
	} else {
		clog.Info("Sqlite Database initialization successful.")
	}
	return db
}

func SetupSqlite(db SqliteDB, opts ...SqliteOption) {
	for _, fn := range opts {
		fn(db)
	}
}
