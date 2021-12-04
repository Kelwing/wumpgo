package router

import (
	"encoding/json"
	"github.com/Postcord/rest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Postcord/objects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestingT is our internal requirements from *testing.T. The weird edgecase is Run since the return type can be different.
type TestingT interface {
	require.TestingT
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Log(args ...interface{})
	Helper()
}

// The more generic runner.
type genericRunner interface {
	Run(name string, f func(t TestingT)) bool
}

// TestComponent is used to run unit tests against the specified component.
func TestComponent(t TestingT, b LoaderBuilder, path string) {
	// Get everything we need from the loader.
	r, _, errHandler, restOrigin, allowedMentions := b.CurrentChain()

	// Get the Postcord regen env var.
	regen := 0
	switch strings.ToLower(os.Getenv("POSTCORD_REGEN")) {
	case "all":
		// 1 | 2
		regen = 3
	case "response":
		regen = 1
	}

	// Make sure the component router isn't nil.
	require.NotNil(t, r)

	// Get the filesystem friendly version of the path.
	fsSafePath := strings.ReplaceAll(path, "/", "_")

	// Create the folder path.
	folderPath := filepath.Join("testframes", "components", fsSafePath)

	// List the file contents of the folder.
	fs, err := os.ReadDir(folderPath)
	if err == os.ErrNotExist {
		// Just return here. There's nothing to actually look at.
		return
	}

	// Make sure there were no errors.
	require.NoError(t, err)

	// Loop through the files.
	for _, f := range fs {
		// Load the file.
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".json") {
			// Skip directories and non-json files.
			continue
		}
		fp := filepath.Join(folderPath, f.Name())
		b, err := os.ReadFile(fp)
		require.NoError(t, err)
		var frameData frame
		err = json.Unmarshal(b, &frameData)
		if err != nil {
			t.Log("unable to run", f.Name(), "because of a json error:", err)
			continue
		}

		// Defines the error handler.
		var returnedErr error
		respExpected := true
		errHandlerOverride := func(err error) *objects.InteractionResponse {
			returnedErr = err
			if errHandler != nil {
				return errHandler(err)
			}
			respExpected = false
			return nil
		}

		// Define the test.
		test := func(t TestingT) {
			// Create the rest client.
			var restClient rest.RESTClient
			if regen&2 == 0 {
				restClient = &restTapePlayer{
					t:    t,
					tape: frameData.RESTRequests,
				}
			} else {
				restClient = &restTape{
					rest: restOrigin,
					tape: &tape{},
				}
			}

			// Create the components handler.
			handler := r.build(loaderPassthrough{
				rest:                  restClient,
				errHandler:            errHandlerOverride,
				globalAllowedMentions: allowedMentions,
				generateFrames:        false,
			})

			// Run the handler.
			resp := handler(frameData.Request)

			// Write the tape.
			if regen&2 != 0 {
				frameData.RESTRequests = *restClient.(*restTape).tape
			}

			// Handle the response.
			if regen&1 == 0 {
				// Compare the data we get back.
				if frameData.Error == "" {
					assert.NoError(t, returnedErr)
				} else {
					assert.EqualError(t, returnedErr, frameData.Error)
				}
				if respExpected {
					assert.Equal(t, frameData.Response, resp)
				}
			} else {
				// Handle setting the response.
				frameData.Response = resp
				frameData.Error = ""
				if returnedErr != nil {
					frameData.Error = returnedErr.Error()
				}
			}

			// Write the response if it was a regen.
			if regen != 0 {
				b := mustMarshal(t, true, frameData)
				require.NoError(t, os.WriteFile(fp, b, 0644))
			}
		}

		// Run the defined test.
		if runner, ok := t.(genericRunner); ok {
			runner.Run(f.Name(), test)
		} else {
			t.(*testing.T).Run(f.Name(), func(t *testing.T) {
				t.Helper()
				test(t)
			})
		}
	}
}

