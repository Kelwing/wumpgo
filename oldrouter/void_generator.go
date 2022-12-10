package router

import "strconv"

// Used to generate void paths.
type voidGenerator int

// VoidCustomID is used to return a unique custom ID for this context that resolves to a void.
func (g *voidGenerator) VoidCustomID() string {
	*g++
	return "/_postcord/void/" + strconv.Itoa(int(*g))
}
