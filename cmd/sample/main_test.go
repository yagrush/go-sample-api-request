package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testFileName = "main.go"

func startTestHttpServer() *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// dump, err := httputil.DumpRequest(r, true)
		// if err != nil {
		// 	fmt.Fprintln(w, err)
		// }
		// log.Println(string(dump))

		// body, err := io.ReadAll(r.Body)
		// if err != nil {
		// 	panic(err)
		// }
		// log.Println(string(body))

		// key "file" で添付送信されてくるファイルのファイル名をJSONに含めて返すだけの仮API
		_, fileHeader, _ := r.FormFile("file")
		fmt.Fprintf(w, "{\"data\": {\"fileId\": \"%s\"}}", fileHeader.Filename)
	})
	ts := httptest.NewServer(handler)

	return ts
}

func TestRequestFileUpload(t *testing.T) {
	ts := startTestHttpServer()
	defer ts.Close()

	ret, err := requestFileUpload(ts.URL, testFileName)

	if err != nil {
		t.Error(err)
	} else if ret != testFileName {
		t.Errorf("want:%s result:%s", testFileName, ret)
	}
}
