package shell

import "os/exec"

type Command struct {
	Run   string
	Shell string
	Sudo  bool
}

//go:generate mockery -name=Shell -output=automock -outpkg=automock -case=underscore
type Shell interface {
	Exec(command Command) (string, error)
}

func New() Shell {
	return &shell{}
}

const DefaultShell = "sh"

type shell struct{}

// Exec executes given command in specified shell
func (s *shell) Exec(command Command) (string, error) {
	shell := command.Shell
	if command.Shell == "" {
		shell = DefaultShell
	}

	var cmd *exec.Cmd
	if command.Sudo {
		cmd = exec.Command("sudo", shell, "-c", command.Run)
	} else {
		cmd = exec.Command(shell, "-c", command.Run)
	}

	return s.runCmd(cmd)
}

func (s *shell) runCmd(cmd *exec.Cmd) (string, error) {
	out, err := cmd.CombinedOutput()
	return string(out), err
}
