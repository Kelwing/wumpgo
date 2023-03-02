package router

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChunk(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	chunked := Chunk(data, 5)
	require.Len(t, chunked, 2)
	chunked = Chunk(data, 1)
	require.Len(t, chunked, 10)
}
