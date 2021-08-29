package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_voidGenerator_VoidCustomID(t *testing.T) {
	tests := [] struct{
		name string

		runs  int
		wants string
	}{
		{
			name:  "single run",
			runs:  1,
			wants: "/_postcord/void/1",
		},
		{
			name:  "multiple runs",
			runs:  10,
			wants: "/_postcord/void/10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := voidGenerator(0)
			res := ""
			for i := 0; i < tt.runs; i++ {
				res = v.VoidCustomID()
			}
			assert.Equal(t, tt.wants, res)
		})
	}
}
