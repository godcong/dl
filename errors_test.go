package dl

import (
	"errors"
	"testing"
)

func TestInvalidTypeError(t *testing.T) {
	err := InvalidTypeError("string")
	if err == nil {
		t.Error("InvalidTypeError should not return nil")
	}

	var e *invalidTypeErr
	if !errors.As(err, &e) {
		t.Error("InvalidTypeError should return an *invalidTypeErr")
	}

	if e.typeString != "string" {
		t.Errorf("InvalidTypeError should set typeString to the correct value, got %s", e.typeString)
	}
}

func TestInvalidErr_Error(t *testing.T) {
	err := &invalidTypeErr{typeString: "string"}
	if err.Error() != "invalid type string" {
		t.Errorf("invalidTypeErr.Error should return the correct string, got %s", err.Error())
	}
}

func TestInvalidErr_Is(t *testing.T) {
	err1 := &invalidTypeErr{typeString: "string"}
	err2 := &invalidTypeErr{typeString: "int"}

	if !err1.Is(err2) {
		t.Error("err1.Is(err2) should return true")
	}

	err3 := &invalidTypeErr{typeString: "string"}
	if !err1.Is(err3) {
		t.Error("err1.Is(err3) should return true")
	}

	if err1.Is(nil) {
		t.Error("err1.Is(nil) should return false")
	}
}

func TestError_Is(t *testing.T) {
	err1 := &invalidTypeErr{typeString: "string"}
	err2 := &invalidTypeErr{typeString: "int"}

	if !errors.Is(err1, err2) {
		t.Error("errors.Is(err1, err2) should return true")
	}
}
