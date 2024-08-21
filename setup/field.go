package setup

import (
	"encoding/json"
	"math"
	"reflect"
	"strconv"
	"time"
)

const (
	KindMax = reflect.UnsafePointer + 1
)

type (
	Kind  = reflect.Kind
	Value = reflect.Value
)

type FieldSetter interface {
	IsZero(value Value) bool
	Init(value Value, val string) (Value, bool)
	Set(value Value, val string)
}

type Field struct {
	Setter   FieldSetter
	Value    Value
	TagValue string
}

func (f Field) IsZero() bool {
	return f.Setter.IsZero(f.Value)
}

func (f Field) Init() (Value, bool) {
	return f.Setter.Init(f.Value, f.TagValue)
}
func (f Field) Fill() bool {
	if !f.IsZero() {
		return false
	}
	f.Setter.Set(f.Value, f.TagValue)
	return true
}

func (f Field) CanSet() bool {
	return f.Value.CanSet()
}

func (f Field) Kind() Kind {
	return f.Value.Kind()
}

func (f Field) Set(convert Value) {
	f.Value.Set(convert)
}

type floatField struct {
	size int
}

func (f floatField) Init(value Value, val string) (Value, bool) {
	return value, false
}

func (f floatField) IsZero(value Value) bool {
	return math.Float64bits(value.Float()) == 0
}

func (f floatField) Set(value Value, val string) {
	if val, err := strconv.ParseFloat(val, f.size); err == nil {
		value.SetFloat(val)
	}
}

type boolField struct{}

func (b boolField) Init(value Value, val string) (Value, bool) {
	return value, false
}

func (b boolField) IsZero(value Value) bool {
	return !value.Bool()
}

func (b boolField) Set(value Value, val string) {
	if val, err := strconv.ParseBool(val); err == nil {
		value.SetBool(val)
	}
}

type intField struct {
	size int
}

func (i intField) Init(value Value, val string) (Value, bool) {
	return value, false
}

func (i intField) IsZero(value Value) bool {
	return value.Int() == 0
}

func (i intField) Set(value Value, val string) {
	if i.size == 64 {
		if val, err := time.ParseDuration(val); err == nil {
			value.Set(reflect.ValueOf(val))
			return
		}
	}
	if val, err := strconv.ParseInt(val, 0, i.size); err == nil {
		value.SetInt(val)
	}
}

type uintField struct {
	size int
}

func (u uintField) Init(value Value, val string) (Value, bool) {
	return value, false
}

func (u uintField) IsZero(value Value) bool {
	return value.Uint() == 0
}

func (u uintField) Set(value Value, val string) {
	if val, err := strconv.ParseUint(val, 0, u.size); err == nil {
		value.SetUint(val)
	}
}

type stringField struct{}

func (s stringField) Init(value Value, val string) (Value, bool) {
	return value, false
}

func (s stringField) IsZero(value Value) bool {
	return value.String() == ""
}

func (s stringField) Set(value Value, val string) {
	value.SetString(val)
}

type sliceField struct{}

func (s sliceField) Init(value Value, val string) (Value, bool) {
	ref := reflect.New(value.Type())
	ref.Elem().Set(reflect.MakeSlice(value.Type(), 0, 0))
	return ref, true
}

func (s sliceField) IsZero(value Value) bool {
	return value.Len() == 0
}

func (s sliceField) Set(value Value, val string) {
	if val != "" && val != "[]" {
		if err := json.Unmarshal([]byte(val), value.Interface()); err != nil {
		}
	}
	value.Set(value.Elem().Convert(value.Type()))
}

type emptyField struct{}

func (e emptyField) Init(value Value, val string) (Value, bool) {
	return value, false
}

func (e emptyField) IsZero(value Value) bool {
	return true
}

func (e emptyField) Set(value Value, val string) {
}

type pointerField struct{}

func (p pointerField) Init(value Value, val string) (Value, bool) {
	ref := reflect.New(value.Type().Elem())
	return ref, true
}

func (p pointerField) IsZero(value Value) bool {
	return value.IsNil()
}

func (p pointerField) Set(value Value, val string) {
	value.Set(reflect.New(value.Type().Elem()))
}

type chanField struct{}

func (c chanField) Init(value Value, val string) (Value, bool) {
	ref := reflect.MakeChan(value.Type(), 0)
	value.Set(ref)
	return ref, true
}

func (c chanField) IsZero(value Value) bool {
	return value.IsNil()
}

func (c chanField) Set(value Value, val string) {
	// TODO: channel capacity
	value.Set(reflect.MakeChan(value.Type(), 0))
}

type mapField struct{}

func (m mapField) Init(value Value, val string) (Value, bool) {
	ref := reflect.New(value.Type())
	ref.Elem().Set(reflect.MakeMap(value.Type()))
	return ref, true
}

func (m mapField) IsZero(value Value) bool {
	return value.Len() == 0
}

func (m mapField) Set(value Value, val string) {
	// TODO: map capacity
	value.Set(reflect.MakeMapWithSize(value.Type(), 0))
}

type funcField struct{}

func (f funcField) Init(value Value, val string) (Value, bool) {
	results := make([]Value, value.Type().NumOut())
	for i := range results {
		t := value.Type().Out(i)
		s := GetSetter(t.Kind())
		ref, _ := s.Init(reflect.New(t), "")
		results[i] = ref.Elem()
	}
	fn := func(_ []Value) []Value {
		return results
	}
	ref := reflect.MakeFunc(value.Type(), fn)
	return ref, true
}

func (f funcField) IsZero(value Value) bool {
	return value.IsNil()
}

func (f funcField) Set(value Value, val string) {

}

type interfaceField struct{}

func (i interfaceField) Init(value Value, val string) (Value, bool) {
	return value, false
}

func (i interfaceField) IsZero(value Value) bool {
	return value.IsNil()
}

func (i interfaceField) Set(value Value, val string) {
	// TODO implement me
	panic("implement me")
}

type arrayField struct{}

func (a arrayField) Init(value Value, val string) (Value, bool) {
	return value, false
}

func (a arrayField) IsZero(value Value) bool {
	return value.Len() == 0
}

func (a arrayField) Set(value Value, val string) {

}

type structField struct{}

func (s structField) Init(value Value, val string) (Value, bool) {
	return value, false
}

func (s structField) IsZero(value Value) bool {
	for i := 0; i < value.NumField(); i++ {
		if !value.Field(i).IsZero() {
			return false
		}
	}
	return true
}

func (s structField) Set(value Value, val string) {
	if val != "" && val != "{}" {
		if err := json.Unmarshal([]byte(val), value.Addr().Interface()); err == nil {
			return
		}
	}
	if err := LoadStruct(value.Addr().Interface()); err == nil {
		return
	}
}

type complexField struct{}

func (c complexField) Init(value Value, val string) (Value, bool) {
	return value, false
}

func (c complexField) IsZero(value Value) bool {
	complex := value.Complex()
	return math.Float64bits(real(complex)) == 0 && math.Float64bits(imag(complex)) == 0
}

func (c complexField) Set(value Value, val string) {
}

type uintptrField struct {
	field FieldSetter
}

func (u uintptrField) IsZero(value Value) bool {
	return true
}

func (u uintptrField) Init(value Value, val string) (Value, bool) {
	return value, false
}

func (u uintptrField) Set(value Value, val string) {
	u.field.Set(value, val)
}

var (
	fieldSetters [KindMax]FieldSetter
)

func GetSetter(kind Kind) FieldSetter {
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
	fieldSetters[reflect.Uintptr] = &uintptrField{field: fieldSetters[reflect.Uint]}
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
