package webox

import (
	"net/http"
	"crypto/rand"
	"errors"
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

	var uploadedFiles []*UploadedFile

	if t.MaxFileSize == 0 {
		t.MaxFileSize = 1024 * 1024 * 1024
	}

	err := req.ParseMultipartForm(int64(t.MaxFileSize))

	if err != nil {
		return nil, errors.New("The uploaded file is too large")
	}

	// check if any files are stored in request
	for _, fHeaders := range req.MultipartForm.File {
		for _, headr := range fHeaders {
			uploadedFiles, err = func(uploadedFiles []*UploadedFile) ([]*UploadedFile, error) {
				var uploadedFile UploadedFile 
				infile, err := headr.Open()
				if err != nil {
					return nil, err
				}
				defer infile.Close()

				// check file type
				buff := make([]byte, 512)
				_, err = infile.Read(buff)
				if err != nil {
					return nil, err
				}

				// Todo : check to see if the file type is permitted
				allowed := false
				fileType := http.DetectContentType(buff)
				allowedTypes := []string("image/jpeg", "image/png", "image/gif")

				if len(allowedTypes) > 0 {
					for _, ftype := range allowedTypes {
						if strings.EqualFold(ftype, fileType) {
							allowed = true
							break
						}
					}
				}

			} (uploadedFiles)
		}
	}
}