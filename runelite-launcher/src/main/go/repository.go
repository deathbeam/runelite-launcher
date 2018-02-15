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
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"strings"
)

type Artifact struct {
	ArtifactId string
	GroupId    string
	Version    string
	Suffix     string
}

type Repository struct {
	Url       string
	LocalPath string
}

type MavenSnapshot struct {
	TimeStamp   string `xml:"timestamp"`
	BuildNumber string `xml:"buildNumber"`
}

type MavenVersions struct {
	Version string `xml:"version"`
}

type MavenVersioning struct {
	Release     string        `xml:"release"`
	Versions    MavenVersions `xml:"versions"`
	Snapshot    MavenSnapshot `xml:"snapshot"`
	LastUpdated string        `xml:"lastUpdated"`
}

type MavenMetadata struct {
	ArtifactId string          `xml:"artifactId"`
	GroupId    string          `xml:"groupId"`
	Versioning MavenVersioning `xml:"versioning"`
}

type Client struct {
	ArtifactId string `json:"artifactId"`
	GroupId    string `json:"groupId"`
	Version    string `json:"version"`
	Extension  string `json:"extension"`
}

type Bootstrap struct {
	Client Client `json:"client"`
}

// Read and parse data about artifacts from the bootstrap.json file
func ReadBootstrap(url string) (Bootstrap, error) {
	var bootstrap Bootstrap

	file, err := FetchFile(url)

	if err != nil {
		return bootstrap, err
	}

	if err = json.Unmarshal(file, &bootstrap); err != nil {
		return bootstrap, err
	}

	return bootstrap, nil
}

// Read and parse maven metadata file
func ReadMavenMetadata(url string) (MavenMetadata, error) {
	var mavenMetadata MavenMetadata

	file, err := FetchFile(url)

	if err != nil {
		return mavenMetadata, err
	}

	if err = xml.Unmarshal(file, &mavenMetadata); err != nil {
		return mavenMetadata, err
	}

	return mavenMetadata, nil
}

// Read maven checksum and decode it
func ReadMavenCheckSum(url string) ([]byte, error) {
	checkSum, err := FetchFile(url)

	if err != nil {
		return []byte{}, err
	}

	decodedCheckSum := make([]byte, hex.DecodedLen(len(checkSum)))
	hex.Decode(decodedCheckSum, checkSum)
	return decodedCheckSum, nil
}

func DownloadArtifact(artifact Artifact, repository Repository) (string, error) {
	// Build effective path to artifact in repository
	groupPath := strings.Replace(artifact.GroupId, ".", "/", -1)
	artifactVersion := artifact.Version
	artifactRepoUrl := fmt.Sprintf("%s/%s/%s", repository.Url, groupPath, artifact.ArtifactId)
	artifactRepoVersionedUrl := fmt.Sprintf("%s/%s", artifactRepoUrl, artifactVersion)

	// If artifact is snapshot, we need to first get correct snapshot version
	if strings.Contains(artifactVersion, "SNAPSHOT") {
		mavenMetadata, err := ReadMavenMetadata(artifactRepoVersionedUrl + "/maven-metadata.xml")

		if err != nil {
			return "", err
		}

		snapshotMetadata := mavenMetadata.Versioning.Snapshot
		artifactVersion = strings.Replace(artifactVersion, "SNAPSHOT", "", 1) +
			snapshotMetadata.TimeStamp + "-" + snapshotMetadata.BuildNumber
	}

	// Build final artifact URL and destination path
	artifactName := fmt.Sprintf("%s-%s%s", artifact.ArtifactId, artifactVersion, artifact.Suffix)
	artifactUrl := fmt.Sprintf("%s/%s", artifactRepoVersionedUrl, artifactName)
	artifactDestination := path.Join(repository.LocalPath, artifactName)

	// Get maven checksum from repository
	checkSumUrl := fmt.Sprintf("%s.%s", artifactUrl, "md5")
	checkSum, err := ReadMavenCheckSum(checkSumUrl)

	if err != nil {
		return artifactDestination, err
	}

	if err := DownloadFile(artifactUrl, artifactDestination, checkSum, md5.New()); err != nil {
		return artifactDestination, err
	}

	return artifactDestination, nil
}

func ProcessArtifact(artifact Artifact, repository Repository, cache string) error {
	// Download artifact
	artifactPath, err := DownloadArtifact(artifact, repository)

	if err != nil {
		return err
	}

	if !CompareFiles(artifactPath, cache) {
		if strings.Contains(artifactPath, ".tar") {
			// Extract artifact if it is .tar*
			os.RemoveAll(cache)
			os.MkdirAll(cache, os.ModePerm)

			if err = ExtractFile(artifactPath, cache); err != nil {
				return err
			}
		} else {
			// Replace artifact with new one
			if err = CopyFile(artifactPath, cache); err != nil {
				return err
			}
		}
	}

	return nil
}
