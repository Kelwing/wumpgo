package scaffolding

import "strings"

func Bashify(in string) string {
	return strings.ReplaceAll(strings.ToLower(in), " ", "_")
}
