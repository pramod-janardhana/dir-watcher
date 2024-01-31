package scanner

import (
	"bufio"
	"io"
	"strings"
)

type Scanner struct {
	str string
}

func NewScanner(str string) *Scanner {
	return &Scanner{str: str}
}

// Scan scans the stream for magic string and returns the count
func (s *Scanner) Scan(reader io.Reader) int64 {
	sc := bufio.NewScanner(reader)
	sc.Split(bufio.ScanLines)

	var count int64 = 0
	for sc.Scan() {
		count = count + int64(strings.Count(sc.Text(), s.str))
	}

	return count
}
