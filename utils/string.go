package utils

import (
	"strings"
	"unicode"
)

func GenerateRepeatedLetter(char string, times int) string {
	return strings.Repeat(char, times)
}

// if IDWallet like that, it will get idwallet
// otherwise, it's fine
func CapitalToUnderScore(word string) string {
	var newWords []rune
	for k, v := range word {
		lv := unicode.ToLower(v)
		if k == 0 {
			newWords = append(newWords, lv)
			continue
		}
		if !unicode.IsUpper(rune(word[k-1])) && unicode.IsUpper(v) {
			newWords = append(newWords, '_')
			newWords = append(newWords, lv)
			continue
		}
		newWords = append(newWords, lv)

	}
	return string(newWords)
}

// func CapitalToUnderScore(word string) (newword string) {
// 	for k, v := range word {
// 		lv := strings.ToLower(string(v))
// 		if k == 0 {
// 			newword += lv
// 			continue
// 		}
// 		if !unicode.IsUpper(rune(word[k-1])) && unicode.IsUpper(v) {
// 			newword += "_" + lv
// 			continue
// 		}
// 		newword += lv
// 	}
// 	return
// }
