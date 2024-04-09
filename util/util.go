package util

import (
	"os"
	"path/filepath"
	"regexp"
)

func GetAbsFilePath(fp string) string {
	absPath, err := filepath.Abs(fp)
	if err != nil {
		panic(err)
	}
	return absPath
}

func OpenFile(fp string) *os.File {
	r, err := os.Open(GetAbsFilePath(fp))
	if err != nil {
		panic(err)
	}
	return r
}

func ReadFile(fp string) string {
	r, err := os.ReadFile(GetAbsFilePath(fp))
	if err != nil {
		panic(err)
	}
	return string(r)
}

func GetOpenAPIToken(path string) (string, error) {
	token := ReadFile(path)

	reg := regexp.MustCompile(`[^\w+]`)

	return reg.ReplaceAllString(string(token), ""), nil
}
