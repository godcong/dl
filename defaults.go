// Copyright (c) 2024 GodCong. All rights reserved.

// Package dl for Default Loader
package dl

import (
	"encoding"
	"encoding/json"
	"reflect"
	"strconv"
	"time"
)

const (
	fieldName = "default"
)

func setDefaults(ptr interface{}) error {
	kind := reflect.TypeOf(ptr).Kind()
	if kind != reflect.Ptr {
		return InvalidTypeError(kind.String())
	}

	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()

	if t.Kind() != reflect.Struct {
		return InvalidTypeError(t.Kind().String())
	}

	for i := 0; i < t.NumField(); i++ {
		if defaultVal := t.Field(i).Tag.Get(fieldName); defaultVal != "-" {
			if err := setField(v.Field(i), defaultVal); err != nil {
				return err
			}
		}
	}

	return nil
}

func setIntField(field reflect.Value, defaultVal string, size int) {
	if size == 64 {
		if val, err := time.ParseDuration(defaultVal); err == nil {
			field.Set(reflect.ValueOf(val).Convert(field.Type()))
			return
		}
	}
	if val, err := strconv.ParseInt(defaultVal, 0, size); err == nil {
		field.SetInt(val)
	}
}

func setUintField(field reflect.Value, defaultVal string, size int) {
	if val, err := strconv.ParseUint(defaultVal, 0, size); err == nil {
		field.SetUint(val)
	}
}

func setFloatField(field reflect.Value, defaultVal string, size int) {
	if val, err := strconv.ParseFloat(defaultVal, size); err == nil {
		field.SetFloat(val)
	}
}

func setObjectField(field reflect.Value, defaultVal string) {
	if val, err := strconv.ParseBool(defaultVal); err == nil {
		field.Set(reflect.ValueOf(val))
	}
}

func setField(field reflect.Value, defaultVal string) error {
	if !field.CanSet() {
		return nil
	}

	if !shouldInitializeField(field, defaultVal) {
		return nil
	}

	isInitial := isInitialValue(field)
	if isInitial {
		if unmarshalByInterface(field, defaultVal) {
			return nil
		}

		switch field.Kind() {
		case reflect.Bool:
			if val, err := strconv.ParseBool(defaultVal); err == nil {
				field.SetBool(val)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			setIntField(field, defaultVal, field.Type().Bits())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			setUintField(field, defaultVal, field.Type().Bits())
		case reflect.Float32, reflect.Float64:
			setFloatField(field, defaultVal, field.Type().Bits())
		case reflect.String:
			field.SetString(defaultVal)
		case reflect.Slice:
			ref := reflect.New(field.Type())
			ref.Elem().Set(reflect.MakeSlice(field.Type(), 0, 0))
			if defaultVal != "" && defaultVal != "[]" {
				if err := json.Unmarshal([]byte(defaultVal), ref.Interface()); err != nil {
					return err
				}
			}
			field.Set(ref.Elem().Convert(field.Type()))
		case reflect.Map:
			ref := reflect.New(field.Type())
			ref.Elem().Set(reflect.MakeMap(field.Type()))
			if defaultVal != "" && defaultVal != "{}" {
				if err := json.Unmarshal([]byte(defaultVal), ref.Interface()); err != nil {
					return err
				}
			}
			field.Set(ref.Elem().Convert(field.Type()))
		case reflect.Struct:
			if defaultVal != "" && defaultVal != "{}" {
				if err := json.Unmarshal([]byte(defaultVal), field.Addr().Interface()); err != nil {
					return err
				}
			}
		case reflect.Ptr:
			field.Set(reflect.New(field.Type().Elem()))
		default:
			// nothing to do
		}
	}

	switch field.Kind() {
	case reflect.Ptr:
		if isInitial || field.Elem().Kind() == reflect.Struct {
			err := setField(field.Elem(), defaultVal)
			if err != nil {
				return err
			}
		}
	case reflect.Struct:
		if err := LoadStruct(field.Addr().Interface()); err != nil {
			return err
		}
	case reflect.Slice:
		for j := 0; j < field.Len(); j++ {
			if err := setField(field.Index(j), defaultVal); err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, e := range field.MapKeys() {
			var v = field.MapIndex(e)

			switch v.Kind() {
			case reflect.Ptr:
				switch v.Elem().Kind() {
				case reflect.Struct, reflect.Slice, reflect.Map:
					if err := setField(v.Elem(), ""); err != nil {
						return err
					}
				default:
					// nothing to do
				}
			case reflect.Struct, reflect.Slice, reflect.Map:
				ref := reflect.New(v.Type())
				ref.Elem().Set(v)
				if err := setField(ref.Elem(), ""); err != nil {
					return err
				}
				field.SetMapIndex(e, ref.Elem().Convert(v.Type()))
			default:
				// nothing to do
			}
		}
	default:
		// nothing to do
	}

	return nil
}

func unmarshalByInterface(field reflect.Value, defaultVal string) bool {
	asText, ok := field.Addr().Interface().(encoding.TextUnmarshaler)
	if ok && defaultVal != "" {
		// if field implements encode.TextUnmarshaler, try to use it before decode by kind
		if err := asText.UnmarshalText([]byte(defaultVal)); err == nil {
			return true
		}
	}
	asJSON, ok := field.Addr().Interface().(json.Unmarshaler)
	if ok && defaultVal != "" && defaultVal != "{}" && defaultVal != "[]" {
		// if field implements json.Unmarshaler, try to use it before decode by kind
		if err := asJSON.UnmarshalJSON([]byte(defaultVal)); err == nil {
			return true
		}
	}
	return false
}

func isInitialValue(field reflect.Value) bool {
	return reflect.DeepEqual(reflect.Zero(field.Type()).Interface(), field.Interface())
}

func shouldInitializeField(field reflect.Value, tag string) bool {
	switch field.Kind() {
	case reflect.Struct:
		return true
	case reflect.Ptr:
		if !field.IsNil() && field.Elem().Kind() == reflect.Struct {
			return true
		}
	case reflect.Slice:
		return field.Len() > 0 || tag != ""
	case reflect.Map:
		return field.Len() > 0 || tag != ""
	default:
		// nothing to do
	}

	return tag != ""
}

// CanUpdate returns true when the given value is an initial value of its type
func CanUpdate(v interface{}) bool {
	return isInitialValue(reflect.ValueOf(v))
}
