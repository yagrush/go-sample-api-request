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
	"strings"

	"github.com/yagrush/go-sample-api-request/util"
)

const (
	formElementNameFile = "file"
	httpMethodPost      = "POST"
)

func main() {
	urlBase := flag.String("urlBase", "https://hogehoge.com/uploadFile?token=%s", "url base")
	pathOpenAPITokenFile := flag.String("token", "./config/token", "path of open api token file")
	filePathForUpload := flag.String("file", "./README.md", "path of upload file")
	flag.Parse()

	// url
	token, err := util.GetOpenAPIToken(*pathOpenAPITokenFile)
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

	file := util.OpenFile(filePath)
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

	return response.Data.FileId, err
}

type ThisResponse struct {
	Data ThisResponseData `json:"data"`
}

type ThisResponseData struct {
	FileId string `json:"fileId"`
}
