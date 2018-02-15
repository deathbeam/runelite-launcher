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
	"encoding/json"
	"encoding/xml"
)

type Client struct {
	ArtifactId string `json:"artifactId"`
	GroupId    string `json:"groupId"`
	Version    string `json:"version"`
	Extension  string `json:"extension"`
}

type Bootstrap struct {
	Client Client `json:"client"`
}

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
