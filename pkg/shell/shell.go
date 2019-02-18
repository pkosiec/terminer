package shell

import (
	"fmt"
	"os/exec"
)

type Command struct {
	Run   string
	Shell string
	Root  bool
}

//go:generate mockery -name=Shell -output=automock -outpkg=automock -case=underscore
type Shell interface {
	Exec(command Command) (string, error)
}

func New() Shell {
	return &shell{}
}

const DefaultShell = "/bin/sh"

type shell struct{}

// Exec executes given command in specified shell
func (s *shell) Exec(command Command) (string, error) {
	if command.Shell == "" {
		command.Shell = DefaultShell
	}

	var cmd *exec.Cmd
	if command.Root {
		cmd = s.runAsRoot(command)
	} else {
		cmd = exec.Command(command.Shell, "-c", command.Run)
	}

	return s.runCmd(cmd)
}

func (s *shell) runCmd(cmd *exec.Cmd) (string, error) {
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (s *shell) runAsRoot(cmd Command) *exec.Cmd {
	if !s.isCommandAvailable("sudo") {
		return exec.Command("su", "-s", cmd.Shell, "-c", cmd.Run)
	}

	return exec.Command("sudo", cmd.Shell, "-c", cmd.Run)
}

func (s *shell) isCommandAvailable(cmdName string) bool {
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("command -v %s", cmdName))
	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}
