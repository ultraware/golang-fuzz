package main

import (
	"os"
	"os/exec"
)

func command(name string, args ...string) {
	b, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		_, _ = os.Stderr.Write(b)
		panic(err)
	}
}
