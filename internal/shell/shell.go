package shell

import "os/exec"

//go:generate mockery -name=Shell -output=automock -outpkg=automock -case=underscore
type Shell interface {
	Exec(shell, command string) (string, error)
}

func New() Shell {
	return &shell{}
}

const DefaultShell = "sh"

type shell struct{}

// Exec executes given command in specified shell
func (s *shell) Exec(shell, command string) (string, error) {
	cmd := s.prepareCmd(shell, command)
	return s.runCmd(cmd)
}

func (s *shell) prepareCmd(shell, command string) *exec.Cmd {
	return exec.Command(shell, "-c", command)
}

func (s *shell) runCmd(cmd *exec.Cmd) (string, error) {
	out, err := cmd.CombinedOutput()
	return string(out), err
}
