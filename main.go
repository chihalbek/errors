package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	ErrInvalidURL       = errors.New("invalid url")
	ErrConnectionFailed = errors.New("connection error")
	ErrDownloadFailed   = errors.New("file upload failed")
	ErrFileNotFound     = errors.New("file not found")
)

func downloadFile(inputUrl, filename string) error {
	_, err := url.Parse(inputUrl)
	if err != nil {
		return ErrInvalidURL
	}

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(inputUrl)
	if err != nil {
		return ErrConnectionFailed
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return ErrFileNotFound
		}
		return ErrDownloadFailed
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		errC := file.Close()
		if err == nil {
			err = errC
		}
	}(file)

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy file", ErrDownloadFailed)
	}
	return nil
}

func main() {
	inputUrl := flag.String("url", "", "URL для загрузки файла")
	output := flag.String("output", "", "Имя сохраненного файла")

	flag.Parse()

	err := downloadFile(*inputUrl, *output)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidURL):
			fmt.Println("invalid URL. please, enter correct address.")
		case errors.Is(err, ErrConnectionFailed):
			fmt.Println("connection error. try again or check your connection.")
		case errors.Is(err, ErrDownloadFailed):
			fmt.Println("file upload failed. check if the file is available for downloading.")
		case errors.Is(err, ErrFileNotFound):
			fmt.Println("file not found. check if the URL is correct or if the file is on the server.")
		}
		return
	}
}
