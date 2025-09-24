package ascii

import "strings"

// ASCII art patterns for digits 0-9, colon, and space
var asciiDigits = map[rune][]string{
	'0': {
		"██████",
		"██  ██",
		"██  ██",
		"██  ██",
		"██████",
	},
	'1': {
		"████  ",
		"  ██  ",
		"  ██  ",
		"  ██  ",
		"██████",
	},
	'2': {
		"██████",
		"    ██",
		"██████",
		"██    ",
		"██████",
	},
	'3': {
		"██████",
		"    ██",
		"██████",
		"    ██",
		"██████",
	},
	'4': {
		"██  ██",
		"██  ██",
		"██████",
		"    ██",
		"    ██",
	},
	'5': {
		"██████",
		"██    ",
		"██████",
		"    ██",
		"██████",
	},
	'6': {
		"██████",
		"██    ",
		"██████",
		"██  ██",
		"██████",
	},
	'7': {
		"██████",
		"    ██",
		"    ██",
		"    ██",
		"    ██",
	},
	'8': {
		"██████",
		"██  ██",
		"██████",
		"██  ██",
		"██████",
	},
	'9': {
		"██████",
		"██  ██",
		"██████",
		"    ██",
		"██████",
	},
	':': {
		"  ",
		"██",
		"  ",
		"██",
		"  ",
	},
	' ': {
		"                                  ",
		"                                  ",
		"                                  ",
		"                                  ",
		"                                  ",
	},
}

// ToASCII converts a time string like "25:00" to ASCII art
func ToASCII(timeStr string) string {
	if len(timeStr) == 0 {
		return ""
	}

	// Get ASCII patterns for each character
	var patterns [][]string
	for _, char := range timeStr {
		if pattern, exists := asciiDigits[char]; exists {
			patterns = append(patterns, pattern)
		}
	}

	if len(patterns) == 0 {
		return timeStr // fallback to original string
	}

	// Build the ASCII art line by line
	var result []string
	height := len(patterns[0])

	for line := 0; line < height; line++ {
		var lineStr strings.Builder
		for i, pattern := range patterns {
			if i > 0 {
				lineStr.WriteString("  ") // spacing between characters
			}
			if line < len(pattern) {
				lineStr.WriteString(pattern[line])
			}
		}
		result = append(result, lineStr.String())
	}

	return strings.Join(result, "\n")
}
