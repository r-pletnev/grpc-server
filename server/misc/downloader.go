package misc

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Downloader struct {
	client *http.Client
}

func NewDownloader() *Downloader {
	return &Downloader{client: http.DefaultClient}
}

func (d Downloader) makeRequest(url string) (*bytes.Buffer, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %v", err)
	}

	response, err := d.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("get response: %v", err)
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("status code is %d, body: %v",response.StatusCode, response.Body )
	}
	defer response.Body.Close()
	var b bytes.Buffer
	if _, err := io.Copy(&b, response.Body); err != nil {
		return nil, err
	}

	return &b, nil
}

func (d Downloader) DownloadFile(url string) (*bytes.Buffer, error) {
	return d.makeRequest(url)
}
