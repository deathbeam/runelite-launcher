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
		distributionCache := path.Join(launcherCache, "distribution")

		if !FileExists(launcherCache) {
			os.MkdirAll(launcherCache, os.ModePerm)
		}

		// Parse bootstrap properties
		bootstrap := ReadBootstrap("http://static.runelite.net/bootstrap.json")

		// Get latest repository tag
		latestTag := GetLatestTag("deathbeam/runelite-launcher")

		// Create distribution repository and distribution artifact
		distributionRepository := Repository{
			Url: "https://github.com/deathbeam/runelite-launcher/raw/mvn-repo",
			LocalPath: launcherCache,
		}

		distributionArtifact := Artifact{
			ArtifactId: "runelite-distribution",
			GroupId: "/*$mvn.project.groupId$*/",
			Version: strings.Replace(latestTag.Name, "v", "", 1),
			Suffix: fmt.Sprintf("-archive-distribution-%s.tar.gz", systemName),
		}

		// Create runelite repository and client artifact
		clientRepository := Repository{
			Url: "http://repo.runelite.net",
			LocalPath: launcherCache,
		}

		clientArtifact := Artifact{
			ArtifactId: bootstrap.Client.ArtifactId,
			GroupId: bootstrap.Client.GroupId,
			Version: bootstrap.Client.Version,
			Suffix: "-shaded.jar",
		}

		// Download, process, unarchive, copy distribution and client
		distributionPath := ProcessRemoteArchive(distributionArtifact, distributionRepository, distributionCache, systemName)
		ProcessRemoteExecutable(clientArtifact, clientRepository, distributionPath)

		// Build path to application executable
		distributionNativePath := distributionPath

		if systemName == "darwin" {
			distributionNativePath = path.Join(distributionNativePath, "Contents", "MacOS", distributionArtifact.ArtifactId)
		} else if strings.Contains(systemName, "windows") {
			distributionNativePath = path.Join(distributionNativePath, distributionArtifact.ArtifactId + ".exe")
		} else {
			distributionNativePath = path.Join(distributionNativePath, distributionArtifact.ArtifactId)
		}

		run(distributionNativePath)
	}

	CreateUI(boot)
}