package shell

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

// PrintFn prints command output
type PrintFn func (string)

// Command represents command to execute in given shell
type Command struct {
	Run   string
	Shell string
	Root  bool
}

// Shell gives an ability to run shell commands
//go:generate mockery -name=Shell -output=automock -outpkg=automock -case=underscore
type Shell interface {
	Exec(command Command, outputPrinter, errPrinter PrintFn) error
}

// New creates a new instance that implements Shell interface
func New() Shell {
	return &shell{}
}

// DefaultShell defines in which shell all commands should be executed by default
const DefaultShell = "/bin/sh"

type shell struct{}

// Exec executes given command in specified shell
func (s *shell) Exec(command Command, outputPrinter, errPrinter PrintFn) error {
	if command.Shell == "" {
		command.Shell = DefaultShell
	}

	var cmd *exec.Cmd
	if command.Root {
		cmd = s.rootCommand(command)
	} else {
		cmd = exec.Command(command.Shell, "-c", command.Run)
	}

	return s.runCmd(cmd, outputPrinter, errPrinter)
}

func (s *shell) runCmd(cmd *exec.Cmd, outputPrinter, errPrinter PrintFn) error {
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	s.preparePipeScan(stdOut, outputPrinter)
	s.preparePipeScan(stdErr, errPrinter)

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	return err
}

func (s *shell) preparePipeScan(pipe io.ReadCloser, printer PrintFn) {
	scanner := bufio.NewScanner(pipe)
	scanner.Split(bufio.ScanLines)

	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			printer(text)
		}
	}()
}

// TODO: How to test it?
func (s *shell) rootCommand(cmd Command) *exec.Cmd {
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
