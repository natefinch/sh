package sh_test

import (
	"fmt"

	"github.com/natefinch/sh"
)

func ExampleCmd() {
	echo := sh.Cmd("echo")

	fmt.Print(echo("Hi there!"))
	// output:
	// Hi there!
}

func ExamplePipe() {
	echo := sh.Cmd("echo")
	wc := sh.Cmd("wc", "-w")
	grep := sh.Cmd("grep", "-o")

	fmt.Print(sh.Pipe(echo("Hi there!!"), grep("Hi"), wc()))
	// output:
	// 1
}

func ExampleCmd_BakedIn() {
	echo := sh.Cmd("echo")

	// You can bake-in arguments when you create the command, the arguments are
	// then always used when we use the Command.
	wc := sh.Cmd("wc", "-w")

	fmt.Print(sh.Pipe(echo("Hi there!"), wc()))
	// output:
	// 2
}

func Example_String() {
	echo := sh.Cmd("echo")

	output := echo("Hi there!")

	prnt := func(s string) {
		fmt.Print(s)
	}

	// Since we're passing the output into a function expecting a string, we
	// have to call String() on the output of the command.
	prnt(output.String())
	// output:
	// Hi there!
}
