package shell_test

import (
	"testing"

	"github.com/pkosiec/terminer/internal/shell"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShell_Exec(t *testing.T) {
	s := shell.New()
	out, err := s.Exec(shell.Command{
		Run:  "echo 'Foo'",
		Sudo: false,
	})
	require.NoError(t, err)
	assert.Equal(t, "Foo\n", out)
}
