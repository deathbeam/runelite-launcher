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
	"flag"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

func main() {
	// Setup CLI flags
	var flagVersion string
	var flagDebug bool

	flag.BoolVar(&flagDebug, "debug", false,
		"Enables debug logging on the RuneLite client")

	flag.StringVar(&flagVersion, "version", "",
		"Forces the launcher to download specific version of RuneLite client")

	flag.Parse()

	// Setup cache directories
	home, err := homedir.Dir()

	if err != nil {
		panic(err)
	}

	runeliteHome := path.Join(home, ".runelite")
	launcherCache := path.Join(runeliteHome, "cache")
	distributionCache := path.Join(launcherCache, "RuneLite")

	if !FileExists(launcherCache) {
		if err := os.MkdirAll(launcherCache, os.ModePerm); err != nil {
			panic(err)
		}
	}

	run := func(path string) error {
		var args []string

		if flagDebug {
			args = []string{"--debug"}
		}

		logger.LogLine("Launching %v...", path)
		cmd := exec.Command(path, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			return err
		}

		os.Exit(0)
		return nil
	}

	boot := func() error {
		// Build system name
		systemName := runtime.GOOS

		if systemName != "darwin" {
			switch runtime.GOARCH {
			case "386":
				systemName += "32"
			case "amd64":
				systemName += "64"
			}
		}

		// Parse distribution bootstrap properties
		distributionBootstrap, err := ReadBootstrap("/*$mvn.project.property.distribution.bootstrap.url$*/")

		if err != nil {
			return err
		}

		// Create distribution repository and distribution artifact
		distributionRepository := Repository{
			Url: "/*$mvn.project.property.distribution.repository.url$*/",
			LocalPath: launcherCache,
		}

		distributionArtifact := Artifact{
			ArtifactId: distributionBootstrap.Client.ArtifactId,
			GroupId: distributionBootstrap.Client.GroupId,
			Version: distributionBootstrap.Client.Version,
			Suffix: fmt.Sprintf("-%s.%s", systemName, distributionBootstrap.Client.Extension),
		}

		// Parse client bootstrap properties
		clientBootstrap, err := ReadBootstrap("http://static.runelite.net/bootstrap.json")

		if err != nil {
			return err
		}

		// Create runelite repository and client artifact
		clientRepository := Repository{
			Url: "http://repo.runelite.net",
			LocalPath: launcherCache,
		}

		clientArtifact := Artifact{
			ArtifactId: clientBootstrap.Client.ArtifactId,
			GroupId: clientBootstrap.Client.GroupId,
			Version: clientBootstrap.Client.Version,
			Suffix: fmt.Sprintf("-shaded.%s", clientBootstrap.Client.Extension),
		}

		// Force set the client version if set from CLI
		if flagVersion != "" {
			clientArtifact.Version = flagVersion
		}

		// Download and unarchive distribution
		if err = ProcessArtifact(distributionArtifact, distributionRepository, distributionCache); err != nil {
			return err
		}

		// Build path to application jar
		distributionJarPath := distributionCache

		if systemName == "darwin" {
			distributionJarPath = path.Join(distributionJarPath, "Contents", "Resources")
		}

		distributionJarPath = path.Join(
			distributionJarPath,
			fmt.Sprintf("%s-%s.jar", distributionArtifact.ArtifactId, distributionArtifact.Version))

		// Download and copy client
		if err = ProcessArtifact(clientArtifact, clientRepository, distributionJarPath); err != nil {
			return err
		}

		// Build path to application executable
		distributionExecutablePath := distributionCache

		if systemName == "darwin" {
			distributionExecutablePath = path.Join(distributionExecutablePath, "Contents", "MacOS", distributionArtifact.ArtifactId)
		} else if strings.Contains(systemName, "windows") {
			distributionExecutablePath = path.Join(distributionExecutablePath, distributionArtifact.ArtifactId + ".exe")
		} else {
			distributionExecutablePath = path.Join(distributionExecutablePath, distributionArtifact.ArtifactId)
		}

		if err = run(distributionExecutablePath); err != nil {
			return err
		}

		return nil
	}

	safeBoot := func() {
		const maxRetries = 3

		for i := 1; i <= maxRetries; i++ {
			err := boot()

			if err == nil {
				break
			}

			logger.LogLine("Unexpected error occurred: %s", err)

			if i == maxRetries {
				panic(err)
				os.Exit(1)
			}

			logger.LogLine("Cleaning cache and retrying (%d of %d)...", i, maxRetries)
			os.RemoveAll(launcherCache)
		}
	}

	CreateUI(safeBoot)
}
