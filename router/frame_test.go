package router

import (
	"os"
	"path/filepath"
	"testing"

	"wumpgo.dev/wumpgo/objects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_frame_write(t *testing.T) {
	// Get the folder path.
	folderPath := []string{t.TempDir(), "folder"}
	folderPathJoined := filepath.Join(folderPath...)

	// Create the frame.
	f := &frame{
		Request: &objects.Interaction{Type: 69},
		RESTRequests: tape{
			{
				FuncName:     "test",
				GenericError: "hello world!",
			},
		},
		Error:    "test!",
		Response: &objects.InteractionResponse{},
	}
	f.write(folderPath...)

	// Check the folder exists.
	x, err := os.ReadDir(folderPathJoined)
	require.NoError(t, err)
	require.Len(t, x, 1)
	assert.False(t, x[0].IsDir())
}
