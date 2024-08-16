package setup

import (
	"reflect"
)

const (
	KindMax = reflect.UnsafePointer + 1
)

type Field struct {
	Type     reflect.Kind
	Value    reflect.Value
	TagValue string
	Func     func(value reflect.Value, val string) error
}

var (
	fieldEmpties = [KindMax]Field{
		reflect.Bool: {
			Type: reflect.Bool,
			Func: func(value reflect.Value, val string) error {
				switch val {
				case "1", "t", "T", "true", "TRUE", "True":
					value.SetBool(true)
				default:
					value.SetBool(false)
				}
				return nil
			},
		},
	}
)

func (f Field) IsEmpty() bool {
	switch f.Type {
	case reflect.Bool:
		return !f.Value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return f.Value.Int() == 0
	case reflect.Float32, reflect.Float64:
		return !(f.Value.Float() != 0)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return f.Value.Uint() == 0
	case reflect.Slice:
		return f.Value.Len() == 0
	case reflect.String:
		return f.Value.String() == ""
	case reflect.Pointer:
		return f.Value.IsNil()
	}
	return true
}

func (f Field) IsZero() bool {
	switch f.Type {
	case reflect.Bool:
		return !f.Value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return f.Value.Int() == 0
	case reflect.Float32, reflect.Float64:
		return !(f.Value.Float() != 0)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return f.Value.Uint() == 0
	case reflect.Slice:
		return f.Value.Len() == 0
	case reflect.String:
		return f.Value.String() == ""
	case reflect.Pointer:
		return f.Value.IsNil()
	case reflect.Map:
		return f.Value.Len() == 0
	case reflect.Chan:
		return f.Value.IsNil()
	default:

	}
	return true
}
