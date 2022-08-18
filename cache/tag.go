package cache

import (
	"github.com/kiririx/amasugi/constx"
	"github.com/kiririx/amasugi/model"
	"github.com/kiririx/krutils/sugar"
	"reflect"
)

var TagM = make(map[string]string)

func InitTagM[T model.IModel](t T) {
	tt := reflect.TypeOf(t)
	sugar.ForIndex(tt.NumField(), func(i int) (bool, bool) {
		field := tt.Field(i)
		tag := field.Tag.Get(constx.TAG)
		TagM[tag] = field.Name
		return false, false
	})
}
