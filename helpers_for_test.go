package router

import (
	"reflect"
	"testing"
)

func callBuilderFunction[T any](t *testing.T, builder T, expectsIface bool, funcName string, args ...any) (err error) {
	// Signify that this is a test helper.
	t.Helper()

	// Recover if this fails to gracefully fail the test.
	defer func() {
		if r := recover(); r != nil {
			err = ungenericError(r)
		}
	}()

	// Check this is an interface and that said interface has the function.
	if expectsIface {
		ifaceType := reflect.TypeOf(&builder).Elem()
		if ifaceType.Kind() != reflect.Interface {
			t.Fatal("not a interface")
		}
		_, has := ifaceType.MethodByName(funcName)
		if !has {
			t.Fatal("type does not have function")
		}
	}

	// Get the function and try to call it.
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
	if res[0].Interface() != (any)(builder) {
		t.Fatal("the argument returned was not the builder")
	}
	return
}
