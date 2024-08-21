// Copyright (c) 2024 GodCong. All rights reserved.

// Package setup for Default Loader
package setup

import (
	"errors"
	"fmt"
)

var ErrInvalidType = errors.New("empty")

type invalidTypeErr struct {
	kind Kind
}

func (i *invalidTypeErr) Error() string {
	return fmt.Sprintf("invalid type %s", i.kind.String())
}

func (i *invalidTypeErr) Is(err error) bool {
	if err == nil {
		return false
	}
	var invalidTypeErr *invalidTypeErr
	if errors.As(err, &invalidTypeErr) {
		return i.kind == invalidTypeErr.kind
	}
	return false
}

func InvalidTypeError(kind Kind) error {
	return &invalidTypeErr{kind: kind}
}
