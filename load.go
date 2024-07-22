// Copyright (c) 2024 GodCong. All rights reserved.

// Package dl for Default Loader
package dl

import (
	"github.com/creasty/defaults"
)

// DefaultLoader is an interface that can be implemented by structs to customize the default
type DefaultLoader interface {
	Default() error
}

// DefaultLoaderFunc is a function type that defines a function to load default values into a struct referenced by a pointer.
type DefaultLoaderFunc[T any] func(*T) error

// DefaultOptionLoader is an interface that specifies a method to load default values into a struct with a parameter.
type DefaultOptionLoader[P any] interface {
	Default(P) error
}

// DefaultOptionLoaderFunc is a function type that defines a function to load default values into a struct with a parameter.
type DefaultOptionLoaderFunc[T any, P any] func(*T, P) error

// Load initializes members in a struct referenced by a pointer.
// Maps and slices are initialized by `make` and other primitive types are set with default values.
// `ptr` should be a struct pointer
func Load[T any, P any](ptr *T, args ...P) error {
	var arg P
	if len(args) > 0 {
		arg = args[0]
	}
	if ok, err := LoadInterface(ptr, arg); ok {
		return err
	}
	return LoadStruct(ptr)
}

// LoadInterface initializes members in a struct referenced by a pointer.
// Maps and slices are initialized by `make` and other primitive types are set with default values.
// `ptr` should be a struct pointer
func LoadInterface[P any](ptr any, arg P) (bool, error) {
	if v, ok := ptr.(DefaultLoader); ok {
		return ok, v.Default()
	}
	switch p := ptr.(type) {
	case DefaultLoader:
		return true, p.Default()
	case DefaultOptionLoader[P]:
		return true, p.Default(arg)
	}

	return false, nil
}

// LoadStruct initializes members in a struct referenced by a pointer.
// Maps and slices are initialized by `make` and other primitive types are set with default values.
// `ptr` should be a struct pointer
func LoadStruct(ptr any) error {
	return defaults.Set(ptr)
}

// Pointer creates a pointer to a value.
func Pointer[T any](v T) *T {
	return &v
}
