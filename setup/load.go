// Copyright (c) 2024 GodCong. All rights reserved.

// Package setup for Default Loader
package setup

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
func Load[T any](ptr *T) error {
	if ok, err := LoadInterface(ptr, any(nil)); ok {
		return err
	}
	return LoadStruct(ptr)
}

// MustLoad initializes members in a struct referenced by a pointer.
// Maps and slices are initialized by `make` and other primitive types are set with default values.
// `ptr` should be a struct pointer
func MustLoad[T any](ptr *T) {
	if err := Load(ptr); err != nil {
		panic(err)
	}
}

// LoadWithOption initializes members in a struct referenced by a pointer.
// Maps and slices are initialized by `make` and other primitive types are set with default values.
// `ptr` should be a struct pointer
func LoadWithOption[T any, P any](ptr *T, arg P) error {
	if ok, err := LoadInterface(ptr, arg); ok {
		return err
	}
	return LoadStruct(ptr)
}

// LoadInterface initializes members in a struct referenced by a pointer.
// Maps and slices are initialized by `make` and other primitive types are set with default values.
// `ptr` should be a struct pointer
func LoadInterface[P any](ptr any, arg P) (bool, error) {
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
	return setDefaults(ptr)
}

// Pointer creates a pointer to a value.
func Pointer[T any](v T) *T {
	return &v
}

// Object creates a value from a pointer.
func Object[T any](v *T) T {
	return *v
}
