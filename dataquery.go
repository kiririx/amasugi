package amasugi

import (
	"database/sql"
	"fmt"
	"github.com/kiririx/amasugi/cache"
	"github.com/kiririx/amasugi/model"
	"github.com/kiririx/krutils/logx"
	"reflect"
	"time"
)

type DataQuery[T model.IModel] struct {
	stat   *sql.Rows
	sqlVal *SQLVal
	offset int64
	limit  int64
	pos    int64
}

// type baseType interface {
// 	int | int64 | int16 | int8 | int32 | string | bool
// }

type SQLVal struct {
	sql    string
	params []any
}

func (dq *DataQuery[T]) reGetResultSet() error {
	dq.offset = dq.offset + dq.limit
	dq.limit = 10
	dq.sqlVal.sql += fmt.Sprintf(" limit %v,%v", dq.offset, dq.limit)
	rows, err := db.Query(dq.sqlVal.sql, dq.sqlVal.params...)
	if err != nil {
		logx.ERR(err)
		return err
	}
	dq.stat = rows
	return nil
}

func (dq *DataQuery[T]) toPage(pageNum, pageSize uint64) {
	dq.offset = int64(pageNum * pageSize)
	dq.limit = int64(pageSize)
	dq.pos = dq.offset
}

// Next dataQuery的迭代方法，每调用一次，就向下返回一条记录，第一次读取时，从数据库里取出limit条数据，但是只返回结果集的第一条。
//
// 第二次读取，如果已读数量小于limit，就继续读取结果集的第二条，以此类推，假如limit是10，第11次读取的时候，结果集已经读完了，那么就从数据库再取limit条，然后继续读取。
func (dq *DataQuery[T]) Next() (*T, error) {
	var t T
	if dq.pos >= dq.offset+dq.limit-1 {
		err := dq.reGetResultSet()
		if err != nil {
			return &t, err
		}
	}
	if !dq.stat.Next() {
		return nil, nil
	}

	dq.pos++
	cache.InitTagM(t)

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
		return nil, err
	}

	tElem := reflect.ValueOf(&t).Elem()
	for i, v := range _colsName {
		if _, ok := cache.TagM[v]; !ok {
			continue
		}
		tVal := tElem.FieldByName(cache.TagM[v])
		// todo 这里直接set bool。。更快
		tVal.Set(reflect.ValueOf(_cols[i]).Elem())
	}
	return &t, nil
}

func (dq *DataQuery[T]) Read(f func(t *T, err error)) {
	dq.ReadLimit(f, -1)
}

func (dq *DataQuery[T]) ReadLimit(f func(t *T, err error), limit int64) {
	var i int64 = 0
	for {
		if limit > 0 && i >= limit {
			break
		}
		t, err := dq.Next()
		if err != nil {
			f(nil, err)
			break
		}
		f(t, nil)
		i++
	}
}

func (dq *DataQuery[T]) Count() (int64, error) {
	sqlCond := fmt.Sprintf("select count(*) from (%v)", dq.sqlVal.sql)
	rows, err := db.Query(sqlCond, dq.sqlVal.params...)
	if err != nil {
		return 0, err
	}
	if !rows.Next() {
		return 0, err
	}
	var c int64 = 0
	err = rows.Scan(&c)
	if err != nil {
		return 0, err
	}
	return c, nil
}
