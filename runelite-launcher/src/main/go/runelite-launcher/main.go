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

var (
	logger func(format string, a ...interface{})
	progress func(value int)
	closeWindow func()
)

func main() {
	run := func(path string) {
		logger("Launching %s", path)
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


		// Parse bootstrap properties
		bootstrapPath := "http://static.runelite.net/bootstrap.json"
		logger("Downloading %s from %s", path.Base(bootstrapPath), bootstrapPath)
		bootstrap := ReadBootstrap(bootstrapPath)
		clientArtifactName := bootstrap.Client.ArtifactId
		clientArtifactVersion := bootstrap.Client.Version
		clientArtifactGroupId := bootstrap.Client.GroupId
		clientJarName := fmt.Sprintf("%s-%s-shaded.jar", clientArtifactName, clientArtifactVersion)

		// Get latest distribution from tag from git
		distributionArtifactName := "runelite-distribution"
		distributionArtifactVersion := strings.Replace(GetLatestTag("deathbeam/runelite-launcher").Name, "v", "", 1)
		distributionArtifactGroupId := "/*$mvn.project.groupId$*/"
		distributionDirName := fmt.Sprintf("%s-%s",
			distributionArtifactName,
			distributionArtifactVersion)
		distributionJarName := fmt.Sprintf("%s-%s.jar",
			distributionArtifactName,
			distributionArtifactVersion)
		distributionArchiveName := fmt.Sprintf("%s-%s-archive-distribution-%s.tar.gz",
			distributionArtifactName,
			distributionArtifactVersion,
			systemName)

		// Setup cache directories
		home, _ := homedir.Dir()
		runeliteHome := path.Join(home, ".runelite")
		launcherCache := path.Join(runeliteHome, "cache")
		systemCache := path.Join(launcherCache, systemName)
		distributionCache := path.Join(systemCache, distributionDirName)
		logger("Found system cache directory at %s", systemCache)

		if !FileExists(systemCache) {
			os.MkdirAll(systemCache, os.ModePerm)
		}

		// Setup versions
		distributionCacheVersionPath := path.Join(launcherCache, ".version-distribution")
		distributionCacheVersion := ReadVersion(distributionCacheVersionPath)
		clientCacheVersionPath := path.Join(launcherCache, ".version-client")
		clientCacheVersion := ReadVersion(clientCacheVersionPath)

		// Try to download distribution if not already downloaded
		distributionArchiveDestination := path.Join(launcherCache, distributionArchiveName)

		if !FileExists(distributionArchiveDestination) || !CompareVersion(distributionCacheVersion, distributionArtifactVersion) {
			baseUrl := "https://github.com/deathbeam/runelite-launcher/raw/mvn-repo"
			groupPath := strings.Replace(distributionArtifactGroupId, ".", "/", -1)
			archiveUrl := fmt.Sprintf("%s/%s/%s/%s/%s",
				baseUrl, groupPath, distributionArtifactName, distributionArtifactVersion, distributionArchiveName)

			os.RemoveAll(archiveUrl)

			DownloadFile(archiveUrl, distributionArchiveDestination, func(percent float64) {
				progress(int(percent))
			})
		} else {
			logger("Found distribution archive at %s", distributionArchiveDestination)
		}

		// Try to extract distribution if not already extracted
		if !FileExists(systemCache) || !CompareVersion(distributionCacheVersion, distributionArtifactVersion) {
			os.RemoveAll(systemCache)
			os.MkdirAll(systemCache, os.ModePerm)
			ExtractFile(distributionArchiveDestination, systemCache)
			SaveVersion(distributionCacheVersionPath, distributionArtifactVersion)
		} else {
			logger("Found distribution at %s", systemCache)
		}

		// Try to download shaded jar if not already present
		distributionPath := distributionCache

		if systemName == "darwin" {
			distributionPath = path.Join(distributionPath, "Contents", "Resources")
		}

		distributionJarDestination := path.Join(distributionPath, distributionJarName)

		if !FileExists(distributionJarDestination) || !CompareVersion(clientCacheVersion, clientArtifactVersion) {
			baseUrl := "http://repo.runelite.net/"
			groupPath := strings.Replace(clientArtifactGroupId, ".", "/", -1)
			shadedJarUrl := fmt.Sprintf("%s/%s/%s/%s/%s",
				baseUrl, groupPath, clientArtifactName, clientArtifactVersion, clientJarName)

			os.RemoveAll(distributionJarDestination)

			DownloadFile(shadedJarUrl, distributionJarDestination, func(percent float64) {
				progress(int(percent))
			})

			SaveVersion(clientCacheVersionPath, clientArtifactVersion)
		} else {
			logger("Found distribution jar at %s", distributionJarDestination)
		}

		// Build path to application executable
		distributionNativePath := distributionCache

		if systemName == "darwin" {
			distributionNativePath = path.Join(distributionNativePath, "Contents", "MacOS", distributionArtifactName)
		} else if strings.Contains(systemName, "windows") {
			distributionNativePath = path.Join(distributionNativePath, distributionArtifactName+".exe")
		} else {
			distributionNativePath = path.Join(distributionNativePath, distributionArtifactName)
		}

		go run(distributionNativePath)
		closeWindow()
	}

	CreateUI(boot)
}