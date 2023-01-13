package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func command(name string, args ...string) {
	if *printCmd {
		fmt.Println(name, strings.Join(args, ` `))
	}

	b, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		_, _ = os.Stderr.Write(b)
		panic(err)
	}

	if *verbose {
		fmt.Println(string(b))
	}
}

func execute(lookPath bool, name string, args ...string) {
	if *printCmd {
		fmt.Println(name, strings.Join(args, ` `))
	}

	var err error
	if lookPath {
		name, err = exec.LookPath(name)
		if err != nil {
			panic(err)
		}
	}

	err = syscall.Exec(name, append([]string{name}, args...), os.Environ())
	if err != nil {
		panic(err)
	}
}
