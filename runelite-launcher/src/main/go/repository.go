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

func ProcessRemoteArchive(artifact Artifact, repository Repository, cache string, systemName string) string {
  // Download artifact
  distributionArtifactPath, distributionArtifactChanged := DownloadArtifact(artifact, repository)

  // Setup path to extracted distribution
  distributionDirName := fmt.Sprintf("%s-%s", artifact.ArtifactId, artifact.Version)
  distributionPath := path.Join(cache, distributionDirName)

  if systemName == "darwin" {
    distributionPath = path.Join(distributionPath, "Contents", "Resources")
  }

  // Extract distribution if it is .tar*
  if strings.Contains(distributionArtifactPath, ".tar") &&
    distributionArtifactChanged ||
    !FileExists(distributionPath) {

    os.RemoveAll(cache)
    os.MkdirAll(cache, os.ModePerm)
    ExtractFile(distributionArtifactPath, cache)
  }

  return distributionPath
}

func ProcessRemoteExecutable(artifact Artifact, repository Repository, cache string) {
  // Download artifact
  clientArtifactPath, clientArtifactChanged := DownloadArtifact(artifact, repository)

  // Setup path to extracted executable
  distributionJarName := fmt.Sprintf("%s.jar", path.Base(cache))
  distributionJarDestination := path.Join(cache, distributionJarName)

  // Replace distribution executable with downloaded one
  if clientArtifactChanged || !CompareFiles(clientArtifactPath, distributionJarDestination) {
    CopyFile(clientArtifactPath, distributionJarDestination)
  }
}
