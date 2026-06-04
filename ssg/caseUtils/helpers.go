package caseUtils

import (
	"strings"
	"unicode"
)

// ToSnakeCase converts a CamelCase string to snake_case.
func ToSnakeCase(str string) string {
	var b strings.Builder
	runes := []rune(str)
	strLen := len(runes)

	// Track if the last character written was an underscore
	// Initialize true to prevent leading underscores
	lastUnderscore := true

	for currentIndex := 0; currentIndex < strLen; currentIndex++ {
		currentRune := runes[currentIndex]
		nextIsLower := false
		if currentIndex+1 < strLen && unicode.IsLower(runes[currentIndex+1]) {
			nextIsLower = true
		}

		// 1. Handle delimiters (-, ., space, :)
		// Convert them all to underscores for processing
		if currentRune == '-' || currentRune == '.' || currentRune == ' ' ||
			currentRune == ':' || currentRune == '_' {
			if !lastUnderscore {
				b.WriteRune('_')
				lastUnderscore = true
			}
			continue
		}

		// 2. Handle Uppercase
		if currentRune >= 'A' && currentRune <= 'Z' {
			// Check if we need to insert an underscore before this capital
			if !lastUnderscore {
				prev := runes[currentIndex-1]
				isPrevLower := (prev >= 'a' && prev <= 'z') || (prev >= '0' && prev <= '9')

				// Insert _ if:
				// a. Previous was lowercase/digit (camelCase -> camel_case)
				// b. Previous was Upper but next is Lower (JSONId -> json_id)
				if isPrevLower || (currentRune >= 'A' && currentRune <= 'Z' && nextIsLower) {
					b.WriteRune('_')
					lastUnderscore = true
				}
			}

			// Convert to lowercase
			currentRune = currentRune + 32
		}

		// 3. Write the character
		b.WriteRune(currentRune)
		lastUnderscore = false
	}

	return b.String()
}

// SnakeToTitle converts a snake_case string to TitleCase.
func SnakeToTitle(str string) string {
	str = ToSnakeCase(str)
	allStrs := strings.Split(str, "_")
	builder := strings.Builder{}

	for _, current := range allStrs {
		builder.WriteString(ToTitle(current))
	}

	return builder.String()
}

// ToPascalCase converts a string to PascalCase.
func ToPascalCase(str string) string {
	title := SnakeToTitle(str)

	return strings.ToUpper(title[:1]) + title[1:]
}

// ToCamelCase converts a string to camelCase.
func ToCamelCase(s string) string {
	title := SnakeToTitle(s)

	return strings.ToLower(title[:1]) + title[1:]
}

// ToTitle function will convert the given string to title case.
func ToTitle(value string) string {
	return _titleCaser.String(value)
}
