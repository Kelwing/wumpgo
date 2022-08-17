package router

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/kelwing/wumpgo/rest"
	"github.com/stretchr/testify/require"
)

type tapeItem struct {
	// Input
	FuncName string            `json:"func_name"`
	Params   []json.RawMessage `json:"params"`

	// Output
	Results      []json.RawMessage `json:"results"`
	GenericError string            `json:"generic_error,omitempty"`
	RESTError    *rest.ErrorREST   `json:"rest_error,omitempty"`
}

func mustMarshal(t TestingT, indent bool, item any) []byte {
	var b []byte
	var err error
	if indent {
		b, err = json.MarshalIndent(item, "", "  ")
	} else {
		b, err = json.Marshal(item)
	}
	require.NoError(t, err)
	return b
}

func (i *tapeItem) match(t TestingT, funcName string, isVard bool, inCount int, items ...any) {
	// Check the right function is called.
	if funcName != i.FuncName {
		t.Fatalf("wrong function called: expected %s, got %s", i.FuncName, funcName)
		return // Here for unit tests - in production this will never be hit.
	}

	// If the function is variadic, the params check is special.
	if isVard {
		// The error will be different.
		if inCount-1 > len(i.Params) {
			t.Fatalf("wrong number of inputs: got %d", (inCount-2)+len(i.Params))
			return // Here for unit tests - in production this will never be hit.
		}
	} else if inCount != len(i.Params) {
		t.Fatalf("wrong number of inputs: expected %d, got %d", len(i.Params), inCount)
		return // Here for unit tests - in production this will never be hit.
	}

	// Check all the params are equal.
	for x, p := range i.Params {
		end := len(items) - 1
		if x >= end && isVard {
			require.JSONEq(t, string(p), string(mustMarshal(t, false, reflect.ValueOf(items[x]).Field(x-end).Interface())))
		}
		require.JSONEq(t, string(p), string(mustMarshal(t, false, items[x])))
	}

	// Get the count of outputs.
	outCount := len(items) - inCount

	// Check if there is an error on the end.
	if outCount > 0 {
		ptr, _ := items[len(items)-1].(*error)
		if ptr != nil {
			if i.GenericError != "" {
				*ptr = errors.New(i.GenericError)
			} else if i.RESTError != nil {
				*ptr = i.RESTError
			}
			outCount--
		}
	}

	// Check the output count is equal to the number of outputs.
	if outCount != len(i.Results) {
		t.Fatalf("wrong number of outputs: expected %d, got %d", len(i.Results), outCount)
		return // Here for unit tests - in production this will never be hit.
	}

	// Handle the remainder of the params.
	for j, item := range i.Results {
		require.NoError(t, json.Unmarshal(item, items[inCount+j]))
	}
}

type tape []*tapeItem

func (t *tape) write(funcName string, isVard bool, params ...any) *tapeItem {
	undynamicLen := len(params)
	if isVard {
		undynamicLen--
	}
	p := make([]json.RawMessage, undynamicLen)
	for i, x := range params {
		if i == len(params)-1 && isVard {
			// Get the item from reflect.
			r := reflect.ValueOf(x)

			// Get each item from the slice and turn it into JSON.
			for j := 0; j < r.Len(); j++ {
				b, err := json.Marshal(r.Index(j).Interface())
				if err != nil {
					panic(err)
				}
				p = append(p, b)
			}

			// Break here.
			break
		}

		// Otherwise just handle it as standard.
		b, err := json.Marshal(x)
		if err != nil {
			panic(err)
		}
		p[i] = b
	}
	x := &tapeItem{
		FuncName: funcName,
		Params:   p,
	}
	*t = append(*t, x)
	return x
}

func (i *tapeItem) end(items ...any) {
	// Check if the last type is an error and if so split it from the items.
	var err error
	var ok bool
	if len(items) > 0 {
		err, ok = items[len(items)-1].(error)
		if ok || items[len(items)-1] == nil {
			items = items[:len(items)-1]
		}
	}

	// Marshal the rest of the results.
	p := make([]json.RawMessage, len(items))
	for i, x := range items {
		b, err := json.Marshal(x)
		if err != nil {
			panic(err)
		}
		p[i] = b
	}
	i.Results = p

	// Figure out how to process the error.
	if err != nil {
		if e, ok := err.(*rest.ErrorREST); ok {
			i.RESTError = e
		} else {
			i.GenericError = err.Error()
		}
	}
}
