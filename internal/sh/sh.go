package sh

import "os/exec"

//go:generate mockery -name=Sh -output=automock -outpkg=automock -case=underscore
type Sh interface {
	Exec(command string) (string, error)
	ExecInDir(command, dir string) (string, error)
}

func New() Sh {
	return &sh{}
}

type sh struct{}

// Exec executes given command in Bourne shell
func (s *sh) Exec(command string) (string, error) {
	cmd := s.prepareCmd(command)
	return s.runCmd(cmd)
}

// RunInDir runs given command in shell in specific directory
func (s *sh) ExecInDir(command, dir string) (string, error) {
	cmd := s.prepareCmd(command)
	cmd.Dir = dir
	return s.runCmd(cmd)
}

func (s *sh) prepareCmd(command string) *exec.Cmd {
	return exec.Command("/bin/sh", "-c", command)
}

func (s *sh) runCmd(cmd *exec.Cmd) (string, error) {
	out, err := cmd.CombinedOutput()
	return string(out), err
}
