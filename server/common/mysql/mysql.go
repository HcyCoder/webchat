package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

func New(dsn string) sqlx.SqlConn {
	return sqlx.NewMysql(dsn)
}
