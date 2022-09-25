package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestReadDir(t *testing.T) {
	_, err := ReadDir("not_exist_dir")
	require.Error(t, err)

	result, err := ReadDir("testdata/env")
	require.NoError(t, err)
	require.Equal(t, getExpectedResult(), result)

	// check if ignores files
	testDirName := "testdata/env/2"
	createTestDir(testDirName)
	result, err = ReadDir("testdata/env")
	require.NoError(t, err)
	require.Equal(t, getExpectedResult(), result)
	removeTestDir(testDirName)

	// check if ignores files with the - in the name
	testFileName := "testdata/env/2=2.txt"
	createTestFile(testFileName, []byte("hello\ngo\n"))
	result, err = ReadDir("testdata/env")
	require.NoError(t, err)
	require.Equal(t, getExpectedResult(), result)
	removeTestFile(testFileName)
}

func createTestDir(dir string) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("Error creating test directory: %w", err))
	}
}

func removeTestDir(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		panic(fmt.Errorf("Error removing test directory: %w", err))
	}
}

func createTestFile(file string, data []byte) {
	err := os.WriteFile(file, data, os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("Error creating test file: %w", err))
	}
}

func removeTestFile(file string) {
	err := os.Remove(file)
	if err != nil {
		panic(fmt.Errorf("Error removing test directory: %w", err))
	}
}

func getExpectedResult() Environment {
	result := make(Environment)
	result["BAR"] = EnvValue{"bar", false}
	result["EMPTY"] = EnvValue{"", false}
	result["FOO"] = EnvValue{"   foo\nwith new line", false}
	result["HELLO"] = EnvValue{"\"hello\"", false}
	result["UNSET"] = EnvValue{"", true}

	return result
}
