package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...)
	setEnvironmentVariables(env)

	command.Stdout = os.Stdout
	command.Stdin = os.Stdin

	err := command.Start()
	if err != nil {
		panic(err)
	}

	err = command.Wait()
	if err != nil {
		panic(err)
	}

	return command.ProcessState.ExitCode()
}

func setEnvironmentVariables(env Environment) bool {
	for name, envValue := range env {
		if envValue.NeedRemove == true {
			os.Unsetenv(name)
			continue
		}

		os.Setenv(name, envValue.Value)
	}

	return true
}
