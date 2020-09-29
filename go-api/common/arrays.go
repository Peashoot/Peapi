package common

import (
	"fmt"
	"reflect"
)

// IndexOf - find index of first item in array(or slice, string)
func IndexOf(array interface{}, compare interface{}) int {
	arrayType := reflect.ValueOf(array)
	if arrayType.Kind() != reflect.Array && arrayType.Kind() != reflect.Slice && arrayType.Kind() != reflect.String {
		panic("The type of array isn't correct")
	}
	for i := 0; i < arrayType.Len(); i++ {
		if arrayType.Index(i).Interface() == compare {
			return i
		}
	}
	return -1
}

// LastIndexOf - find index of last item in array(or slice, string)
func LastIndexOf(array interface{}, compare interface{}) int {
	arrayType := reflect.ValueOf(array)
	if arrayType.Kind() != reflect.Array && arrayType.Kind() != reflect.Slice && arrayType.Kind() != reflect.String {
		panic("The type of array isn't correct")
	}
	compareType := reflect.ValueOf(compare)
	fmt.Println(compareType.Kind())
	for i := arrayType.Len() - 1; i >= 0; i-- {
		fmt.Println(arrayType.Index(i).Kind())
		fmt.Println(arrayType.Index(i).Interface())
		if arrayType.Index(i).Interface() == compare {
			return i
		}
	}
	return -1
}
