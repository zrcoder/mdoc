package osutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsFile(t *testing.T) {
	tests := []struct {
		path   string
		expVal bool
	}{
		{
			path:   "osutil.go",
			expVal: true,
		}, {
			path:   "../osutil",
			expVal: false,
		}, {
			path:   "not_found",
			expVal: false,
		},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, test.expVal, IsFile(test.path))
		})
	}
}

func TestIsDir(t *testing.T) {
	tests := []struct {
		path   string
		expVal bool
	}{
		{
			path:   "osutil.go",
			expVal: false,
		}, {
			path:   "../osutil",
			expVal: true,
		}, {
			path:   "not_found",
			expVal: false,
		},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, test.expVal, IsDir(test.path))
		})
	}
}

func TestIsExist(t *testing.T) {
	tests := []struct {
		path   string
		expVal bool
	}{
		{
			path:   "osutil.go",
			expVal: true,
		}, {
			path:   "../osutil",
			expVal: true,
		}, {
			path:   "not_found",
			expVal: false,
		},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, test.expVal, IsExist(test.path))
		})
	}
}
