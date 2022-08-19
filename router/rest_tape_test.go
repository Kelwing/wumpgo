package router

import (
	"encoding/json"
	"reflect"
	"testing"

	"wumpgo.dev/wumpgo/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ rest.RESTClient = (*restTape)(nil)

func Test_restTape(t *testing.T) {
	// Get all methods we wish to mock.
	refImpl := reflect.TypeOf((*rest.Client)(nil))
	numMethods := refImpl.NumMethod()
	methods := make([]reflect.Method, numMethods)
	for i := 0; i < numMethods; i++ {
		methods[i] = refImpl.Method(i)
	}

	// Make sub-tests for all the methods.
	for _, v := range methods {
		t.Run(v.Name, func(t *testing.T) {
			// Define the tape item.
			item := tapeItem{FuncName: v.Name}

			// Create zero values for each of the params.
			inCount := v.Type.NumIn()
			if v.Type.IsVariadic() {
				inCount--
			}
			in := make([]json.RawMessage, inCount-1)
			zeroIn := make([]reflect.Value, inCount-1)
			for i := 1; i < inCount; i++ {
				x := v.Type.In(i)
				n := reflect.New(x)
				b, err := json.Marshal(n.Interface())
				require.NoError(t, err)
				in[i-1] = b
				zeroIn[i-1] = n.Elem()
			}
			item.Params = in

			// Create zero values for the return values.
			outCount := v.Type.NumOut()
			if outCount > 0 {
				o := v.Type.Out(outCount - 1)
				if _, ok := reflect.New(o).Interface().(*error); ok {
					outCount--
				}
			}
			out := make([]json.RawMessage, outCount)
			for i := 0; i < outCount; i++ {
				x := v.Type.Out(i)
				b, err := json.Marshal(reflect.New(x).Interface())
				require.NoError(t, err)
				out[i] = b
			}
			item.Results = out

			// Create the tape recorder taking in this specific item as a REST client.
			tapeObj := tape{}
			rec := &restTape{
				tape: &tapeObj,
				rest: &restTapePlayer{t: t, tape: tape{&item}},
			}

			// Call the method.
			m := reflect.ValueOf(rec).MethodByName(v.Name)
			if m.IsZero() {
				t.Fatalf("method %s not found", v.Name)
			}
			res := m.Call(zeroIn)

			// Check the results.
			for _, x := range res {
				if !x.IsZero() {
					t.Errorf("%s is not zero", x.Type())
				}
			}

			// Check the tape has the expected result.
			require.Len(t, tapeObj, 1)
			assert.Equal(t, &tapeItem{
				FuncName: v.Name,
				Params:   in,
				Results:  out,
			}, tapeObj[0])
		})
	}
}
