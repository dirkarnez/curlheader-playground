package main

import (
  "os"
	"path/filepath"
)

func AuthorizatedGet(client *http.Client, url string, header *http.Header) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = *header

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		return nil, err
	}

	return body, nil
}

func main() {
  ex, _ := os.Executable()
	exPath := filepath.Dir(ex)
	seedFileName := filepath.Join(exPath, "seed.txt")
	// MustHaveFile(seedFileName)

	header, _ := curlheader.GetCurlHeader(seedFileName)
	client := &http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	var get = func(url string) ([]byte, error) {
		return AuthorizatedGet(client, url, &header)
	}

  var download = func(url, fileName string) error {
		bytes, err := get(url)
		if err != nil {
			return err
		}

		return os.WriteFile(fileName, bytes, 0644)
	}
}
