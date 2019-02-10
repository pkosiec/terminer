package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateInstallRollbackArgs(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		err := validateInstallRollbackArgs(nil, []string{"path/path"})

		assert.NoError(t, err)
	})

	t.Run("No Arguments", func(t *testing.T) {
		err := validateInstallRollbackArgs(nil, []string{})

		assert.Error(t, err)
	})
}
