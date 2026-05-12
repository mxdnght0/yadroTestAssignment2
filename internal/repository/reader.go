package repository

import (
	"bufio"
	"os"
	"strings"
)

type NameReader interface {
	ReadNames(path string) (<-chan string, <-chan error)
}

type fileNameReader struct{}

func NewFileNameReader() NameReader {
	return &fileNameReader{}
}

func (r *fileNameReader) ReadNames(path string) (<-chan string, <-chan error) {
	names := make(chan string)
	errs := make(chan error, 1)

	go func() {
		defer close(names)
		defer close(errs)

		f, err := os.Open(path)
		if err != nil {
			errs <- err
			return
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				names <- line
			}
		}

		if err := scanner.Err(); err != nil {
			errs <- err
		}
	}()

	return names, errs
}
