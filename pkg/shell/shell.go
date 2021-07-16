package shell

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os/exec"
	"strings"
)

// PrintFn prints command output
type PrintFn func(string)

// Command represents command to execute in given shell
type Command struct {
	Run   []string `yaml:"run" json:"run"`
	Shell string   `yaml:"shell" json:"shell"`
	Root  bool     `yaml:"root" json:"root"`
}

// Shell gives an ability to run shell commands
//go:generate mockery -name=Shell -output=automock -outpkg=automock -case=underscore
type Shell interface {
	Exec(command Command, stopOnError bool) error
}

// New creates a new instance that implements Shell interface
func New(printCmd PrintFn, printOut PrintFn, printErr PrintFn) Shell {
	return &shell{printCmd: printCmd, printOut: printOut, printErr: printErr}
}

// DefaultShell defines in which shell all commands should be executed by default
const DefaultShell = "/bin/sh"

type shell struct {
	printCmd PrintFn
	printOut PrintFn
	printErr PrintFn
}

// Exec executes given command in specified shell
func (s *shell) Exec(command Command, stopOnError bool) error {
	if command.Shell == "" {
		command.Shell = DefaultShell
	}

	var errMessages []string

	for _, singleCmd := range command.Run {
		var prefix string
		if command.Root {
			prefix = "$ "
		}

		s.printCmd(fmt.Sprintf("%s%s", prefix, singleCmd))

		var cmd *exec.Cmd
		if command.Root {
			cmd = s.rootCommand(command.Shell, singleCmd)
		} else {
			cmd = exec.Command(command.Shell, "-c", singleCmd)
		}

		err := s.runCmd(cmd)
		if err != nil {
			wrappedErr := errors.Wrapf(err, "while executing %s", singleCmd)
			if stopOnError {
				return wrappedErr
			}

			errMessages = append(errMessages, wrappedErr.Error())
		}
	}

	if len(errMessages) > 0 {
		return errors.New(strings.Join(errMessages, ",\n"))
	}

	return nil
}

func (s *shell) runCmd(cmd *exec.Cmd) error {
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	s.readAndPrint(stdOut, s.printOut)
	s.readAndPrint(stdErr, s.printErr)

	err = cmd.Wait()
	return err
}

func (s *shell) readAndPrint(pipe io.ReadCloser, printer PrintFn) {
	scanner := bufio.NewScanner(pipe)
	scanner.Split(bufio.ScanLines)

	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			printer(text)
		}
	}()
}

// TODO: Test it
func (s *shell) rootCommand(shell, cmd string) *exec.Cmd {
	if !s.isCommandAvailable("sudo") {
		return exec.Command("su", "-s", shell, "-c", cmd)
	}

	return exec.Command("sudo", shell, "-c", cmd)
}

func (s *shell) isCommandAvailable(cmdName string) bool {
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("command -v %s", cmdName))
	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}
