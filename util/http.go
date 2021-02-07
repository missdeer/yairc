package util

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/gzip"

	"github.com/andybalholm/brotli"
)

func uncompressReader(r *http.Response) (io.ReadCloser, bool, error) {
	header := strings.ToLower(r.Header.Get("Content-Encoding"))
	switch header {
	case "":
		return r.Body, false, nil
	case "br":
		rc := brotli.NewReader(r.Body)
		if rc == nil {
			log.Println("creating brotli reader failed")
			return nil, false, errors.New("creating brotli reader failed")
		}
		return ioutil.NopCloser(rc), true, nil
	case "gzip":
		rc, err := gzip.NewReader(r.Body)
		if err != nil {
			log.Println("creating gzip reader failed:", err)
			return nil, false, err
		}
		return rc, true, nil
	case "deflate":
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("reading inflate failed:", err)
			return nil, false, err
		}
		rc := flate.NewReader(bytes.NewReader(content[2:]))
		if rc == nil {
			log.Println("creating deflate reader failed")
			return nil, false, errors.New("creating deflate reader failed")
		}
		return rc, true, nil
	}
	return nil, false, errors.New("unexpected encoding type")
}

func OpenURI(uri string) (rc io.ReadCloser, err error) {
	if strings.HasPrefix(uri, "https://") || strings.HasPrefix(uri, "http://") {
		req, err := http.NewRequest("GET", uri, nil)
		if err != nil {
			return os.Open(uri)
		}

		req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")

		httpClient := http.Client{}
		resp, err := httpClient.Do(req)
		if err != nil {
			return os.Open(uri)
		}

		if resp.StatusCode != 200 {
			return os.Open(uri)
		}

		rc, _, err = uncompressReader(resp)
		if err == nil {
			return rc, err
		}
	}

	return os.Open(uri)
}
