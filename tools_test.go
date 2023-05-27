package webox

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"testing"
)

func TestTools_RandomString(t *testing.T) {
	var testTools Tools
	s := testTools.RandomString(10)
	if len(s) != 10 {
		t.Error("RandomString() failed - expected 10 characters")
	}
}

var uploadOneTests = []struct {
	name          string
	uploadDir     string
	errorExpected bool
}{
	{name: "valid", uploadDir: "./testdata/uploads/", errorExpected: false},
	{name: "invalid", uploadDir: "//", errorExpected: true},
}

func TestTools_UploadOneFile(t *testing.T) {
	for _, e := range uploadOneTests {

		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)

		go func() {
			defer writer.Close()

			part, err := writer.CreateFormFile("file", "./testdata/img.png")
			if err != nil {
				t.Error(err)
			}

			f, err := os.Open("./testdata/img.png")
			if err != nil {
				t.Error(err)
			}
			defer f.Close()
			img, _, err := image.Decode(f)
			if err != nil {
				t.Error("error decoding image", err)
			}

			err = png.Encode(part, img)
			if err != nil {
				t.Error(err)
			}
		}()

		request := httptest.NewRequest("POST", "/", pr)
		request.Header.Add("Content-Type", writer.FormDataContentType())

		var testTools Tools
		testTools.AllowedFileTypes = []string{"image/png"}

		uploadedFiles, err := testTools.UploadFile(request, e.uploadDir, true)
		if e.errorExpected && err == nil {
			t.Errorf("%s: error expected, but none received", e.name)
		}

		if !e.errorExpected {
			if _, err := os.Stat(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles.NewFileName)); os.IsNotExist(err) {
				t.Errorf("%s: expected file to exist: %s", e.name, err.Error())
			}

			_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles.NewFileName))
		}
	}
}

var uploadTests = []struct {
	name          string
	allowedTypes  []string
	renameFile    bool
	errorExpected bool
	maxSize       int
	uploadDir     string
}{
	{name: "allowed no rename", allowedTypes: []string{"image/jpeg", "image/png"}, renameFile: false, errorExpected: false, maxSize: 0, uploadDir: ""},
	{name: "allowed rename", allowedTypes: []string{"image/jpeg", "image/png"}, renameFile: true, errorExpected: false, maxSize: 0, uploadDir: ""},
	{name: "allowed no filetype specified", allowedTypes: []string{}, renameFile: true, errorExpected: false, maxSize: 0, uploadDir: ""},
	{name: "not allowed", allowedTypes: []string{"image/jpeg"}, errorExpected: true, maxSize: 0, uploadDir: ""},
	{name: "too big", allowedTypes: []string{"image/jpeg,", "image/png"}, errorExpected: true, maxSize: 10, uploadDir: ""},
	{name: "invalid directory", allowedTypes: []string{"image/jpeg,", "image/png"}, errorExpected: true, maxSize: 0, uploadDir: "//"},
}
