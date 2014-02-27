// Package sh is intended to make working with shell commands more shell-like.
// This package is basically just syntactic sugar wrapped around
// labix.org/v2/pipe, which in turn just wraps os/exec, but it can make your
// code a lot easier to read when you have a simple shell-like line rather than
// a huge mess of pipes and commands and conversions etc.
//
//	// create functions that runs the echo, grep, and wc shell commands
//	echo := sh.Cmd("echo")
//	grep := sh.Cmd("grep")
//	wc := sh.Cmd("wc")
//
//	// run echo, pipe the output through grep and then through wc
//	// effectively the same as
//	// $ echo Hi there! | grep -o Hi | wc -w
//	fmt.Print(sh.Pipe(echo("Hi there!"), grep("-o", "Hi"), wc("-w")))
//
//	// output:
//	// 1
package sh

import (
	"io"
	"strings"

	"labix.org/v2/pipe"
)

// Cmd returns a function that will return an Executable for the given command
// with the given args.  This is a convenience for defining functions for
// commonly reused Executables, such as grep, ls, mkdir, etc.
//
// The args that are passed to Cmd are passed to the Executable when the
// returned function is run, allowing you to pre-set some common arguments.
func Cmd(name string, args0 ...string) func(args ...string) Executable {
	return func(args1 ...string) Executable {
		return Executable{pipe.Exec(name, append(args0, args1...)...)}
	}
}

// Runner returns a function that will run the given shell command with
// specified arguments. This is a convenience for creating one-off commands that
// aren't going to be put into Pipe, but instead just run standalone.
//
// The args that are passed to Runner are passed to the shell command when the
// returned function is run, allowing you to pre-set some common arguments.
func Runner(name string, args0 ...string) func(args ...string) string {
	return func(args1 ...string) string {
		return Executable{pipe.Exec(name, append(args0, args1...)...)}.String()
	}
}

// Dump returns an excutable that will read the given file and dump its contents
// as the Executable's stdout.
func Dump(filename string) Executable {
	return Executable{pipe.ReadFile(filename)}
}

// Read returns an executable that will read from the given reader and use it as
// the Executable's stdout.
func Read(r io.Reader) Executable {
	return Executable{pipe.Read(r)}
}

// Pipe connects the output of one Executable to the input of the next
// Executable in the list.  The result is an Executable that, when run, returns
// the output of the last Executable run, and any error it might have had.
//
// If any of the Executables fails, no further Executables are run, and the
// failing Executable's stderr and error are returned.
func Pipe(cmds ...Executable) Executable {
	ps := make([]pipe.Pipe, len(cmds))
	for i, c := range cmds {
		ps[i] = c.Pipe
	}
	return Executable{pipe.Line(ps...)}
}

// PipeWith functions like Pipe, but runs the first command with stdin as the
// input.
func PipeWith(stdin string, cmds ...Executable) Executable {
	ps := make([]pipe.Pipe, len(cmds)+1)
	ps[0] = pipe.Read(strings.NewReader(stdin))
	for i, c := range cmds {
		ps[i+1] = c.Pipe
	}
	return Executable{pipe.Line(ps...)}
}

// Executable is a runnable construct.  You can run it by calling Run(), or by
// calling String() (which is automatically done when passing it into a
// fmt.Print style function).  It can be passed into Pipe to form a chain of
// Executables that are executed in series.
type Executable struct {
	pipe.Pipe
}

// Run executes the command with the given string as standard input, and returns
// stdout and a nil error on success, or stderr and a non-nil error on failure.
func (c Executable) Run(stdin string) (string, error) {
	stdout, stderr, err := pipe.DividedOutput(
		pipe.Line(pipe.Read(strings.NewReader(stdin)), c.Pipe),
	)
	if err != nil {
		return string(stderr), err
	}
	return string(stdout), nil
}

// String runs the Executable and returns the standard output as a string,
// ignoring any error.  This is most useful for passing an executable into a
// fmt.Print style function.
func (c Executable) String() string {
	s, _ := pipe.Output(c.Pipe)
	return string(s)
}
