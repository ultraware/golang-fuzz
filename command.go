package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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
