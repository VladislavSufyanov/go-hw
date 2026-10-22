package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "🙃0", expected: ""},
		{input: "aaф0b", expected: "aab"},
		{input: "т1е2ст", expected: "теест"},
		{input: "Щ1", expected: "Щ"},
		{input: "a0b0c0d0v7", expected: "vvvvvvv"},
		{input: "🙃1😀3🙃1", expected: "🙃😀😀😀🙃"},
		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: `\32~1!2`, expected: `33~!!`},
		{input: `\\3abc`, expected: `\\\abc`},
		{input: `\5w2П1`, expected: `5wwП`},
		{input: `ф\70л4\22`, expected: `флллл22`},
		{input: `\\`, expected: `\`},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", `\`, `\G`, `abcd\`}
	for _, tc := range invalidStrings {
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func TestCalcUnpackedLen(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{input: "a4bc2d5e", expected: 13},
		{input: "abccd", expected: 5},
		{input: "", expected: 0},
		{input: "aaa0b", expected: 3},
		{input: "🙃0", expected: 0},
		{input: "aaф0b", expected: 3},
		{input: "т1е2ст", expected: 10},
		{input: "Щ1", expected: 2},
		{input: "a0b0c0d0v7", expected: 7},
		{input: "🙃1😀3🙃1", expected: 20},
		{input: "П9v9d0U9i1r9R9c9[9]9`9~9", expected: 100},
		{input: `qwe\4\5`, expected: 5},
		{input: `qwe\45`, expected: 8},
		{input: `qwe\\5`, expected: 8},
		{input: `qwe\\\3`, expected: 5},
		{input: `\32~1!2`, expected: 5},
		{input: `\\3abc`, expected: 6},
		{input: `\5w2П1`, expected: 5},
		{input: `ф\70л4\22`, expected: 12},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := CalcUnpackedLen(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestCalcUnpackedLenInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", `\`, `\G`, `abcd\`}
	for _, tc := range invalidStrings {
		t.Run(tc, func(t *testing.T) {
			_, err := CalcUnpackedLen(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
