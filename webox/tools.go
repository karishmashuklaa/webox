package webox

import (
	"crypto/rand"
)

const randomStringSource = "abcdefghiklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_+"

type Tools struct {}

func (t *Tools) RandomString(n int) string {
	generatedString, sourceCharacters := make([]rune, n), []rune(randomStringSource)
	for i := range generatedString {
		primeNum, _ := rand.Prime(rand.Reader, len(sourceCharacters))
		index1, index2 := primeNum.Uint64(), uint64(len(sourceCharacters))
		generatedString[i] = sourceCharacters[index1%index2]
	}
	
	return string(generatedString)
}