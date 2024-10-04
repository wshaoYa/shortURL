package sequence

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	replaceSql = `REPLACE INTO sequence (stub) VALUES ('a')`
)

// MySQL 发号器
type MySQL struct {
	conn sqlx.SqlConn
}

// NewMySQL 构造函数
func NewMySQL(dsn string) *MySQL {
	return &MySQL{
		conn: sqlx.NewMysql(dsn),
	}
}

// Next 取号
func (m *MySQL) Next() (uint64, error) {
	//预编译sql
	stmt, err := m.conn.Prepare(replaceSql)
	if err != nil {
		logx.Errorw("sequence mysql m.conn.Prepare failed", logx.Field("err", err))
		return 0, err
	}
	defer stmt.Close()

	//执行sql
	result, err := stmt.Exec()
	if err != nil {
		logx.Errorw("sequence mysql stmt.Exec() failed", logx.Field("err", err))
		return 0, err
	}
	//取last id
	lastID, err := result.LastInsertId()
	if err != nil {
		logx.Errorw("sequence mysql result.LastInsertId() failed", logx.Field("err", err))
		return 0, err
	}

	return uint64(lastID), nil
}
