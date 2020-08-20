package convert

import (
	"log"
	"reflect"

	"github.com/lovego/strs"
	"github.com/lovego/structs"
)

func ValidateQuery(typ reflect.Type) {
	if typ.Kind() != reflect.Struct {
		log.Panic("req.Query must be a struct.")
	}
}

func Query(value reflect.Value, map2strs map[string][]string) (err error) {
	structs.Traverse(value, true, func(v reflect.Value, f reflect.StructField) bool {
		if f.Tag.Get("json") == "-" {
			return true
		}

		var lowercaseName string

		values := map2strs[f.Name]
		if len(values) == 0 {
			lowercaseName = strs.FirstLetterToLower(f.Name)
			values = map2strs[lowercaseName]
		}
		if len(values) == 0 {
			switch f.Type.Kind() {
			case reflect.Slice, reflect.Array:
				name := f.Name + "[]"
				if values = map2strs[name]; len(values) == 0 {
					values = map2strs[lowercaseName+"[]"]
				}
			}
		}
		if len(values) > 0 {
			err = SetArray(v, values)
		}
		return err == nil // if err == nil, go on Traverse
	})
	return
}