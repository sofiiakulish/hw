package main

import (
	"errors"
	"fmt"
	"github.com/cheggaaa/pb"
	"io"
	"os"
)

var (
	// ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrTheSameFile           = errors.New("you are trying to rewrite the source file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == toPath {
		return ErrTheSameFile
	}

	inFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("open input file: %w", err)
	}

	defer inFile.Close()

	fileStat, err := os.Stat(fromPath)
	if err != nil {
		return fmt.Errorf("read input file info: %w", err)
	}

	fileSize := fileStat.Size()
	if fileSize < offset {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 {
		limit = fileSize
	}

	if limit > fileSize-offset {
		limit = fileSize - offset
	}

	outFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("open output file: %w", err)
	}

	defer outFile.Close()

	if offset > 0 {
		inFile.Seek(offset, io.SeekStart)
	}

	bar := pb.StartNew(int(limit))
	bar.Start()

	reader := bar.NewProxyReader(inFile)

	_, err = io.CopyN(outFile, reader, limit)
	if err != nil {
		return fmt.Errorf("copy data to file: %w", err)
	}

	bar.Finish()

	return nil
}
