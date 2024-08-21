package setup

import (
	"errors"
	"reflect"
	"testing"
)

func TestInvalidTypeError(t *testing.T) {
	err := InvalidTypeError(reflect.String)
	if err == nil {
		t.Error("InvalidTypeError should not return nil")
	}

	var e *invalidTypeErr
	if !errors.As(err, &e) {
		t.Error("InvalidTypeError should return an *invalidTypeErr")
	}

	if e.kind.String() != "string" {
		t.Errorf("InvalidTypeError should set typeString to the correct value, got %s", e.kind.String())
	}
}

func TestInvalidErr_Error(t *testing.T) {
	err := &invalidTypeErr{kind: reflect.String}
	if err.Error() != "invalid type string" {
		t.Errorf("invalidTypeErr.Error should return the correct string, got %s", err.Error())
	}
}

func TestInvalidErr_Is(t *testing.T) {
	err1 := &invalidTypeErr{kind: reflect.String}
	err2 := &invalidTypeErr{kind: reflect.Int}

	if !err1.Is(err2) {
		t.Error("err1.Is(err2) should return true")
	}

	err3 := &invalidTypeErr{kind: reflect.String}
	if !err1.Is(err3) {
		t.Error("err1.Is(err3) should return true")
	}

	if err1.Is(nil) {
		t.Error("err1.Is(nil) should return false")
	}
}

func TestError_Is(t *testing.T) {
	err1 := &invalidTypeErr{kind: reflect.String}
	err2 := &invalidTypeErr{kind: reflect.Int}

	if !errors.Is(err1, err2) {
		t.Error("errors.Is(err1, err2) should return true")
	}
}
