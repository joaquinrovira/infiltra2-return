package util

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// ToSlug transforms a string into a slug
func ToSlug(s string) string {
	// Remove diacritics
	s = removeDiacritics(s)

	// Convert to lower case
	s = strings.ToLower(s)

	// Replace spaces and underscores with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")

	// Remove all non-alphanumeric characters except for hyphens
	re := regexp.MustCompile("[^a-z0-9-]+")
	s = re.ReplaceAllString(s, "")

	// Remove leading and trailing hyphens
	s = strings.Trim(s, "-")

	// Remove duplicate hyphens
	s = removeDuplicateHyphens(s)

	return s
}

// RemoveDiacritics removes diacritics from a string
func removeDiacritics(s string) string {
	// Normalize the string to decompose combined characters into base characters and diacritics
	t := norm.NFD.String(s)
	// Use a strings.Builder to efficiently build the resulting string
	var b strings.Builder
	for _, r := range t {
		// If the rune is not a combining mark, write it to the builder
		if !unicode.Is(unicode.Mn, r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// removeDuplicateHyphens removes duplicate hyphens from a string
func removeDuplicateHyphens(s string) string {
	var builder strings.Builder
	var prev rune
	for _, current := range s {
		if current != '-' || prev != '-' {
			builder.WriteRune(current)
		}
		prev = current
	}
	return builder.String()
}