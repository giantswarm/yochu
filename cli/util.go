package cli

import (
	"fmt"
	"os"
)

// ExitStderr exits the process with code 1 and prints output to Stderr, using the error message of err.
// If --debug is set to true, a helpful error will be printed.
func ExitStderr(err error) {
	debug(err)
	stderrf(err.Error())
	os.Exit(1)
}

// ExitStderrf exits the process with code 1 and print output to Stderr, using all objects of v in format f.
// If --debug is set to true, helpful errors will be printed if given.
func ExitStderrf(f string, v ...interface{}) {
	debug(v...)
	stderrf(f, v...)
	os.Exit(1)
}

// ExitStdoutf exits the process with code 0 and print output to Stdout, using all objects of v in format f.
func ExitStdoutf(f string, v ...interface{}) {
	stdoutf(f, v...)
	os.Exit(0)
}

// Verbosef prints output that is interesting for the user to Stdout, using all objects of v in format f.
func Verbosef(f string, v ...interface{}) {
	if globalFlags.verbose {
		stdoutf(f, v...)
	}
}

// Stdoutf prints output for the user to Stdout, using all objects of v in format f.
func Stdoutf(f string, v ...interface{}) {
	stdoutf(f, v...)
}

////////////////////////////////////////////////////////////////////////////////
// private

func debug(v ...interface{}) {
	if globalFlags.debug {
		for _, obj := range v {
			if _, isErr := obj.(error); isErr {
				fmt.Println()
				fmt.Printf("DEBUG: %#v", obj)
				fmt.Println()
				fmt.Println()
			}
		}
	}
}

func stderrf(f string, v ...interface{}) {
	printf(os.Stderr, f, v...)
}

func stdoutf(f string, v ...interface{}) {
	printf(os.Stdout, f, v...)
}

func printf(file *os.File, f string, v ...interface{}) {
	if f == "" {
		return
	}

	fmt.Fprintf(file, f, v...)
	fmt.Println()
}
