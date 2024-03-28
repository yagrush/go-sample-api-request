package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	formElementNameFile = "file"
	httpMethodPost      = "POST"
)

func getAbsFilePath(filePath string) string {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		panic(err)
	}
	return absPath
}

func openFile(filePath string) *os.File {
	r, err := os.Open(getAbsFilePath(filePath))
	if err != nil {
		panic(err)
	}
	return r
}

func readFile(filePath string) string {
	r, err := os.ReadFile(getAbsFilePath(filePath))
	if err != nil {
		panic(err)
	}
	return string(r)
}

func getOpenAPIToken(path string) (string, error) {
	token := readFile(path)

	reg := regexp.MustCompile(`[^\w+]`)

	return reg.ReplaceAllString(string(token), ""), nil
}

func main() {
	urlBase := flag.String("urlBase", "https://hogehoge.com/uploadFile?token=%s", "url base")
	pathOpenAPITokenFile := flag.String("token", "./config/token", "path of open api token file")
	filePathForUpload := flag.String("file", "./README.md", "path of upload file")
	flag.Parse()

	// url
	token, err := getOpenAPIToken(*pathOpenAPITokenFile)
	if err != nil {
		panic(err)
	}

	url := fmt.Sprintf(*urlBase, token)

	// request
	ret, err := requestFileUpload(url, *filePathForUpload)
	if err != nil {
		panic(err)
	}

	log.Println(ret)
	log.Println("--- Finished ---")
}

func requestFileUpload(url, filePath string) (string, error) {
	var err error

	client := &http.Client{}

	var body bytes.Buffer
	bodyWriter := multipart.NewWriter(&body)

	file := openFile(filePath)
	defer file.Close()

	// ファイルを添付する
	var tempWriter io.Writer
	if tempWriter, err = bodyWriter.CreateFormFile(formElementNameFile, file.Name()); err != nil {
		return "", err
	}
	if _, err = io.Copy(tempWriter, file); err != nil {
		return "", err
	}

	// その他要素を付加するとき
	if tempWriter, err = bodyWriter.CreateFormField("greeting"); err != nil {
		return "", err
	}
	if _, err = io.Copy(tempWriter, strings.NewReader("hello world!")); err != nil {
		return "", err
	}

	bodyWriter.Close()

	// request
	req, err := http.NewRequest(httpMethodPost, url, &body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad request: %s", res.Status)
		return "", err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	// 例
	// このようなJSONレスポンスBODYを期待している場合
	// {"data": {"fileId": "1234567"}}
	var response ThisResponse
	json.Unmarshal(resBody, &response)

	return response.Data.FileId, nil
}

type ThisResponse struct {
	Data ThisResponseData `json:"data"`
}

type ThisResponseData struct {
	FileId string `json:"fileId"`
}
