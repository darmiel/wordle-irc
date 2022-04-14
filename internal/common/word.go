package common

import (
	"fmt"
	"strings"
	"unicode"
)

type Word interface {
	Word() string
	At(idx int) rune
	ContainsChar(char rune) bool
	ColorForCharAt(index int, char rune) Color
	Print(input Word) string
}

type word string

func (w word) Word() string {
	return string(w)
}

func (w word) At(idx int) rune {
	return rune(string(w)[idx])
}

func (w word) ContainsChar(char rune) bool {
	char = unicode.ToLower(char)
	for _, r := range w {
		if r == char {
			return true
		}
	}
	return false
}

func (w word) ColorForCharAt(index int, char rune) Color {
	// correct char at correct index
	if w.At(index) == char {
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
func (w word) Print(input Word) string {
	var bob strings.Builder
	for i := range w {
		if bob.Len() > 0 {
			bob.WriteRune(' ')
		}
		guessedChar := input.At(i)
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
func (w word) IsHeterogram() bool {
	u := make(map[rune]bool)
	for _, char := range w {
		if _, ok := u[char]; ok {
			return false
		}
		u[char] = true
	}
	return false
}

///

func WordOf(w string) Word {
	return word(strings.TrimSpace(strings.ToLower(w)))
}
