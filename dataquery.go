package amasugi

import (
	"database/sql"
	"fmt"
	"github.com/kiririx/krutils/sugar"
	"reflect"
	"time"
)

type DataQuery[T IModel] struct {
	stat   *sql.Rows
	sqlVal *SQLVal
	offset uint64
	limit  uint64
}

type baseType interface {
	int | int64 | int16 | int8 | int32 | string | bool
}

type SQLVal struct {
	sql    string
	params []any
}

var tagM = make(map[string]string)

func (dq *DataQuery[T]) Next() (T, bool) {
	var t T
	ok := false
	dq.offset = dq.offset + dq.limit
	dq.limit = 10
	dq.sqlVal.sql += fmt.Sprintf(" limit %v,%v", dq.offset, dq.limit)
	if dq.stat.Next() {
		tt := reflect.TypeOf(t)
		sugar.ForIndex(tt.NumField(), func(i int) (bool, bool) {
			field := tt.Field(i)
			tag := field.Tag.Get("ami")
			tagM[tag] = field.Name
			return false, false
		})

		_cols := make([]any, 0)
		// dbModel := make(map[string]any)
		// todo 这里应该用map，更快
		_colsName := make([]string, 0)
		colTypes, err := dq.stat.ColumnTypes()
		for _, v := range colTypes {
			colName := v.Name()
			_colsName = append(_colsName, colName)
			colType := v.DatabaseTypeName()
			switch colType {
			case "INT":
				vv := 0
				_cols = append(_cols, &vv)
			case "VARCHAR":
				vv := ""
				_cols = append(_cols, &vv)
			case "DATETIME":
				var vv time.Time
				_cols = append(_cols, &vv)
			default:
				vv := ""
				_cols = append(_cols, &vv)
			}

		}
		err = dq.stat.Scan(_cols...)
		if err != nil {
			panic(err)
		}

		tElem := reflect.ValueOf(&t).Elem()
		for i, v := range _colsName {
			if _, ok := tagM[v]; !ok {
				continue
			}
			tVal := tElem.FieldByName(tagM[v])
			// todo 这里直接set bool。。更快
			tVal.Set(reflect.ValueOf(_cols[i]).Elem())
		}
		ok = true
	}
	return t, ok
}

func (dq *DataQuery[T]) Read(f func(t T)) {
	for t, ok := dq.Next(); ok; {
		f(t)
	}
}
