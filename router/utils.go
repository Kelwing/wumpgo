package router

import "wumpgo.dev/wumpgo/objects"

func Chunk[T any](s []T, size int) [][]T {
	var divided [][]T

	for i := 0; i < len(s); i += size {
		end := i + size

		if end > len(s) {
			end = len(s)
		}

		divided = append(divided, s[i:end])
	}

	return divided
}

func ComponentsToRows(in []*objects.Component) []*objects.Component {
	chunks := Chunk(in, 5)

	componentRows := make([]*objects.Component, len(chunks))

	for i, c := range chunks {
		componentRows[i] = &objects.Component{
			Type:       objects.ComponentTypeActionRow,
			Components: c,
		}
	}

	return componentRows
}
