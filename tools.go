package webox

import (
	"net/http"
	"crypto/rand"
)

const randomStringSource = "abcdefghiklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_+"

type Tools struct {
	MaxFileSize int
}

func (t *Tools) RandomString(n int) string {
	generatedString, sourceCharacters := make([]rune, n), []rune(randomStringSource)
	for i := range generatedString {
		primeNum, _ := rand.Prime(rand.Reader, len(sourceCharacters))
		index1, index2 := primeNum.Uint64(), uint64(len(sourceCharacters))
		generatedString[i] = sourceCharacters[index1%index2]
	}
	
	return string(generatedString)
}

type UploadedFile struct {
	NewFileName string
	OriginalFileName string
	FileSize int64
}

func (t *Tools) UploadFiles(req *http.Request, uploadDir string, rename ...bool) ([]*UploadedFile, error) {
	renameFile := true
	if len(rename) > 0 {
		renameFile = rename[0]
	}
}