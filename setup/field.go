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

type FieldSetter interface {
	IsZero(value reflect.Value) bool
	Init(value reflect.Value) reflect.Value
	Set(value reflect.Value, val string)
}

type Field struct {
	Setter   FieldSetter
	Value    reflect.Value
	TagValue string
}

func (f Field) IsZero() bool {
	return f.Setter.IsZero(f.Value)
}

func (f Field) Init() reflect.Value {
	return f.Setter.Init(f.Value)
}
func (f Field) Fill() {
	if !f.IsZero() {
		return
	}
	f.Setter.Set(f.Value, f.TagValue)
}

func (f Field) CanSet() bool {
	return f.Value.CanSet()
}

func (f Field) Kind() reflect.Kind {
	return f.Value.Kind()
}

func (f Field) Set(convert reflect.Value) {
	f.Value.Set(convert)
}

type floatField struct {
	size int
}

func (f floatField) Init(value reflect.Value) reflect.Value {
	// TODO implement me
	panic("implement me")
}

func (f floatField) IsZero(value reflect.Value) bool {
	return math.Float64bits(value.Float()) == 0
}

func (f floatField) Set(value reflect.Value, val string) {
	if val, err := strconv.ParseFloat(val, f.size); err == nil {
		value.SetFloat(val)
	}
}

type boolField struct {
}

func (b boolField) Init(value reflect.Value) reflect.Value {
	// TODO implement me
	panic("implement me")
}

func (b boolField) IsZero(value reflect.Value) bool {
	return !value.Bool()
}

func (b boolField) Set(value reflect.Value, val string) {
	if val, err := strconv.ParseBool(val); err == nil {
		value.SetBool(val)
	}
}

type intField struct {
	size int
}

func (i intField) Init(value reflect.Value) reflect.Value {
	// TODO implement me
	panic("implement me")
}

func (i intField) IsZero(value reflect.Value) bool {
	return value.Int() == 0
}

func (i intField) Set(value reflect.Value, val string) {
	if i.size == 64 {
		if val, err := time.ParseDuration(val); err == nil {
			value.Set(reflect.ValueOf(val).Convert(value.Type()))
			return
		}
	}
	if val, err := strconv.ParseInt(val, 0, i.size); err == nil {
		value.SetInt(val)
	}
}

type field struct {
	size int
}

func (u field) Init(value reflect.Value) reflect.Value {
	// TODO implement me
	panic("implement me")
}

func (u field) IsZero(value reflect.Value) bool {
	return value.Uint() == 0
}

func (u field) Set(value reflect.Value, val string) {
	if val, err := strconv.ParseUint(val, 0, u.size); err == nil {
		value.SetUint(val)
	}
}

type stringField struct {
}

func (s stringField) Init(value reflect.Value) reflect.Value {
	// TODO implement me
	panic("implement me")
}

func (s stringField) IsZero(value reflect.Value) bool {
	return value.String() == ""
}

func (s stringField) Set(value reflect.Value, val string) {
	value.SetString(val)
}

type sliceField struct {
}

func (s sliceField) Init(value reflect.Value) reflect.Value {
	ref := reflect.New(value.Type())
	ref.Elem().Set(reflect.MakeSlice(value.Type(), 0, 0))
	return ref
}

func (s sliceField) IsZero(value reflect.Value) bool {
	return value.Len() == 0
}

func (s sliceField) Set(value reflect.Value, val string) {
	if val != "" && val != "[]" {
		if err := json.Unmarshal([]byte(val), value.Interface()); err != nil {
		}
	}
	value.Set(value.Elem().Convert(value.Type()))
}

type emptyField struct {
}

func (e emptyField) Init(value reflect.Value) reflect.Value {
	// TODO implement me
	panic("implement me")
}

func (e emptyField) IsZero(value reflect.Value) bool {
	return true
}

func (e emptyField) Set(value reflect.Value, val string) {
}

type pointerField struct {
}

func (p pointerField) Init(value reflect.Value) reflect.Value {
	ref := reflect.New(value.Type().Elem())
	value.Set(ref)
	return ref
}

