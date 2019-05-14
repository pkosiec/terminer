package cmd_test

import (
	"github.com/pkosiec/terminer/cmd"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrintVersion(t *testing.T) {
	assert.NotPanics(t, func() {
		cmd.PrintVersion(nil, nil)
	})
}