func testCommand(t TestingT, b LoaderBuilder, autocomplete bool, commandRoute ...string) {
	// Get everything we need from the loader.
	_, r, errHandler, restOrigin, allowedMentions := b.CurrentChain()

	// Get the Postcord regen env var.
	regen := 0
	switch strings.ToLower(os.Getenv("POSTCORD_REGEN")) {
	case "all":
		// 1 | 2
		regen = 3
	case "response":
		regen = 1
	}

	// Make sure the command router isn't nil.
	require.NotNil(t, r)

	// Create the folder path.
	pathParts := []string{"testframes", "commands"}
	if autocomplete {
		pathParts = []string{"testframes", "autocompletes"}
	}
	pathParts = append(pathParts, commandRoute...)
	folderPath := filepath.Join(pathParts...)

	// List the file contents of the folder.
	fs, err := os.ReadDir(folderPath)
	if os.IsNotExist(err) {
		// Just return here. There's nothing to actually look at.
		return
	}

	// Make sure there were no errors.
	require.NoError(t, err)

	// Loop through the files.
	for _, f := range fs {
		// Load the file.
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".json") {
			// Skip directories and non-json files.
			continue
		}
		fp := filepath.Join(folderPath, f.Name())
		b, err := os.ReadFile(fp)
		require.NoError(t, err)
		var frameData frame
		err = json.Unmarshal(b, &frameData)
		if err != nil {
			t.Log("unable to run", f.Name(), "because of a json error:", err)
			continue
		}

		// Defines the error handler.
		var returnedErr error
		respExpected := true
		errHandlerOverride := func(err error) *objects.InteractionResponse {
			returnedErr = err
			if errHandler != nil {
				return errHandler(err)
			}
			respExpected = false
			return nil
		}

		// Define the test.
		test := func(t TestingT) {
			// Defines the helper.
			t.Helper()

			// Create the rest client.
			var restClient rest.RESTClient
			if regen&2 == 0 {
				restClient = &restTapePlayer{
					t:    t,
					tape: frameData.RESTRequests,
				}
			} else {
				restClient = &restTape{
					rest: restOrigin,
					tape: &tape{},
				}
			}

			// Create the handler.
			cmdHandler, autoCompleteHandler := r.build(loaderPassthrough{
				rest:                  restClient,
				errHandler:            errHandlerOverride,
				globalAllowedMentions: allowedMentions,
				generateFrames:        false,
			})

			// Run the handler.
			var resp *objects.InteractionResponse
			if autocomplete {
				resp = autoCompleteHandler(frameData.Request)
			} else {
				resp = cmdHandler(frameData.Request)
			}

			// Write the tape.
			if regen&2 != 0 {
				frameData.RESTRequests = *restClient.(*restTape).tape
			}

			// Handle the response.
			if regen&1 == 0 {
				// Compare the data we get back.
				if frameData.Error == "" {
					assert.NoError(t, returnedErr)
				} else {
					assert.EqualError(t, returnedErr, frameData.Error)
				}
				if respExpected {
					assert.Equal(t, frameData.Response, resp)
				}
			} else {
				// Handle setting the response.
				frameData.Response = resp
				frameData.Error = ""
				if returnedErr != nil {
					frameData.Error = returnedErr.Error()
				}
			}

			// Write the response if it was a regen.
			if regen != 0 {
				b := mustMarshal(t, true, frameData)
				require.NoError(t, os.WriteFile(fp, b, 0644))
			}
		}

		// Run the defined test.
		if runner, ok := t.(genericRunner); ok {
			runner.Run(f.Name(), test)
		} else {
			t.(*testing.T).Run(f.Name(), func(t *testing.T) {
				t.Helper()
				test(t)
			})
		}
	}
}

// TestCommand is used to run unit tests against the specified command.
func TestCommand(t TestingT, b LoaderBuilder, commandRoute ...string) {
	testCommand(t, b, false, commandRoute...)
}

// TestAutocomplete is used to run unit tests against the specified commands auto-complete.
func TestAutocomplete(t TestingT, b LoaderBuilder, commandRoute ...string) {
	testCommand(t, b, true, commandRoute...)
}
