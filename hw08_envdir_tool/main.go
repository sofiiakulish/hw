package main

import (
	"os"
)

func main() {
	// go-envdir /path/to/env/dir command arg1 arg2
	args := os.Args
	env, err := ReadDir(args[1])
	if err != nil {
		panic(err)
	}

	RunCmd(args[2:], env)
}
