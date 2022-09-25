package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
)

const INVALID_FILE_NAME_CHARACTER = "="

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	var valName string
	var valValue string
	var needRemove bool
	var ok bool

	result := make(Environment)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return result, fmt.Errorf("read dir: %w", err)
	}

	for _, f := range files {

		valName, ok = getFileName(f)
		if ok == false {
			continue
		}

		valValue, needRemove, err = getValue(dir, f)
		if err != nil {
			return result, fmt.Errorf("get value error: %w", err)
		}

		result[valName] = EnvValue{valValue, needRemove}
	}

	return result, nil
}

func getFileName(file fs.FileInfo) (string, bool) {
	if file.IsDir() {
		return "", false
	}

	if strings.Contains(file.Name(), INVALID_FILE_NAME_CHARACTER) {
		return "", false
	}

	return file.Name(), true
}

func getValue(dir string, fileInfo fs.FileInfo) (string, bool, error) {
	if fileInfo.Size() == 0 {
		return "", true, nil
	}

	file, err := os.Open(dir + "/" + fileInfo.Name())

	if err != nil {
		return "", false, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	value := scanner.Text()

	err = scanner.Err()
	if err != nil {
		return "", false, err
	}

	result := strings.TrimRight(value, " \t")

	result = strings.Replace(result, "\x00", "\n", -1)

	return result, false, nil
}
