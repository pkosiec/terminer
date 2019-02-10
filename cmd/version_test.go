package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunVersion(t *testing.T) {
	assert.NotPanics(t, func() {
		runVersion(nil, nil)
	})
}
