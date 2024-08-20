package setup

import (
	"math"
	"reflect"
	"strconv"
	"time"
)

const (
	KindMax = reflect.UnsafePointer + 1
)

type FieldSetter interface {
	IsZero(value reflect.Value) bool
	Set(value reflect.Value, val string) error
}

type Field struct {
	Setter   FieldSetter
	Value    reflect.Value
	TagValue string
}

func (f Field) IsZero() bool {
	return f.Setter.IsZero(f.Value)
}

type floatField struct {
	size int
}

func (f floatField) IsZero(value reflect.Value) bool {
	return math.Float64bits(value.Float()) == 0
}

func (f floatField) Set(value reflect.Value, val string) error {
	if val, err := strconv.ParseFloat(val, f.size); err == nil {
		value.SetFloat(val)
	}
	return nil
}

type boolField struct {
}

func (b boolField) IsZero(value reflect.Value) bool {
	return !value.Bool()
}

func (b boolField) Set(value reflect.Value, val string) error {
	if val, err := strconv.ParseBool(val); err == nil {
		value.SetBool(val)
	}
	return nil
}

type intField struct {
	size int
}

func (i intField) IsZero(value reflect.Value) bool {
	return value.Int() == 0
}

func (i intField) Set(value reflect.Value, val string) error {
	if i.size == 64 {
		if val, err := time.ParseDuration(val); err == nil {
			value.Set(reflect.ValueOf(val).Convert(value.Type()))
			return nil
		}
	}
	if val, err := strconv.ParseInt(val, 0, i.size); err == nil {
		value.SetInt(val)
	}
	return nil
}

type uintField struct {
	size int
}

func (u uintField) IsZero(value reflect.Value) bool {
	return value.Uint() == 0
}

func (u uintField) Set(value reflect.Value, val string) error {
	if val, err := strconv.ParseUint(val, 0, u.size); err == nil {
		value.SetUint(val)
	}
	return nil
}

type stringField struct {
}

func (s stringField) IsZero(value reflect.Value) bool {
	return value.String() == ""
}

func (s stringField) Set(value reflect.Value, val string) error {
	value.SetString(val)
	return nil
}

type sliceField struct {
}

func (s sliceField) IsZero(value reflect.Value) bool {
	return value.Len() == 0
}

func (s sliceField) Set(value reflect.Value, val string) error {
	value.Set(reflect.MakeSlice(value.Type(), 0, 0))
	return nil
}

type emptyField struct {
}

func (e emptyField) IsZero(value reflect.Value) bool {
	return true
}

func (e emptyField) Set(value reflect.Value, val string) error {
	return nil
}

type pointerField struct {
}

func (p pointerField) IsZero(value reflect.Value) bool {
	return value.IsNil()
}

func (p pointerField) Set(value reflect.Value, val string) error {
	value.Set(reflect.New(value.Type().Elem()))
	return nil
}

type chanField struct {
}

func (c chanField) IsZero(value reflect.Value) bool {
	return value.IsNil()
}

func (c chanField) Set(value reflect.Value, val string) error {
	// TODO: channel capacity
	value.Set(reflect.MakeChan(value.Type(), 0))
	return nil
}

type mapField struct {
}

func (m mapField) IsZero(value reflect.Value) bool {
	return value.Len() == 0
}

func (m mapField) Set(value reflect.Value, val string) error {
	// TODO: map capacity
	value.Set(reflect.MakeMapWithSize(value.Type(), 0))
	return nil
}

type funcField struct {
}

func (f funcField) IsZero(value reflect.Value) bool {
	return value.IsNil()
}

func (f funcField) Set(value reflect.Value, val string) error {
	fn := func(args []reflect.Value) (results []reflect.Value) {
		return
	}
	value.Set(reflect.MakeFunc(value.Type(), fn))
	return nil
}

type interfaceField struct {
}

func (i interfaceField) IsZero(value reflect.Value) bool {
	return value.IsNil()
}

func (i interfaceField) Set(value reflect.Value, val string) error {
	// TODO implement me
	panic("implement me")
}

type arrayField struct {
}

func (a arrayField) IsZero(value reflect.Value) bool {
	return value.Len() == 0
}

func (a arrayField) Set(value reflect.Value, val string) error {
	return nil
}

type structField struct {
}

func (s structField) IsZero(value reflect.Value) bool {
	for i := 0; i < value.NumField(); i++ {
		if !value.Field(i).IsZero() {
			return false
		}
	}
	return true
}

func (s structField) Set(value reflect.Value, val string) error {
	// TODO implement me
	panic("implement me")
}

type complexField struct {
}

func (c complexField) IsZero(value reflect.Value) bool {
	complex := value.Complex()
	return math.Float64bits(real(complex)) == 0 && math.Float64bits(imag(complex)) == 0
}

func (c complexField) Set(value reflect.Value, val string) error {
	// TODO implement me
	panic("implement me")
}

var (
	fieldSetters [KindMax]FieldSetter
)

func GetSetter(kind reflect.Kind) FieldSetter {
	return fieldSetters[kind]
}

func init() {
	fieldSetters[reflect.Invalid] = &emptyField{}
	fieldSetters[reflect.Bool] = &boolField{}
	// 2
	fieldSetters[reflect.Int] = &intField{size: strconv.IntSize}
	fieldSetters[reflect.Int8] = &intField{size: 8}
	fieldSetters[reflect.Int16] = &intField{size: 16}
	fieldSetters[reflect.Int32] = &intField{size: 32}
	fieldSetters[reflect.Int64] = &intField{size: 64}
	// 7
	fieldSetters[reflect.Uint] = &uintField{size: strconv.IntSize}
	fieldSetters[reflect.Uint8] = &uintField{size: 8}
	fieldSetters[reflect.Uint16] = &uintField{size: 16}
	fieldSetters[reflect.Uint32] = &uintField{size: 32}
	fieldSetters[reflect.Uint64] = &uintField{size: 64}
	// 12
	fieldSetters[reflect.Uintptr] = &intField{size: strconv.IntSize}
	// 13
	fieldSetters[reflect.Float32] = &floatField{size: 32}
	fieldSetters[reflect.Float64] = &floatField{size: 64}
	// 15
	fieldSetters[reflect.Complex64] = &complexField{}
	fieldSetters[reflect.Complex128] = &complexField{}
	// 17
	fieldSetters[reflect.Array] = &arrayField{}
	fieldSetters[reflect.Chan] = &chanField{}
	fieldSetters[reflect.Func] = &funcField{}
	fieldSetters[reflect.Interface] = &interfaceField{}
	fieldSetters[reflect.Map] = &mapField{}
	fieldSetters[reflect.Pointer] = &pointerField{}
	fieldSetters[reflect.Slice] = &sliceField{}
	fieldSetters[reflect.String] = &stringField{}
	fieldSetters[reflect.Struct] = &structField{}
	fieldSetters[reflect.UnsafePointer] = &pointerField{}
	// 27
}
