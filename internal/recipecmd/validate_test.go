package recipecmd_test

import (
	"fmt"
	"github.com/pkosiec/terminer/internal/recipecmd"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateArgs(t *testing.T) {
	filePathBak := recipecmd.FilePath
	urlBak := recipecmd.URL

	testCases := []struct{
		FilePath string
		URL string
		args []string
		expectedErr bool
	}{
		{
			FilePath: "./test.md",
			expectedErr: false,
		},
		{
			URL: "https://example.com",
			expectedErr: false,
		},
		{
			args: []string{"test-recipe"},
			expectedErr: false,
		},
		{
			args: []string{"test-recipe", "test-recipe2"},
			expectedErr: true,
		},
		{
			expectedErr: true,
		},
	}

	for tN, tC := range testCases {
		t.Run(fmt.Sprintf("Test Case %d", tN), func(t *testing.T) {
			recipecmd.URL = tC.URL
			recipecmd.FilePath = tC.FilePath
			err := recipecmd.ValidateArgs(nil, tC.args)

			if tC.expectedErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}

	recipecmd.FilePath = filePathBak
	recipecmd.URL = urlBak
}
