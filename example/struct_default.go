// Code generated by github.com/godcong/dl. DO NOT EDIT.
// Version: devel
// Commit: d9114ea87d890c13bc7bf9e3d6d1f5510fe66e73
// Build Date: 2024-07-25T08:24:23
// Built By: unknown

package example

import (
	"github.com/godcong/dl"
)

// Default loads default values for StructStruct
func (obj *StructStruct) Default() error {
	obj.Key = "key"
	obj.Value = "value"
	return nil
}

// Default loads default values for StructStd
func (obj *StructStd) Default() error {
	obj.FieldString = "test"
	obj.FieldInt = 1
	obj.FieldFloat64 = 1.1
	obj.FieldBytes = []byte("test")
	obj.FieldPBytes = dl.Pointer([]byte("test"))
	obj.FieldIntSlice = []int{1, 2, 3}
	obj.FieldPIntSlice = []*int{dl.Pointer(1), dl.Pointer(2), dl.Pointer(3)}
	obj.FieldPIntPSlice = dl.Pointer([]*int{dl.Pointer(1), dl.Pointer(2), dl.Pointer(3)})
	obj.FieldStringSlice = []string{"test", "test2"}
	obj.FieldStringPSlice = []*string{dl.Pointer("test"), dl.Pointer("test2")}
	obj.FieldPStringSlice = dl.Pointer([]string{"test", "test2"})
	obj.FieldPStringPSlice = dl.Pointer([]*string{dl.Pointer("test"), dl.Pointer("test2")})
	obj.FieldBool = true
	obj.FieldMapStringString = map[string]string{"key1": "value1", "key2": "value2"}
	obj.FieldMapPStringPString = map[*string]*string{dl.Pointer("key1"): dl.Pointer("value1"), dl.Pointer("key2"): dl.Pointer("value2")}
	obj.FieldMapPBytesPBytes = map[*[]byte]*[]byte{dl.Pointer([]byte("key1")): dl.Pointer([]byte("value1")), dl.Pointer([]byte("key2")): dl.Pointer([]byte("value2"))}
	obj.FieldMapIntString = map[string]string{"key1": "value1", "key2": "value2"}
	obj.FieldMapIntInt = map[int]int{1: 11, 2: 22}
	obj.FieldMapStringInt = map[string]int{"value1": 11, "value2": 22}
	return nil
}

// Default loads default values for StructInner
func (obj *StructInner) Default() error {
	if err := dl.Load(&obj.FieldInnerStruct); err != nil {
		return err
	}
	return nil
}
