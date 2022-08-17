package router

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kelwing/wumpgo/rest"
)

var _ rest.RESTClient = (*restTapePlayer)(nil)

type restTapePlayerTestStoreFatal struct {
	*testing.T
	fatal string
}

func (t *restTapePlayerTestStoreFatal) Fatal(args ...any) {
	t.fatal = fmt.Sprint(args...)
}

func Test_restTapePlayer(t *testing.T) {
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
			call := func(t *testing.T, overrun bool) {
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

				// Create the player.
				var tester TestingT = t
				startIndex := 0
				if overrun {
					tester = &restTapePlayerTestStoreFatal{T: t}
					startIndex++
				}
				player := &restTapePlayer{
					t:     tester,
					index: startIndex,
					tape:  []*tapeItem{&item},
				}

				// Do reflection on the player.
				r := reflect.ValueOf(player)

				// Check the method exists and call it.
				m := r.MethodByName(v.Name)
				if m.IsZero() {
					t.Fatalf("method %s not found", v.Name)
				}
				res := m.Call(zeroIn)

				if overrun {
					// Check the fatal was thrown if this is an overrun.
					assert.Equal(t, "unexpected "+v.Name+" at end of tape", tester.(*restTapePlayerTestStoreFatal).fatal)
				} else {
					// Check the results.
					for _, x := range res {
						if !x.IsZero() {
							t.Errorf("%s is not zero", x.Type())
						}
					}
				}

				// Expect the index to be 1.
				assert.Equal(t, 1, player.index)
			}
			t.Run("success", func(t *testing.T) {
				call(t, false)
			})
			t.Run("overrun", func(t *testing.T) {
				call(t, true)
			})
		})
	}
}
