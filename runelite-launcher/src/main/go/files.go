/*
 * Copyright (c) 2018, Tomas Slusny <slusnucky@gmail.com>
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
 * ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
package main

import (
	"bytes"
	"fmt"
	"github.com/verybluebot/tarinator-go"
	"io"
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

func CompareFiles(left string, right string) bool {
	if !FileExists(left) || !FileExists(right) {
		return false
	}

	leftFile, err := os.Stat(left)

	if err != nil {
		return false
	}

	rightFile, err := os.Stat(right)

	if err != nil {
		return false
	}

	return leftFile.IsDir() || rightFile.IsDir() || leftFile.Size() == rightFile.Size()
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
			fi, err := os.Stat(path)

			if err != nil {
				panic(err)
			}

			size := fi.Size()

			if size == 0 {
				size = 1
			}

			var percent = float64(size) / float64(total) * 100
			logger.UpdateProgress(int(percent))
		}

		if stop {
			break
		}

		time.Sleep(time.Second)
	}
}

func DownloadFile(url string, dest string) {
	logger.LogLine("Downloading %s to %s", url, dest)

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
	logger.LogLine("Download completed in %s", elapsed)
}

func ExtractFile(file string, dest string) {
	logger.LogLine("Extracting file %s to %s", file, dest)

	start := time.Now()
	err := tarinator.UnTarinate(dest, file)

	if err != nil {
		panic(err)
	}

	elapsed := time.Since(start)
	logger.LogLine("Extracting completed in %s", elapsed)
}

func CopyFile(src, dst string) {
	logger.LogLine("Copying file %s to %s", src, dst)
	start := time.Now()

	in, err := os.Open(src)

	if err != nil {
		panic(err)
	}

	defer in.Close()

	out, err := os.Create(dst)

	if err != nil {
		panic(err)
	}

	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		panic(err)
	}

	elapsed := time.Since(start)
	err = out.Close()

	if err != nil {
		panic(err)
	}

	logger.LogLine("Copying completed in %s", elapsed)
}
