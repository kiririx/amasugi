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

func (*AmiRepo[T]) GetTableName() string {
	var t T
	return t.GetTableName()
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

func (ar *AmiRepo[T]) Insert(t *T) (int64, error) {
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
	sqlStr := fmt.Sprintf("insert into %s %s values %v", ar.GetTableName(), sqlColumn, sqlValues)
	result, err := db.Exec(sqlStr, values...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (ar *AmiRepo[T]) Update(t *T) (int64, error) {
	tType := reflect.TypeOf(*t)
	sqlColumn := ""
	values := make([]any, 0, tType.NumField())
	for i := 0; i < tType.NumField(); i++ {
		field := tType.Field(i)
		value := reflect.ValueOf(*t).Field(i)
		tag := field.Tag.Get(constx.TAG)
		if tag == "id" {
			continue
		}
		sqlColumn += fmt.Sprintf("%s=?", tag)
		values = append(values, ReflectValParse(value))
		if i < tType.NumField()-1 {
			sqlColumn += ","
		}
	}
	id := reflect.ValueOf(*t).FieldByName("Id").Int()
	sqlStr := fmt.Sprintf("update %s set %s where id = %v", ar.GetTableName(), sqlColumn, id)
	result, err := db.Exec(sqlStr, values...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (*AmiRepo[T]) UpdateColumns(columns []string, t T) {

}

func (ar *AmiRepo[T]) DeleteById(id uint64) (int64, error) {
	sqlVal := fmt.Sprintf("delete from %s where id = ?", ar.GetTableName())
	result, err := db.Exec(sqlVal, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (ar *AmiRepo[T]) Delete(sql string, args ...any) (int64, error) {
	sqlVal := fmt.Sprintf("delete from %s where %v", ar.GetTableName(), sql)
	result, err := db.Exec(sqlVal, args)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (*AmiRepo[T]) Execute(sql string, args ...any) (map[string]any, error) {
	return nil, nil
}
