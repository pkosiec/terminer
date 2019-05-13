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
		cmdPrinter := func(s string) {
			assert.Equal(t, "echo 'Foo'", s)
		}
		outPrinter := func(s string) {
			assert.Equal(t, "Foo", s)
		}
		errPrinter := func(s string) {
			assert.Empty(t, s)
		}

		s := shell.New(cmdPrinter, outPrinter, errPrinter)

		err := s.Exec(shell.Command{
			Run: []string{
				"echo 'Foo'",
			},
			Root: false,
		}, true)
		require.NoError(t, err)
	})

	t.Run("With Custom Shell", func(t *testing.T) {
		cmdPrinter := func(s string) {
			assert.Equal(t, "echo 'Foo'", s)
		}
		outPrinter := func(s string) {
			assert.Equal(t, "Foo", s)
		}
		errPrinter := func(s string) {
			assert.Empty(t, s)
		}

		s := shell.New(cmdPrinter, outPrinter, errPrinter)
		err := s.Exec(shell.Command{
			Run: []string{
				"echo 'Foo'",
			},
			Shell: "bash",
		}, true)
		require.NoError(t, err)
	})

	t.Run("Print errors", func(t *testing.T) {
		cmdPrinter := func(s string) {
			assert.Equal(t, ">&2 echo 'error!'", s)
		}
		outPrinter := func(s string) {
			assert.Empty(t, s)
		}
		errPrinter := func(s string) {
			assert.Equal(t, "error!", s)
		}

		s := shell.New(cmdPrinter, outPrinter, errPrinter)

		err := s.Exec(shell.Command{
			Run: []string{
				">&2 echo 'error!'",
			},
			Root: false,
		}, true)
		require.NoError(t, err)
	})

	t.Run("Run multiple commands", func(t *testing.T) {
		cmdPrinter := printerAssertFn(t, func(i int) string {
			switch i {
			case 0:
				return "echo 'Foo'"
			case 1:
				return "echo 'Bar'"
			case 2:
				return ">&2 echo 'Error here'"
			}

			return ""
		})
		outPrinter := printerAssertFn(t, func(i int) string {
			switch i {
			case 0:
				return "Foo"
			case 1:
				return "Bar"
			case 2:
				return ""
			}

			return ""
		})
		errPrinter := printerAssertFn(t, func(i int) string {
			switch i {
			case 0:
				return "Error here"
			}

			return ""
		})

		s := shell.New(cmdPrinter, outPrinter, errPrinter)

		err := s.Exec(shell.Command{
			Run: []string{
				"echo 'Foo'",
				"echo 'Bar'",
				">&2 echo 'Error here'",
			},
			Root: false,
		}, true)
		require.NoError(t, err)
	})

	t.Run("Stop on error", func(t *testing.T) {
		cmdPrinter := printerAssertFn(t, func(i int) string {
			switch i {
			case 0:
				return "echo 'Foo'"
			case 1:
				return "exit 1"
			}

			return ""
		})
		outPrinter := printerAssertFn(t, func(i int) string {
			switch i {
			case 0:
				return "Foo"
			case 1:
				return "Bar"
			}

			return ""
		})
		errPrinter := func(s string) {
			assert.Fail(t, "Should not be called")
		}

		s := shell.New(cmdPrinter, outPrinter, errPrinter)

		err := s.Exec(shell.Command{
			Run: []string{
				"echo 'Foo'",
				"exit 1",
				"echo 'Bar'",
			},
			Root: false,
		}, true)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "while executing exit 1")
	})

	t.Run("Do not stop on error", func(t *testing.T) {
		cmdPrinter := printerAssertFn(t, func(i int) string {
			switch i {
			case 0:
				return "echo 'Foo'"
			case 1:
				return "exit 1"
			case 2:
				return "echo 'Bar'"
			}

			return ""
		})
		outPrinter := printerAssertFn(t, func(i int) string {
			switch i {
			case 0:
				return "Foo"
			case 1:
				return "Bar"
			}

			return ""
		})
		errPrinter := func(s string) {
			assert.Fail(t, "Should not be called")
		}

		s := shell.New(cmdPrinter, outPrinter, errPrinter)

		err := s.Exec(shell.Command{
			Run: []string{
				"echo 'Foo'",
				"exit 1",
				"echo 'Bar'",
			},
			Root: false,
		}, false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "while executing exit 1")
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

func printerAssertFn(t *testing.T, expectedStringFn func(i int) string) func(s string) {
	var i int

	return func(s string) {
		assert.Equal(t, expectedStringFn(i), s)
		i++
	}
}
