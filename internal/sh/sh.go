package sh

import "os/exec"

// Exec executes given command in Bourne shell
func Exec(command string) (string, error) {
	cmd := prepareCmd(command)
	return runCmd(cmd)
}

// RunInDir runs given command in shell in specific directory
func ExecInDir(command, dir string) (string, error) {
	cmd := prepareCmd(command)
	cmd.Dir = dir
	return runCmd(cmd)
}

func prepareCmd(command string) *exec.Cmd {
	return exec.Command("/bin/sh", "-c", command)
}

func runCmd(cmd *exec.Cmd) (string, error) {
	out, err := cmd.CombinedOutput()
	return string(out), err
}