package objects

import "golang.org/x/exp/constraints"

func HasFlag[T constraints.Integer](bitset T, flag T) bool {
	return (bitset & flag) == flag
}
