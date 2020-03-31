package parser

import (
	"bufio"
	"io"
	"strings"

	"github.com/sirupsen/logrus"
)

type StreamScanner struct {
	scanner *bufio.Scanner
}

func NewStreamScanner(src io.Reader) *StreamScanner {
	return &StreamScanner{bufio.NewScanner(src)}
}

func (s *StreamScanner) Scan() *Scanner {
	if err := s.scanner.Err(); err != nil {
		panic(err)
	}
	var lines []string
	for s.scanner.Scan() {
		line := s.scanner.Text()
		if len(line) == 1 && line[0] == '\x1D' {
			break
		}
		lines = append(lines, line)
	}
	if err := s.scanner.Err(); err != nil {
		logrus.Fatal(err)
		return nil
	}

	if len(lines) > 0 {
		return NewScanner(strings.Join(lines, "\n"))
	}
	return nil
}
