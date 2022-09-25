package main

import (
	"github.com/stretchr/testify/require"
	"testing"
	//"fmt"
	"os"
)

func TestCanSetEnvironmentVariables(t *testing.T) {
	// backup env variables before the test if any was set
	backup := backupEnvironment()

	os.Unsetenv("HELLO")

	env := make(Environment)
	env["HELLO"] = EnvValue{"hello", false}

	setEnvironmentVariables(env)

	value, present := os.LookupEnv("HELLO")
	require.Equal(t, true, present)
	require.Equal(t, "hello", value)

	// restore backup variables
	restoreEnvironmentFromTheBackup(backup)
}

func TestCanRemoveEnvironmentVariables(t *testing.T) {
	// backup env variables before the test if any was set
	backup := backupEnvironment()

	os.Setenv("HELLO", "hello")

	env := make(Environment)
	env["HELLO"] = EnvValue{"hello", true}

	setEnvironmentVariables(env)
	_, present := os.LookupEnv("HELLO")
	require.Equal(t, false, present)

	// restore backup variables
	restoreEnvironmentFromTheBackup(backup)
}

func TestRunCmd(t *testing.T) {
	cmd := make([]string, 2)
	cmd[0] = "testdata/echo.sh"
	cmd[1] = "arg1"
	code := RunCmd(cmd, getAllValues())
	require.Equal(t, 0, code)
}

func backupEnvironment() map[string]string {
	var testEnvNames = []string{"HELLO", "BAR", "FOO", "UNSET", "ADDED", "EMPTY"}
	result := make(map[string]string)

	for _, name := range testEnvNames {
		value, present := os.LookupEnv(name)
		if present == true {
			result[name] = value
		}
	}

	return result
}

func restoreEnvironmentFromTheBackup(env map[string]string) {
	for name, value := range env {
		os.Setenv(name, value)
	}
}

func getAllValues() Environment {
	result := make(Environment)
	result["HELLO"] = EnvValue{"hello", false}
	result["BAR"] = EnvValue{"bar", false}
	result["FOO"] = EnvValue{"foo", false}
	result["UNSET"] = EnvValue{"unset", false}
	result["ADDED"] = EnvValue{"added", false}
	result["EMPTY"] = EnvValue{"empty", false}

	return result
}