func (p pointerField) IsZero(value reflect.Value) bool {
	return value.IsNil()
}

func (p pointerField) Set(value reflect.Value, val string) {
	value.Set(reflect.New(value.Type().Elem()))
}

type chanField struct {
}

func (c chanField) Init(value reflect.Value) reflect.Value {
	// ref := reflect.MakeChan(value.Type(), 0)
	// value.Set(ref)
	// return ref
	return value
}

func (c chanField) IsZero(value reflect.Value) bool {
	return value.IsNil()
}

func (c chanField) Set(value reflect.Value, val string) {
	// TODO: channel capacity
	value.Set(reflect.MakeChan(value.Type(), 0))
}

type mapField struct {
}

func (m mapField) Init(value reflect.Value) reflect.Value {
	ref := reflect.New(value.Type())
	ref.Elem().Set(reflect.MakeMap(value.Type()))
	return ref
}

func (m mapField) IsZero(value reflect.Value) bool {
	return value.Len() == 0
}

func (m mapField) Set(value reflect.Value, val string) {
	// TODO: map capacity
	value.Set(reflect.MakeMapWithSize(value.Type(), 0))
}

type funcField struct {
}

func (f funcField) Init(value reflect.Value) reflect.Value {
	// TODO implement me
	panic("implement me")
}

func (f funcField) IsZero(value reflect.Value) bool {
	return value.IsNil()
}

func (f funcField) Set(value reflect.Value, val string) {
	fn := func(args []reflect.Value) (results []reflect.Value) {
		return
	}
	value.Set(reflect.MakeFunc(value.Type(), fn))
}

type interfaceField struct {
}

func (i interfaceField) Init(value reflect.Value) reflect.Value {
	// TODO implement me
	panic("implement me")
}

func (i interfaceField) IsZero(value reflect.Value) bool {
	return value.IsNil()
}

func (i interfaceField) Set(value reflect.Value, val string) {
	// TODO implement me
	panic("implement me")
}

type arrayField struct {
}

func (a arrayField) Init(value reflect.Value) reflect.Value {
	// TODO implement me
	panic("implement me")
}

func (a arrayField) IsZero(value reflect.Value) bool {
	return value.Len() == 0
}

func (a arrayField) Set(value reflect.Value, val string) {
}

type structField struct {
}

func (s structField) Init(value reflect.Value) reflect.Value {
	// TODO implement me
	panic("implement me")
}

func (s structField) IsZero(value reflect.Value) bool {
	for i := 0; i < value.NumField(); i++ {
		if !value.Field(i).IsZero() {
			return false
		}
	}
	return true
}

func (s structField) Set(value reflect.Value, val string) {
	// TODO implement me
	panic("implement me")
}

type complexField struct {
}

func (c complexField) Init(value reflect.Value) reflect.Value {
	// TODO implement me
	panic("implement me")
}

func (c complexField) IsZero(value reflect.Value) bool {
	complex := value.Complex()
	return math.Float64bits(real(complex)) == 0 && math.Float64bits(imag(complex)) == 0
}

func (c complexField) Set(value reflect.Value, val string) {
	// TODO implement me
	panic("implement me")
}

var (
	fieldSetters [KindMax]FieldSetter
)

func GetSetter(kind reflect.Kind) FieldSetter {
	return fieldSetters[kind]
}

type uintptrField struct {
	field FieldSetter
}

func (u uintptrField) IsZero(value reflect.Value) bool {
	return true
}

func (u uintptrField) Init(value reflect.Value) reflect.Value {
	return value
}

func (u uintptrField) Set(value reflect.Value, val string) {
	u.field.Set(value, val)
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
	fieldSetters[reflect.Uint] = &field{size: strconv.IntSize}
	fieldSetters[reflect.Uint8] = &field{size: 8}
	fieldSetters[reflect.Uint16] = &field{size: 16}
	fieldSetters[reflect.Uint32] = &field{size: 32}
	fieldSetters[reflect.Uint64] = &field{size: 64}
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
