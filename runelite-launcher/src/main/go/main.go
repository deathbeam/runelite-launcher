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
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

func main() {
	run := func(path string) {
		logger.LogLine("Launching %s", path)
		cmd := exec.Command(path)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
		os.Exit(0)
	}

	boot := func() {
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

		// Setup cache directories
		home, _ := homedir.Dir()
		runeliteHome := path.Join(home, ".runelite")
		launcherCache := path.Join(runeliteHome, "cache")
		distributionCache := path.Join(launcherCache, "runelite")

		if !FileExists(launcherCache) {
			os.MkdirAll(launcherCache, os.ModePerm)
		}

		// Create distribution repository and distribution artifact
		distributionRepository := Repository{
			Url: "https://github.com/deathbeam/runelite-launcher/raw/mvn-repo",
			LocalPath: launcherCache,
		}

		// Get latest repository metadata
		mavenMetadata := ReadMavenMetadata(fmt.Sprintf("%s/maven-metadata.xml", distributionRepository.Url))

		distributionArtifact := Artifact{
			ArtifactId: mavenMetadata.ArtifactId,
			GroupId: mavenMetadata.GroupId,
			Version: mavenMetadata.Versioning.Release,
			Suffix: fmt.Sprintf("-%s.tar.gz", systemName),
		}

		// Create runelite repository and client artifact
		clientRepository := Repository{
			Url: "http://repo.runelite.net",
			LocalPath: launcherCache,
		}

		// Parse bootstrap properties
		bootstrap := ReadBootstrap("http://static.runelite.net/bootstrap.json")

		clientArtifact := Artifact{
			ArtifactId: bootstrap.Client.ArtifactId,
			GroupId: bootstrap.Client.GroupId,
			Version: bootstrap.Client.Version,
			Suffix: "-shaded.jar",
		}

		// Download and unarchive distribution
		ProcessArtifact(distributionArtifact, distributionRepository, distributionCache)

		// Build path to application jar
		distributionJarPath := distributionCache

		if systemName == "darwin" {
			distributionJarPath = path.Join(distributionJarPath, "Contents", "Resources")
		}

		distributionJarPath = path.Join(
			distributionJarPath,
			fmt.Sprintf("%s-%s.jar", distributionArtifact.ArtifactId, distributionArtifact.Version))

		// Download and copy client
		ProcessArtifact(clientArtifact, clientRepository, distributionJarPath)

		// Build path to application executable
		distributionExecutablePath := distributionCache

		if systemName == "darwin" {
			distributionExecutablePath = path.Join(distributionExecutablePath, "Contents", "MacOS", distributionArtifact.ArtifactId)
		} else if strings.Contains(systemName, "windows") {
			distributionExecutablePath = path.Join(distributionExecutablePath, distributionArtifact.ArtifactId + ".exe")
		} else {
			distributionExecutablePath = path.Join(distributionExecutablePath, distributionArtifact.ArtifactId)
		}

		run(distributionExecutablePath)
	}

	CreateUI(boot)
}