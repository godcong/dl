package dl

import (
	"github.com/creasty/defaults"
)

// DefaultLoader is an interface that can be implemented by structs to customize the default
type DefaultLoader interface {
	Default() error
}

// Load initializes members in a struct referenced by a pointer.
// Maps and slices are initialized by `make` and other primitive types are set with default values.
// `ptr` should be a struct pointer
func Load[T any](ptr *T) error {
	if ok, err := LoadInterface(ptr); ok {
		return err
	}

	return LoadStruct(ptr)
}

// LoadInterface initializes members in a struct referenced by a pointer.
// Maps and slices are initialized by `make` and other primitive types are set with default values.
// `ptr` should be a struct pointer
func LoadInterface(ptr any) (bool, error) {
	if v, ok := ptr.(DefaultLoader); ok {
		return ok, v.Default()
	}
	return false, nil
}

// LoadStruct initializes members in a struct referenced by a pointer.
// Maps and slices are initialized by `make` and other primitive types are set with default values.
// `ptr` should be a struct pointer
func LoadStruct(ptr any) error {
	return defaults.Set(ptr)
}
