package rest

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_generateInterface(t *testing.T) {
	s := generateInterface()
	if os.Getenv("GENERATE") == "1" {
		assert.NoError(t, os.WriteFile("interface_gen.go", []byte(s), 0666))
	} else {
		b, err := os.ReadFile("interface_gen.go")
		assert.NoError(t, err)
		assert.Equal(t, s, string(b))
	}
}
