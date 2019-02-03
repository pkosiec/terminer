package sh_test

import (
	"testing"

	"github.com/pkosiec/terminer/internal/sh"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {
	s := sh.New()
	out, err := s.Exec("echo 'Foo'")
	require.NoError(t, err)
	assert.Equal(t, "Foo\n", out)
}

func TestExecInDir(t *testing.T) {
	s := sh.New()
	out, err := s.ExecInDir("pwd", "/")
	require.NoError(t, err)
	assert.Equal(t, "/\n", out)
}
