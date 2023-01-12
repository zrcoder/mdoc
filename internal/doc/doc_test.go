package doc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_joinUrlPath(t *testing.T) {
	tests := []struct {
		name string
		base string
		last string
		want string
	}{
		{
			name: "regular",
			base: "/docs",
			last: "1-hello-world",
			want: "/docs/1-hello-world",
		},
		{
			name: "contains blanks",
			base: "/docs",
			last: "1. Hello world",
			want: "/docs/1.%20Hello%20world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, joinUrlPath(tt.base, tt.last), "joinUrlPath(%s, %s)", tt.base, tt.last)
		})
	}
}
