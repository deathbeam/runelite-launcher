package main

import (
	"bytes"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/verybluebot/tarinator-go"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"
)

func fileExists(name string) bool {
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

		if stop == true {
			break
		}

		time.Sleep(time.Second)
	}
}

func downloadFile(url string, dest string) {
	file := path.Base(url)

	log.Printf("Downloading file %s from %s\n", file, url)

	var buffer bytes.Buffer
	buffer.WriteString(dest)
	buffer.WriteString("/")
	buffer.WriteString(file)

	start := time.Now()

	out, err := os.Create(buffer.String())

	if err != nil {
		fmt.Println(buffer.String())
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

	go printDownloadPercent(done, buffer.String(), int64(size))

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

func extractFile(file string, dest string) {
	start := time.Now()
	log.Printf("Extracting file %s to %s\n", file, dest)

	err := tarinator.UnTarinate(dest, file)

	if err != nil {
		panic(err)
	}

	elapsed := time.Since(start)
	log.Printf("Extracting completed in %s", elapsed)
}

func main() {
	url := "https://sigterm.info/runelite-launcher.tar.gz"
	home, _ := homedir.Dir()
	runelite := path.Join(home, ".runelite")
	downloaded := path.Join(runelite, "runelite-launcher.tar.gz")
	cache := path.Join(runelite, "cache")
	executable := path.Join(cache, "out", "RuneLite.jar")

	if !fileExists(runelite) {
		os.MkdirAll(runelite, os.ModePerm)
	}

	if !fileExists(downloaded) {
		downloadFile(url, runelite)
	}

	if !fileExists(cache) {
		os.MkdirAll(cache, os.ModePerm)
		extractFile(downloaded, cache)
	}

	log.Printf("Launching %s\n", executable)
	cmnd := exec.Command("java", "-jar", executable)
	cmnd.Run()
}
