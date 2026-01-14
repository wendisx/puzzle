package database

import (
	"context"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/wendisx/puzzle/pkg/clog"
)

type (
	// mysql database from sqlx
	MysqlDB *sqlx.DB
	// mysql database instance functional configuration
	MysqlOption func(MysqlDB)
)

// InitMysql return a new mysql database instance with specific data source name.
// It does not check the validity of the dsn and will panic if dsn is invalid.
func InitMysql(dsn string) MysqlDB {
	if dsn == "" {
		dsn = _default_mysql_dsn
	}
	ctx, cancle := context.WithTimeout(context.Background(), _conn_sql_timeout)
	defer cancle()
	db, err := sqlx.ConnectContext(ctx, _driver_mysql, dsn)
	if err != nil {
		clog.Panic(err.Error())
	} else {
		clog.Info("Mysql Database initialization successful.")
	}
	return db
}

func SetupMysql(db MysqlDB, opts ...MysqlOption) {
	for _, fn := range opts {
		fn(db)
	}
}
