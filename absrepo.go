package amasugi

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kiririx/krutils/confx"
	"time"
)

var db *sql.DB

func init() {
	conf, err := confx.ResolveProperties("./config.properties")
	if err != nil {
		panic("config file not found")
	}
	db, err = sql.Open("mysql", conf["database.mysql.url"])
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}

type IModel interface {
	GetTableName() string
}

type AbsRepo[T IModel] struct {
}

func (*AbsRepo[T]) Insert(t T) {
}

func (*AbsRepo[T]) Query(sql string, args ...any) *DataQuery[T] {
	var t T
	sql = fmt.Sprintf("select * from %s where %v", t.GetTableName(), sql)
	// sql = fmt.Sprintf("select * from %s where %v limit %v,%v", t.GetTableName(), sql, 0, 10)
	// rows, err := db.Query(sql, args...)
	// if err != nil {
	// 	return nil
	// }

	return &DataQuery[T]{
		sqlVal: &SQLVal{
			sql:    sql,
			params: args,
		},
	}
}
