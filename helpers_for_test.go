package router

import (
	"reflect"
	"testing"
)

func callBuilderFunction(t *testing.T, builder interface{}, funcName string, args ...interface{}) (err error) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			err = ungenericError(r)
		}
	}()
	r := reflect.ValueOf(builder).MethodByName(funcName)
	if r.IsZero() {
		t.Fatal("function does not exist")
	}
	reflectArgs := make([]reflect.Value, len(args))
	for i, v := range args {
		reflectArgs[i] = reflect.ValueOf(v)
	}
	if r.Kind() != reflect.Func {
		t.Fatal("not a function")
	}
	res := r.Call(reflectArgs)
	if len(res) != 1 {
		t.Fatal("arg count not correct for builder:", res)
	}
	if res[0].Interface() != builder {
		t.Fatal("the argument returned was not the builder")
	}
	return
}
