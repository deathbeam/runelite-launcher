package main

import "encoding/json"

type Client struct {
	ArtifactId string `json:"artifactId"`
	GroupId    string `json:"groupId"`
	Version    string `json:"version"`
}

type Bootstrap struct {
	Client Client `json:"client"`
}

func ReadBootstrap(url string) Bootstrap {
	file := FetchFile(url)

	var bootstrap Bootstrap

	if err := json.Unmarshal(file, &bootstrap); err != nil {
		panic(err)
	}

	return bootstrap
}
