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
		return fmt.Errorf("invalid URL: %w: %w", ErrInvalidURL, err)
	}

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(inputUrl)
	if err != nil {
		return fmt.Errorf("connection error: &w: &w", ErrConnectionFailed, err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return ErrFileNotFound
		}
		return fmt.Errorf("file upload failed: %w: %w", ErrDownloadFailed, resp.StatusCode)
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
		return fmt.Errorf("failed to copy file: %w: %w", ErrDownloadFailed, err)
	}
	return nil
}

func main() {
	inputUrl := flag.String("url", "", "URL for downloading file")
	output := flag.String("output", "", "filename")

	flag.Parse()

	err := downloadFile(*inputUrl, *output)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidURL):
			fmt.Println("Invalid URL. Please, enter correct address.")
		case errors.Is(err, ErrConnectionFailed):
			fmt.Println("Connection error. Try again or check your connection.")
		case errors.Is(err, ErrDownloadFailed):
			fmt.Println("File upload failed. Check if the file is available for downloading.")
		case errors.Is(err, ErrFileNotFound):
			fmt.Println("File not found. Check if the URL is correct or if the file is on the server.")
		}
		return
	}
}
