package amasugi

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kiririx/amasugi/constx"
	"github.com/kiririx/amasugi/model"
	"github.com/kiririx/krutils/confx"
	"github.com/kiririx/krutils/sugar"
	"reflect"
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

type AmiRepo[T model.IModel] struct {
}

func (*AmiRepo[T]) Query(sql string, args ...any) *DataQuery[T] {
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
		pos: -1,
	}
}

func ReflectValParse(v reflect.Value) any {
	switch v.Kind() {
	// reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Uint,reflect.Uint8,reflect.Int
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint16, reflect.Uint64:
		return v.Uint()
	case reflect.String:
		return v.String()
	case reflect.Struct:
		return v.Interface()
	default:
		panic("can not convert")
	}
}

func (*AmiRepo[T]) Insert(t *T) error {
	tType := reflect.TypeOf(*t)
	sqlColumn := `(`
	sqlValues := "("
	values := make([]any, 0, tType.NumField())
	sugar.ForIndex(tType.NumField(), func(i int) (bool, bool) {
		field := tType.Field(i)
		value := reflect.ValueOf(*t).Field(i)
		sqlColumn += field.Tag.Get(constx.TAG)
		values = append(values, ReflectValParse(value))
		sqlValues += "?"
		if i < tType.NumField()-1 {
			sqlColumn += ","
			sqlValues += ","
		}
		return false, false
	})
	sqlColumn += ")"
	sqlValues += ")"
	tableName := (*t).GetTableName()
	sqlStr := fmt.Sprintf("insert into %s %s values %v", tableName, sqlColumn, sqlValues)
	_, err := db.Exec(sqlStr, values...)
	if err != nil {
		return err
	}
	return nil
}

func (*AmiRepo[T]) Update(t *T) (uint64, error) {
	return 0, nil
}

func (*AmiRepo[T]) DeleteById(id uint64) error {
	return nil
}

func (*AmiRepo[T]) Delete(sql string, args ...any) (uint64, error) {
	return 0, nil
}

func (*AmiRepo[T]) Execute(sql string, args ...any) (map[string]any, error) {
	return nil, nil
}
