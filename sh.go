// Package sh is intended to make working with shell commands more shell-like.
// This package is basically just syntactic sugar wrapped around os/exec, but it
// can make your code a lot easier to read when you have a simple shell-like
// line rather than a huge mess of pipes and commands.
//
//	echo := sh.Cmd("echo")
//
//	fmt.Print(echo("Hi there!"))
//	// output:
//	// Hi there!
//
//	grep := sh.Cmd("grep", "-o")
//	wc := sh.Cmd("wc", "-w")
//
//	fmt.Print(sh.Pipe(echo("Hi there!"), grep("Hi"), wc()))
//	// output:
//	// 1
package sh

import (
	"bufio"
	"os/exec"
)

// Cmd returns a function that runs the given command with the given args.  The
// args that are passed to Cmd cannot be overridden, but the function that is
// returned can add more arguments when you call it.
func Cmd(name string, args ...string) func(args ...string) Command {
	c := exec.Command(name, args...)

	f := func(args ...string) Command {
		c.Args = append(c.Args, args...)

		return func(args ...string) (string, error) {
			if len(args) > 0 {
				p, _ := c.StdinPipe()
				b := bufio.NewWriter(p)
				go func() { b.WriteString(args[0]); b.Flush(); p.Close() }()
			}
			out, err := c.CombinedOutput()
			return string(out), err
		}
	}

	return f
}

// Pipe connects the output of one Command to the input of the next Command,
// returning immediately if any of the commands produce an error.  The result is
// a function that returns the output of the last function run, and any error it
// might have had.  Alternatively, you can call String() on the result to get
// just the output. You may also pass the result into a fmt.Print-style function
// as if it were a string.
func Pipe(cmds ...Command) pipeResult {
	stdin := ""
	var err error
	for _, cmd := range cmds {
		stdin, err = cmd(stdin)
		if err != nil {
			break
		}
	}
	return func() (string, error) { return stdin, err }
}

// Command is the type that is returned from running a Cmd.  You can pass it
// into a fmt.Print style function as if it were a string and it'll do the right
// thing, or you can execute it to get the output and any errors, or call
// String() to just get the output.
type Command func(arg ...string) (string, error)

// String implements Stringer.String
func (c Command) String() string {
	s, _ := c()
	return s
}

// pipeResult is a function that you can call to get a string and an error, or
// call String() to get just the string.
type pipeResult func() (string, error)

// String implements Stringer.String
func (p pipeResult) String() string {
	s, _ := p()
	return s
}
