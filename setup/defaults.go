// Copyright (c) 2024 GodCong. All rights reserved.

// Package setup for Default Loader
package setup

import (
	"encoding/json"
	"fmt"
	"reflect"
)

const (
	tagName = "default"
)

func setDefaults(ptr interface{}) error {
	kind := reflect.TypeOf(ptr).Kind()
	if kind != reflect.Ptr {
		return InvalidTypeError(kind)
	}

	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()

	if t.Kind() != reflect.Struct {
		return InvalidTypeError(t.Kind())
	}

	for i := 0; i < t.NumField(); i++ {
		if defaultVal := t.Field(i).Tag.Get(tagName); defaultVal != "-" {
			field := v.Field(i)
			fmt.Println("kind", field.Kind(), field.String())
			if err := setField(&Field{
				Setter:   GetSetter(field.Kind()),
				Value:    field,
				TagValue: defaultVal,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func setField(field *Field) error {
	if !field.CanSet() {
		return nil
	}

	if !shouldInitializeField(field.Value, field.TagValue) {
		return nil
	}

	needInit := field.IsZero()
	if needInit {
		if unmarshal(field.Value, field.TagValue) {
			return nil
		}

		switch field.Kind() {
		case reflect.Bool:
			field.Fill()
			// if val, err := strconv.ParseBool(defaultVal); err == nil {
			// 	uintField.SetBool(val)
			// }
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			field.Fill()
			// setIntField(uintField, defaultVal, uintField.Type().Bits())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			field.Fill()
			// setUintField(uintField, defaultVal, uintField.Type().Bits())
		case reflect.Float32, reflect.Float64:
			field.Fill()
			// setFloatField(uintField, defaultVal, uintField.Type().Bits())
		case reflect.String:
			field.Fill()
			// uintField.SetString(defaultVal)
		case reflect.Slice:
			ref, truth := field.Init()
			// ref := reflect.New(uintField.Type())
			// ref.Elem().Set(reflect.MakeSlice(uintField.Type(), 0, 0))
			if truth {
				if field.TagValue != "" && field.TagValue != "[]" {
					if err := json.Unmarshal([]byte(field.TagValue), ref.Interface()); err != nil {
						return err
					}
				}
				field.SetRef(ref.Elem().Convert(field.Value.Type()))
			}
		case reflect.Map:
			ref, truth := field.Init()
			if truth {
				if field.TagValue != "" && field.TagValue != "{}" {
					if err := json.Unmarshal([]byte(field.TagValue), ref.Interface()); err != nil {
						return err
					}
				}
				field.SetRef(ref.Elem().Convert(field.Value.Type()))
			}
		case reflect.Struct:
			ref, truth := field.Init()
			if truth {
				if field.TagValue != "" && field.TagValue != "{}" {
					if err := json.Unmarshal([]byte(field.TagValue), ref.Addr().Interface()); err != nil {
						return err
					}
				}
			}
		case reflect.Pointer:
			ref, truth := field.Init()
			if truth {
				field.SetRef(ref)
			}
		case reflect.Func:
			ref, truth := field.Init()
			if truth {
				field.SetRef(ref)
			}
		case reflect.Chan:
			ref, truth := field.Init()
			if truth {
				field.SetRef(ref)
			}
		default:
			// nothing to do
		}
	}

	switch field.Kind() {
	case reflect.Pointer:
		if needInit || field.Value.Elem().Kind() == reflect.Struct {
			err := setField(&Field{
				Setter:   GetSetter(field.Value.Elem().Kind()),
				Value:    field.Value.Elem(),
				TagValue: field.TagValue,
			})
			if err != nil {
				return err
			}
		}
	case reflect.Struct:
		if err := LoadStruct(field.Value.Addr().Interface()); err != nil {
			return err
		}
	case reflect.Slice:
		for j := 0; j < field.Value.Len(); j++ {
			v := field.Value.Index(j)
			if err := setField(&Field{
				Setter:   GetSetter(v.Kind()),
				Value:    v,
				TagValue: field.TagValue,
			}); err != nil {
				return err
			}
		}
	case reflect.Map:
		iter := field.Value.MapRange()
		for iter.Next() {
			v := iter.Value()
			switch v.Kind() {
			case reflect.Ptr:
				elem := v.Elem()
				switch elem.Kind() {
				case reflect.Struct, reflect.Slice, reflect.Map:
					if err := setField(&Field{
						Setter: GetSetter(elem.Kind()),
						Value:  elem,
					}); err != nil {
						return err
					}
				default:
					// nothing to do
				}
			case reflect.Struct, reflect.Slice, reflect.Map:
				ref := reflect.New(v.Type())
				ref.Elem().Set(v)
				if err := setField(&Field{
					Setter: GetSetter(ref.Kind()),
					Value:  ref,
				}); err != nil {
					return err
				}
				field.Value.SetMapIndex(iter.Key(), ref.Elem())
			default:
				// nothing to do
			}
		}
	default:
		// nothing to do
	}

	return nil
}

func unmarshal(field reflect.Value, defaultVal string) bool {
	if defaultVal == "" {
		return false
	}
	if asText, ok := field.Addr().Interface().(interface {
		UnmarshalText(text []byte) error
	}); ok {
		// if uintField implements encode.TextUnmarshaler, try to use it before decode by kind
		if err := asText.UnmarshalText([]byte(defaultVal)); err == nil {
			return true
		}
	}
	if asJSON, ok := field.Addr().Interface().(interface {
		UnmarshalJSON([]byte) error
	}); ok && defaultVal != "{}" && defaultVal != "[]" {
		// if uintField implements json.Unmarshaler, try to use it before decode by kind
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
