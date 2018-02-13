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
	"github.com/cavaliercoder/grab"
	"github.com/mholt/archiver"
	"io"
	"net/http"
	"os"
	"time"
)

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

func FetchFile(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return []byte{}, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	if err = resp.Body.Close(); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func DownloadFile(url string, dest string) error {
	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest(dest, url)

	// start download
	logger.LogLine("Downloading %v...", req.URL())
	resp := client.Do(req)
	logger.LogLine("%v", resp.HTTPResponse.Status)

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			logger.UpdateProgress(int(100 * resp.Progress()))
		case <-resp.Done:
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		return err
	}

	logger.LogLine("Download saved to %v", resp.Filename)
	return nil
}

func ExtractFile(file string, dest string) error {
	logger.LogLine("Extracting %v...", file)

	if err := archiver.TarGz.Open(file, dest); err != nil {
		return err
	}

	logger.LogLine("Archive extracted to %v", dest)
	return nil
}

func CopyFile(file string, dest string) error {
	logger.LogLine("Copying %v...", file)

	in, err := os.Open(file)

	if err != nil {
		return err
	}

	out, err := os.Create(dest)

	if err != nil {
		return err
	}

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	if err = in.Close(); err != nil {
		return err
	}

	if err = out.Close(); err != nil {
		return err
	}

	logger.LogLine("File copied to %v", dest)
	return nil
}
