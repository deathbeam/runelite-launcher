package main

import (
  "fmt"
  "os"
  "path"
  "strings"
)

type Artifact struct {
  ArtifactId string
  GroupId string
  Version string
  Suffix string
}

type Repository struct {
  Url string
  LocalPath string
}

func DownloadArtifact(artifact Artifact, repository Repository) (string, bool) {
  groupPath := strings.Replace(artifact.GroupId, ".", "/", -1)
  artifactName := fmt.Sprintf("%s-%s%s", artifact.ArtifactId, artifact.Version, artifact.Suffix)
  artifactUrl := fmt.Sprintf("%s/%s/%s/%s/%s",
    repository.Url, groupPath, artifact.ArtifactId, artifact.Version, artifactName)

  artifactDestination := path.Join(repository.LocalPath, artifactName)
  changed := false

  if !FileExists(artifactDestination) {
    changed = true
    DownloadFile(artifactUrl, artifactDestination)
  }

  return artifactDestination, changed
}

func ProcessArtifact(artifact Artifact, repository Repository, cache string) {
  // Download artifact
  artifactPath, artifactChanged := DownloadArtifact(artifact, repository)

  if artifactChanged || !CompareFiles(artifactPath, cache) {

    if strings.Contains(artifactPath, ".tar") {
      // Extract artifact if it is .tar*
      os.RemoveAll(cache)
      os.MkdirAll(cache, os.ModePerm)
      ExtractFile(artifactPath, cache)
    } else {
      // Replace artifact with new one
      CopyFile(artifactPath, cache)
    }
  }
}
