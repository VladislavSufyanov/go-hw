package hw02unpackstring

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

// CalcUnpackedLen determines the number of bytes in the unpacked string.
// If the string does not match the format, it returns (0, ErrInvalidString).
//
//nolint:gocognit,gocyclo,nestif
func CalcUnpackedLen(s string) (int, error) {
	firstRune, firstRuneSize := utf8.DecodeRuneInString(s)
	if firstRune == utf8.RuneError && firstRuneSize == 0 {
		return 0, nil
	}
	if firstRune == utf8.RuneError || (firstRune >= 48 && firstRune <= 57) {
		return 0, ErrInvalidString
	}
	totalSize := 0
	escapes := 0
	if firstRune == '\\' {
		escapes = 1
	}
	previousRune := firstRune
	for _, r := range s[firstRuneSize:] {
		if escapes == 0 {
			switch {
			case r == '\\':
				if previousRune < 48 || previousRune > 57 {
					totalSize += utf8.RuneLen(previousRune)
				}
				escapes = 1
			case r >= 48 && r <= 57 && previousRune >= 48 && previousRune <= 57:
				return 0, ErrInvalidString
			case r >= 48 && r <= 57:
				totalSize += int(r-'0') * utf8.RuneLen(previousRune)
			case previousRune < 48 || previousRune > 57:
				totalSize += utf8.RuneLen(previousRune)
			}
		} else {
			if escapes == 1 {
				if r != '\\' && (r < 48 || r > 57) {
					return 0, ErrInvalidString
				}
				escapes = 2
			} else {
				switch {
				case r == '\\':
					totalSize += utf8.RuneLen(previousRune)
					escapes = 1
				case r >= 48 && r <= 57:
					totalSize += int(r-'0') * utf8.RuneLen(previousRune)
					escapes = 0
				case r < 48 || r > 57:
					totalSize += utf8.RuneLen(previousRune)
					escapes = 0
				default:
					panic(fmt.Sprintf("unexpected rune %q", r))
				}
			}
		}
		previousRune = r
	}
	switch {
	case escapes == 1:
		return 0, ErrInvalidString
	case escapes == 2 || (escapes == 0 && (previousRune < 48 || previousRune > 57)):
		totalSize += utf8.RuneLen(previousRune)
	}
	return totalSize, nil
}

// Unpack performs string unpacking according to the format.
// It supports character escaping.
// "a4bc2d5e" => "aaaabccddddde"
// "abcd" => "abcd"
// "qwe\45"" => "qwe44444"
//
//nolint:gocognit
func Unpack(s string) (string, error) {
	totalSize, err := CalcUnpackedLen(s)
	if err != nil {
		return "", fmt.Errorf("error calc unpacked string len: %w", err)
	}
	if totalSize == 0 {
		return "", nil
	}
	firstRune, firstRuneSize := utf8.DecodeRuneInString(s)
	var builder strings.Builder
	builder.Grow(totalSize)
	escapes := 0
	if firstRune == '\\' {
		escapes = 1
	}
	previousRune := firstRune
	for _, r := range s[firstRuneSize:] {
		if escapes == 0 {
			switch {
			case r == '\\':
				if previousRune < 48 || previousRune > 57 {
					builder.WriteRune(previousRune)
				}
				escapes = 1
			case r >= 48 && r <= 57:
				for range int(r - '0') {
					builder.WriteRune(previousRune)
				}
			case previousRune < 48 || previousRune > 57:
				builder.WriteRune(previousRune)
			}
		} else {
			if escapes == 1 {
				escapes = 2
			} else {
				switch {
				case r == '\\':
					builder.WriteRune(previousRune)
					escapes = 1
				case r >= 48 && r <= 57:
					for range int(r - '0') {
						builder.WriteRune(previousRune)
					}
					escapes = 0
				case r < 48 || r > 57:
					builder.WriteRune(previousRune)
					escapes = 0
				default:
					panic(fmt.Sprintf("unexpected rune %q", r))
				}
			}
		}
		previousRune = r
	}
	if escapes == 2 || (escapes == 0 && (previousRune < 48 || previousRune > 57)) {
		builder.WriteRune(previousRune)
	}
	return builder.String(), nil
}
