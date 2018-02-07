package main

import (
	"bytes"
	"fmt"
	"github.com/verybluebot/tarinator-go"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func FetchFile(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return buf.Bytes()
}

func FileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func printDownloadPercent(done chan int64, path string, total int64) {
	stop := false

	for {
		select {
		case <-done:
			stop = true
		default:

			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}

			fi, err := file.Stat()
			if err != nil {
				log.Fatal(err)
			}

			size := fi.Size()

			if size == 0 {
				size = 1
			}

			var percent = float64(size) / float64(total) * 100

			fmt.Printf("%.0f", percent)
			fmt.Println("%")
		}

		if stop {
			break
		}

		time.Sleep(time.Second)
	}
}

func DownloadFile(url string, dest string) {
	log.Printf("Downloading %s to %s\n", url, dest)

	start := time.Now()

	out, err := os.Create(dest)

	if err != nil {
		fmt.Println(dest)
		panic(err)
	}

	defer out.Close()

	headResp, err := http.Head(url)

	if err != nil {
		panic(err)
	}

	defer headResp.Body.Close()

	size, err := strconv.Atoi(headResp.Header.Get("Content-Length"))

	if err != nil {
		panic(err)
	}

	done := make(chan int64)

	go printDownloadPercent(done, dest, int64(size))

	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)

	if err != nil {
		panic(err)
	}

	done <- n

	elapsed := time.Since(start)
	log.Printf("Download completed in %s", elapsed)
}

func ExtractFile(file string, dest string) {
	start := time.Now()
	log.Printf("Extracting file %s to %s\n", file, dest)

	err := tarinator.UnTarinate(dest, file)

	if err != nil {
		panic(err)
	}

	elapsed := time.Since(start)
	log.Printf("Extracting completed in %s", elapsed)
}
