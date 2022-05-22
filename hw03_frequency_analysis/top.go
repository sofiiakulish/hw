package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type Word struct {
	Value string
	Count int
}

func Top10(str string) []string {
	var words = strings.Fields(str)
	var groupedWords = GroupWords(words)
	sort.Slice(groupedWords, func(i, j int) bool {
		if groupedWords[i].Count == groupedWords[j].Count {
			return groupedWords[i].Value < groupedWords[j].Value
		}

		return groupedWords[i].Count > groupedWords[j].Count
	})

	var result = make([]string, 0)
	for k, v := range groupedWords {
		if k > 9 {
			return result
		}
		result = append(result, v.Value)
	}

	return result
}

func GroupWords(words []string) []Word {
	var result = make([]Word, 0)
	for _, word := range words {
		_, pos := ValueIsExists(word, result)
		if pos == -1 {
			result = append(result, Word{word, 0})
			pos = len(result) - 1
		}

		result[pos].Count = result[pos].Count + 1
	}

	return result
}

func ValueIsExists(word string, Words []Word) (Word, int) {
	for i, value := range Words {
		if value.Value == word {
			return value, i
		}
	}

	return Word{"", 0}, -1
}
