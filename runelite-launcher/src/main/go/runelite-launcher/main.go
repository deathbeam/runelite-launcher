package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

func main() {
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
	bootstrap := ReadBootstrap("http://static.runelite.net/bootstrap.json")
	clientArtifactName := bootstrap.Client.ArtifactId
	clientArtifactVersion := bootstrap.Client.Version
	clientJarName := fmt.Sprintf("%s-%s-shaded.jar", clientArtifactName, clientArtifactVersion)

	// TODO: Parse distribution properties from somewhere
	distributionArtifactName := "runelite-distribution"
	distributionArtifactVersion := "1.0.0-SNAPSHOT"
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
	log.Printf("System cache directory: %s", systemCache)

	// TODO: Try to download distribution if not already downloaded
	distributionArchiveDestination := path.Join(launcherCache, distributionArchiveName)

	// Try to extract distribution if not already extracted
	if !FileExists(systemCache) {
		os.MkdirAll(systemCache, os.ModePerm)
		ExtractFile(distributionArchiveDestination, systemCache)
	}

	// Try to download shaded jar if not already present
	distributionPath := distributionCache

	if systemName == "darwin" {
		distributionPath = path.Join(distributionPath, "Contents", "Resources")
	}

	distributionJarDestination := path.Join(distributionPath, distributionJarName)

	if !FileExists(distributionJarDestination) {
		baseUrl := "http://repo.runelite.net/"
		groupPath := strings.Replace(bootstrap.Client.GroupId, ".", "/", -1)
		shadedJarUrl := fmt.Sprintf("%s/%s/%s/%s/%s",
			baseUrl, groupPath, clientArtifactName, clientArtifactVersion, clientJarName)

		DownloadFile(shadedJarUrl, distributionJarDestination)
	}

	// Launch application
	distributionNativePath := distributionCache

	if systemName == "darwin" {
		distributionNativePath = path.Join(distributionNativePath, "Contents", "MacOS", distributionArtifactName)
	} else if strings.Contains(systemName, "windows") {
		distributionNativePath = path.Join(distributionNativePath, distributionArtifactName+".exe")
	} else {
		distributionNativePath = path.Join(distributionNativePath, distributionArtifactName)
	}

	log.Printf("Launching %s\n", distributionNativePath)
	cmd := exec.Command(distributionNativePath)
	cmd.Run()
}
