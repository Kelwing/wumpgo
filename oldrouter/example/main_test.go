package main

import (
	"testing"

	"wumpgo.dev/wumpgo/router"
)

func TestComponent_set(t *testing.T) {
	_, _, b := builder()
	router.TestComponent(t, b, "/set/:number")
}

func Test_subgroups_group2_autocomplete(t *testing.T) {
	_, _, b := builder()
	router.TestCommand(t, b, "subgroups", "group2", "autocomplete")
}

func Test_Autocomplete_subgroups_group2_autocomplete(t *testing.T) {
	_, _, b := builder()
	router.TestAutocomplete(t, b, "subgroups", "group2", "autocomplete")
}

func Test_rest(t *testing.T) {
	_, _, b := builder()
	router.TestCommand(t, b, "rest")
}
