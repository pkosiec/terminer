package shell_test

import (
	"testing"

	"github.com/pkosiec/terminer/pkg/shell"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShell_Exec(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		s := shell.New()
		out, err := s.Exec(shell.Command{
			Run:  "echo 'Foo'",
			Root: false,
		})
		require.NoError(t, err)
		assert.Equal(t, "Foo\n", out)
	})

	t.Run("With Custom Shell", func(t *testing.T) {
		s := shell.New()
		out, err := s.Exec(shell.Command{
			Run:   "echo 'Foo'",
			Shell: "sh",
		})
		require.NoError(t, err)
		assert.Equal(t, "Foo\n", out)
	})
}
