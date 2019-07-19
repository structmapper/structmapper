package structmapper

import (
	"reflect"

	"github.com/pkg/errors"
)

func StringerModule(m Mapper) {
	// *.String() -> string
	m.RegisterTransformerFunc(
		func(target Target) bool {
			return target.From.AssignableTo(stringerType) && target.To.AssignableTo(stringType)
		},
		func(from reflect.Value, _ reflect.Type) (reflect.Value, error) {
			if from.Type().Kind() == reflect.Ptr && from.IsNil() {
				return reflect.ValueOf(""), nil
			}

			str, ok := from.Interface().(stringer)
			if !ok {
				return reflect.ValueOf(""), errors.New("Invalid value type")
			}

			return reflect.ValueOf(str.String()), nil
		},
	)
}

type stringer interface {
	String() string
}

var (
	stringerType = reflect.TypeOf((*stringer)(nil)).Elem()
	stringType   = reflect.TypeOf("")
)
