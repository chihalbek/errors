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
	ErrInvalidURL       = errors.New("неверный url")
	ErrConnectionFailed = errors.New("ошибка подсключения")
	ErrDownloadFailed   = errors.New("ошибка. проверьте доступность файла для скачивания")
	ErrFileNotFound     = errors.New("файл не найден")
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

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return ErrDownloadFailed
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
			fmt.Println("неверный url")
		case errors.Is(err, ErrConnectionFailed):
			fmt.Println("ошибка подсключения")
		case errors.Is(err, ErrDownloadFailed):
			fmt.Println("ошибка. проверьте доступность файла для скачивания")
		case errors.Is(err, ErrFileNotFound):
			fmt.Println("файл не найден")
		}
		return
	}
}
