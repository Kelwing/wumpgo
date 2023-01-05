package router

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestHandler func()

func TestInsert(t *testing.T) {
	root := newTrieNode[TestHandler]()

	buf := &bytes.Buffer{}

	root.Insert("/role/{role_id}", func() {
		fmt.Fprint(buf, "testing")
	})
}

func TestSearch(t *testing.T) {
	root := newTrieNode[TestHandler]()

	buf := &bytes.Buffer{}

	root.Insert("/role/{role_id}", func() {
		fmt.Fprint(buf, "testing")
	})

	h, ph, ok := root.Search("/role/1234567")
	require.True(t, ok)

	h()
	require.Equal(t, "testing", buf.String())
	require.Contains(t, ph, "role_id")

	buf2 := &bytes.Buffer{}

	root.Insert("my_cool_custom_id", func() {
		fmt.Fprint(buf2, "testing2")
	})

	h, ph, ok = root.Search("my_cool_custom_id")
	require.True(t, ok)

	require.Empty(t, ph)

	h()

	require.Equal(t, "testing2", buf2.String())

	h, ph, ok = root.Search("does_not_exist")
	require.False(t, ok)

	require.Nil(t, h)
	require.Empty(t, ph)
}
