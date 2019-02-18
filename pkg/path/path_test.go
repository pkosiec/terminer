package path_test

import (
	"github.com/pkosiec/terminer/pkg/path"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsURL(t *testing.T) {
	cases := []struct {
		path           string
		expectedResult bool
	}{
		{path: "https://test.com/file.yaml", expectedResult: true},
		{path: "https://test.com", expectedResult: true},
		{path: "http://example.pl", expectedResult: true},
		{path: "ftp://example.com/test.yaml", expectedResult: true},
		{path: "../example.yaml", expectedResult: false},
		{path: "/test/example.yml", expectedResult: false},
		{path: "/users/test/example.yml", expectedResult: false},
		{path: "example.yml", expectedResult: false},
		{path: "./foo.yaml", expectedResult: false},
	}

	for tN, tC := range cases {
		t.Logf("Test case %d: %s should be %t", tN+1, tC.path, tC.expectedResult)
		result := path.IsURL(tC.path)
		assert.Equal(t, tC.expectedResult, result)
	}
}
