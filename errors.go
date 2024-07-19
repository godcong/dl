// Copyright (c) 2024 GodCong. All rights reserved.

// Package dl for Default Loader
package dl

import (
	"errors"
	"fmt"
)

var ErrInvalidType = errors.New("empty")

type invalidTypeErr struct {
	typeString string
}

func (i *invalidTypeErr) Error() string {
	return fmt.Sprintf("invalid type %s", i.typeString)
}

func (i *invalidTypeErr) Is(err error) bool {
	if err == nil {
		return false
	}
	var invalidTypeErr *invalidTypeErr
	return errors.As(err, &invalidTypeErr)
}

func InvalidTypeError(typeString string) error {
	return &invalidTypeErr{typeString: typeString}
}
