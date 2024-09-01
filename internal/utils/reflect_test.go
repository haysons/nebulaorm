package utils

import (
	"reflect"
	"testing"
)

func TestSliceSetElem(t *testing.T) {
	sliceInt := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	sliceIntNew := make([]int, 0)
	err := SliceSetElem(reflect.ValueOf(&sliceIntNew).Elem(), len(sliceInt), func(i int, elem reflect.Value) (bool, error) {
		if i >= len(sliceInt) {
			return false, nil
		}
		elem.SetInt(int64(sliceInt[i]))
		return true, nil
	})
	if err != nil {
		t.Errorf("slice int set elem err: %v", err)
		return
	}
	if !reflect.DeepEqual(sliceInt, sliceIntNew) {
		t.Errorf("slice int set elem err: %v != %v", sliceInt, sliceIntNew)
		return
	}

	arrayInt := [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	arrayIntNew := [9]int{}
	err = SliceSetElem(reflect.ValueOf(&arrayIntNew).Elem(), len(arrayInt), func(i int, elem reflect.Value) (bool, error) {
		if i >= len(arrayInt) {
			return false, nil
		}
		elem.SetInt(int64(sliceInt[i]))
		return true, nil
	})
	if err != nil {
		t.Errorf("array int set elem err: %v", err)
		return
	}
	if !reflect.DeepEqual(arrayInt, arrayIntNew) {
		t.Errorf("array int set elem err: %v != %v", sliceInt, sliceIntNew)
		return
	}

	arrayIntNew = [9]int{}
	err = SliceSetElem(reflect.ValueOf(&arrayIntNew).Elem(), 50, func(i int, elem reflect.Value) (bool, error) {
		if i >= 50 {
			return false, nil
		}
		elem.SetInt(int64(i + 1))
		return true, nil
	})
	if err != nil {
		t.Errorf("array int set elem err: %v", err)
		return
	}
	if !reflect.DeepEqual(arrayInt, arrayIntNew) {
		t.Errorf("array int set elem err: %v != %v", sliceInt, sliceIntNew)
		return
	}

	a, b, c := 1, 2, 3
	slicePtr := []*int{&a, &b, &c}
	slicePtrNew := make([]*int, 0, 1)
	err = SliceSetElem(reflect.ValueOf(&slicePtrNew).Elem(), len(slicePtr), func(i int, elem reflect.Value) (bool, error) {
		if i >= len(slicePtr) {
			return false, nil
		}
		elem.SetInt(int64(*slicePtr[i]))
		return true, nil
	})
	if err != nil {
		t.Errorf("slice ptr set elem err: %v", err)
		return
	}
	if !reflect.DeepEqual(slicePtr, slicePtrNew) {
		t.Errorf("slice ptr set elem err: %v != %v", slicePtr, slicePtrNew)
		return
	}
}

func TestPtrValue(t *testing.T) {
	var a *int
	aValue := reflect.ValueOf(&a)
	aValue = PtrValue(aValue)
	aValue.SetInt(1)
	if *a != 1 {
		t.Errorf("a int set value err: %v != %v", *a, 1)
		return
	}

	var b int
	bValue := reflect.ValueOf(&b)
	bValue = PtrValue(bValue)
	bValue.SetInt(1)
	if b != 1 {
		t.Errorf("b int set value err: %v != %v", b, 1)
		return
	}

	var c ***string
	cValue := reflect.ValueOf(&c)
	cValue = PtrValue(cValue)
	cValue.SetString("hello")
	if ***c != "hello" {
		t.Errorf("c string set value err: %v != %v", ***c, "hello")
		return
	}
}
