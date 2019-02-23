package shell_test

import (
	"fmt"
	"testing"

	"github.com/pkosiec/terminer/pkg/shell"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShell_Exec(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		s := shell.New()
		outPrinter := func(s string) {
			assert.Equal(t, "Foo", s)
		}
		errPrinter := func(s string) {}

		err := s.Exec(shell.Command{
			Run:  "echo 'Foo'",
			Root: false,
		}, outPrinter, errPrinter)
		require.NoError(t, err)
	})

	t.Run("With Custom Shell", func(t *testing.T) {
		s := shell.New()
		outPrinter := func(s string) {
			assert.Equal(t, "Foo", s)
		}
		errPrinter := func(s string) {}
		err := s.Exec(shell.Command{
			Run:   "echo 'Foo'",
			Shell: "sh",
		}, outPrinter, errPrinter)
		require.NoError(t, err)
	})
}

func TestShell_IsCommandAvailable(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		for _, testCase := range []string{"ls", "echo", "sh", "cd", "mkdir"} {
			t.Run(fmt.Sprintf("Test case %s", testCase), func(t *testing.T) {
				s := shell.ExposeInternalShell()
				result := s.IsCommandAvailable(testCase)
				assert.True(t, result)
			})
		}
	})

	t.Run("Error", func(t *testing.T) {
		s := shell.ExposeInternalShell()
		result := s.IsCommandAvailable("thiscommanddoesnotexist")

		require.False(t, result)
	})
}
