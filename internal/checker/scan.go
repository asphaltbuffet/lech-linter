package checker

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

// Finding represents a single detected misspelling.
type Finding struct {
	File string
	Line int
	Col  int
	Word string
}

// Scan reads from r (labeled filename) and returns all misspelling findings.
func Scan(r io.Reader, filename string) ([]Finding, error) {
	var findings []Finding

	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		col := 1

		for word := range strings.FieldsSeq(line) {
			// Track column before stripping punctuation.
			wordCol := strings.Index(line[col-1:], word) + col
			col = wordCol + len(word)

			clean := strings.TrimFunc(word, func(r rune) bool {
				return !unicode.IsLetter(r) && !unicode.IsDigit(r)
			})

			if clean == "" {
				continue
			}

			if IsMisspelling(clean) {
				findings = append(findings, Finding{
					File: filename,
					Line: lineNum,
					Col:  wordCol,
					Word: clean,
				})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return findings, nil
}
