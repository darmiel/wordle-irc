package main

import (
	"fmt"
	"strings"
)

type Word string

func (w Word) Normalize() Word {
	return Word(strings.TrimSpace(strings.ToLower(string(w))))
}

func (w Word) ContainsChar(char rune) bool {
	for _, r := range w {
		if r == char {
			return true
		}
	}
	return false
}

func (w Word) ColorForCharAt(index int, char rune) Color {
	correctChar := rune(w[index])
	// correct char at correct index
	if correctChar == char {
		return ColorGreenBG
	}
	// correct char at incorrect index
	if w.ContainsChar(char) {
		return ColorYellowBG
	}
	// incorrect char
	return ColorGreyBG
}

// Print returns the word as a colored string
// NOTE: input is required to be normalized
func (w Word) Print(input Word) string {
	var bob strings.Builder
	for i := range w {
		if bob.Len() > 0 {
			bob.WriteRune(' ')
		}
		guessedChar := rune(input[i])
		// background color for char
		color := w.ColorForCharAt(i, guessedChar)
		bob.WriteString(color.String())
		// print char
		bob.WriteString(strings.ToUpper(fmt.Sprintf(" %c ", guessedChar)))
		// reset
		bob.WriteString(StyleReset.String())
	}
	return bob.String()
}

// IsHeterogram checks if the given word is a Heterogram
// NOTE: Input is required to be normalized
func (w Word) IsHeterogram() bool {
	u := make(map[rune]bool)
	for _, char := range w {
		if _, ok := u[char]; ok {
			return false
		}
		u[char] = true
	}
	return false
}
