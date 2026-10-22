package hw03frequencyanalysis

import (
	"cmp"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

type wordFrequency struct {
	word  string
	count int
}

func Top10(s string) []string {
	if s == "" {
		return []string{}
	}
	words := strings.Fields(s)
	frequencyMap := map[string]int{}
	for _, w := range words {
		firstRune, size := utf8.DecodeRuneInString(w)
		if firstRune == utf8.RuneError {
			continue
		}
		if len(w) > size {
			trimed := strings.TrimFunc(w, func(r rune) bool {
				return unicode.IsPunct(r) || unicode.IsSymbol(r)
			})
			frequencyMap[strings.ToLower(trimed)]++
		} else if !unicode.IsSymbol(firstRune) && !unicode.IsPunct(firstRune) {
			frequencyMap[strings.ToLower(w)]++
		}
	}
	frequencySlice := make([]wordFrequency, 0, len(frequencyMap))
	for k, v := range frequencyMap {
		frequencySlice = append(frequencySlice, wordFrequency{k, v})
	}
	slices.SortFunc(frequencySlice, func(a wordFrequency, b wordFrequency) int {
		if n := cmp.Compare(b.count, a.count); n != 0 {
			return n
		}
		return cmp.Compare(a.word, b.word)
	})
	topLen := min(10, len(frequencySlice))
	topTen := make([]string, 0, topLen)
	for _, wf := range frequencySlice[:topLen] {
		topTen = append(topTen, wf.word)
	}
	return topTen
}
