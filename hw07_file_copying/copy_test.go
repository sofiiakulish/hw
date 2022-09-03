package main

import (
	"crypto/md5"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

func FileMD5(path string) string {
	h := md5.New()
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = io.Copy(h, f)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func TestCopy(t *testing.T) {
	err := Copy("testdata/file_does_not_exist.txt", "/tmp/copy_offset0_limit0.txt", 0, 0)
	require.Error(t, err)

	err = Copy("testdata/input.txt", "testdata/input.txt", 0, 0)
	require.Error(t, err)

	err = Copy("testdata/input.txt", "/tmp/offset_is_bigger_then_file_size.txt", 100000, 0)
	require.Error(t, err)

	Copy("testdata/input.txt", "/tmp/copy_offset0_limit0.txt", 0, 0)
	require.Equal(t, FileMD5("testdata/out_offset0_limit0.txt"), FileMD5("/tmp/copy_offset0_limit0.txt"))

	Copy("testdata/input.txt", "/tmp/copy_offset0_limit10.txt", 0, 10)
	require.Equal(t, FileMD5("testdata/out_offset0_limit10.txt"), FileMD5("/tmp/copy_offset0_limit10.txt"))

	Copy("testdata/input.txt", "/tmp/copy_offset0_limit1000.txt", 0, 1000)
	require.Equal(t, FileMD5("testdata/out_offset0_limit1000.txt"), FileMD5("/tmp/copy_offset0_limit1000.txt"))

	Copy("testdata/input.txt", "/tmp/copy_offset0_limit10000.txt", 0, 10000)
	require.Equal(t, FileMD5("testdata/out_offset0_limit10000.txt"), FileMD5("/tmp/copy_offset0_limit10000.txt"))

	Copy("testdata/input.txt", "/tmp/copy_offset100_limit1000.txt", 100, 1000)
	require.Equal(t, FileMD5("testdata/out_offset100_limit1000.txt"), FileMD5("/tmp/copy_offset100_limit1000.txt"))

	Copy("testdata/input.txt", "/tmp/copy_offset6000_limit1000.txt", 6000, 1000)
	require.Equal(t, FileMD5("testdata/out_offset6000_limit1000.txt"), FileMD5("/tmp/copy_offset6000_limit1000.txt"))
}
