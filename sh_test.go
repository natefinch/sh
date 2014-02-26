package sh_test

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/natefinch/sh"
)

const SWCrawl = `
A long time ago, in a galaxy far, far away....

It is a period of civil war. Rebel
spaceships, striking from a hidden
base, have won their first victory
against the evil Galactic Empire.

During the battle, Rebel spies managed
to steal secret plans to the Empire's
ultimate weapon, the Death Star, an
armored space station with enough
power to destroy an entire planet.

Pursued by the Empire's sinister agents,
Princess Leia races home aboard her
starship, custodian of the stolen plans
that can save her people and restore
freedom to the galaxy...`

func ExampleCmd() {
	echo := sh.Cmd("echo")

	fmt.Print(echo("Hi there!"))
	// output:
	// Hi there!
}

func ExamplePipe() {
	echo := sh.Cmd("echo")

	// note, you can "bake in" arguments when you create these functions.
	upper := sh.Cmd("tr", "[:lower:]", "[:upper:]")
	grep := sh.Cmd("grep")

	// Equivalent of shell command:
	// $ echo Hi there! | grep -o Hi | wc -w
	fmt.Print(sh.Pipe(echo(SWCrawl), grep("far"), upper()))
	// output:
	// A LONG TIME AGO, IN A GALAXY FAR, FAR AWAY....
}

func ExamplePipeWith() {
	upper := sh.Cmd("tr", "[:lower:]", "[:upper:]")
	grep := sh.Cmd("grep")

	fmt.Print(sh.PipeWith(SWCrawl, grep("far"), upper()))
	// output:
	// A LONG TIME AGO, IN A GALAXY FAR, FAR AWAY....
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

func ExampleDump() {
	grep := sh.Cmd("grep")

	name := "ExampleDumpTest"
	defer writeTempFile(name, SWCrawl)()

	// Equivalent of shell command
	// $ cat ExampleDumpTest | grep far
	fmt.Print(sh.Pipe(sh.Dump(name), grep("far")))
	// output:
	// A long time ago, in a galaxy far, far away....
}

func ExampleRead() {
	grep := sh.Cmd("grep")

	name := "ExampleReadTest"
	f, cleanup := openTempFile(name, SWCrawl)
	defer cleanup()

	fmt.Print(sh.Pipe(sh.Read(f), grep("far")))
	// output:
	// A long time ago, in a galaxy far, far away....
}

func openTempFile(name, content string) (f *os.File, cleanup func()) {
	err := ioutil.WriteFile(name, []byte(content), 0777)
	if err != nil {
		panic(err)
	}
	f, err = os.Open(name)
	if err != nil {
		panic(err)
	}
	return f, func() { f.Close(); os.Remove(name) }
}

func writeTempFile(name, content string) (cleanup func()) {
	err := ioutil.WriteFile(name, []byte(content), 0777)
	if err != nil {
		panic(err)
	}
	return func() { os.Remove(name) }
}
