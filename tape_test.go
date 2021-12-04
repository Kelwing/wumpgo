package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/Postcord/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testingTTapeItemMatch struct {
	*testing.T

	format string
	args   []interface{}
}

func (t *testingTTapeItemMatch) Fatalf(format string, args ...interface{}) {
	t.format = format
	t.args = args
}

func Test_tapeItem_match(t *testing.T) {
	t.Run("wrong func name", func(t *testing.T) {
		x := tapeItem{FuncName: "foo"}
		m := &testingTTapeItemMatch{T: t}
		x.match(m, "bar", false, 0)
		assert.Equal(t, "wrong function called: expected foo, got bar", fmt.Sprintf(m.format, m.args...))
	})

	t.Run("variadic wrong count", func(t *testing.T) {
		x := tapeItem{FuncName: "foo", Params: []json.RawMessage{[]byte(`"a"`), []byte(`"b"`)}}
		m := &testingTTapeItemMatch{T: t}
		x.match(m, "foo", true, 4, "a", []string{"b", "c", "d"})
		assert.Equal(t, "wrong number of inputs: got 4", fmt.Sprintf(m.format, m.args...))
	})

	t.Run("non-variadic wrong count", func(t *testing.T) {
		x := tapeItem{FuncName: "foo", Params: []json.RawMessage{[]byte(`"a"`), []byte(`"b"`)}}
		m := &testingTTapeItemMatch{T: t}
		x.match(m, "foo", false, 3, "a", "b", []string{"c", "d"})
		assert.Equal(t, "wrong number of inputs: expected 2, got 3", fmt.Sprintf(m.format, m.args...))
	})

	t.Run("wrong output count", func(t *testing.T) {
		x := tapeItem{
			FuncName: "foo",
			Params:   []json.RawMessage{[]byte(`"a"`), []byte(`"b"`)},
			Results:  []json.RawMessage{[]byte("123")},
		}
		m := &testingTTapeItemMatch{T: t}
		x.match(m, "foo", false, 2, "a", "b", 123, 456)
		assert.Equal(t, "wrong number of outputs: expected 1, got 2", fmt.Sprintf(m.format, m.args...))
	})
}

func Test_tape_write(t *testing.T) {
	tests := []struct {
		name string

		vard bool
	}{
		{
			name: "variadic",
			vard: true,
		},
		{
			name: "non-variadic",
			vard: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := []interface{}{
				"a", "b", "c", []string{"d", "e"},
			}
			tapeObj := tape{}
			tapeObj.write("testing", tt.vard, items...)
			require.Len(t, tapeObj, 1)
			assert.Equal(t, "testing", tapeObj[0].FuncName)
			if tt.vard {
				assert.Equal(t, []json.RawMessage{[]byte(`"a"`),
					[]byte(`"b"`), []byte(`"c"`), []byte(`"d"`), []byte(`"e"`)}, tapeObj[0].Params)
			} else {
				assert.Equal(t, []json.RawMessage{[]byte(`"a"`),
					[]byte(`"b"`), []byte(`"c"`), []byte(`["d","e"]`)}, tapeObj[0].Params)
			}
		})
	}
}

func Test_tapeItem_end(t *testing.T) {
	tests := []struct {
		name string

		restError    *rest.ErrorREST
		genericError string
	}{
		{
			name: "no error",
		},
		{
			name: "rest error",
			restError: &rest.ErrorREST{
				Message: "message",
			},
		},
		{
			name:         "generic error",
			genericError: "error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the params.
			params := []interface{}{"a", "b", "c"}
			if tt.restError != nil {
				params = append(params, tt.restError)
			} else if tt.genericError != "" {
				params = append(params, errors.New(tt.genericError))
			}

			// Create the tape item and call the end function.
			ti := tapeItem{}
			ti.end(params...)

			// Check the results.
			assert.Equal(t, []json.RawMessage{[]byte(`"a"`), []byte(`"b"`), []byte(`"c"`)}, ti.Results)

			// Check the error.
			if tt.restError != nil {
				assert.Equal(t, tt.restError, ti.RESTError)
			} else if tt.genericError != "" {
				assert.Equal(t, tt.genericError, ti.GenericError)
			} else {
				assert.Nil(t, ti.RESTError)
				assert.Equal(t, "", ti.GenericError)
			}
		})
	}
}
