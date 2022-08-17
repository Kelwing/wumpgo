package router

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/kelwing/wumpgo/objects"
)

// Defines all information relating to a Postcord data frame.
type frame struct {
	// Request is the interaction body that was sent by Discord.
	Request *objects.Interaction `json:"request"`

	// RESTRequests is the REST requests made during this interaction.
	RESTRequests tape `json:"rest_requests"`

	// Error is the string of the error that was encountered during this interaction.
	Error string `json:"error,omitempty"`

	// Response is the response that was sent back to Discord.
	Response *objects.InteractionResponse `json:"response"`
}

// Used to write the frame.
func (f *frame) write(subfolders ...string) {
	// Ensure the folder exists.
	joined := filepath.Join(subfolders...)
	if err := os.MkdirAll(joined, 0777); err != nil {
		panic(err)
	}

	// Defines the filename.
	filename := filepath.Join(joined, time.Now().In(time.UTC).Format("02-01-2006_15-04-03")+"_untitled_frame.json")
	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filename, b, 0666)
	if err != nil {
		panic(err)
	}
}
