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

func DownloadArtifact(artifact Artifact, repository Repository) (string, bool, error) {
  groupPath := strings.Replace(artifact.GroupId, ".", "/", -1)
  artifactName := fmt.Sprintf("%s-%s%s", artifact.ArtifactId, artifact.Version, artifact.Suffix)
  artifactUrl := fmt.Sprintf("%s/%s/%s/%s/%s",
    repository.Url, groupPath, artifact.ArtifactId, artifact.Version, artifactName)

  artifactDestination := path.Join(repository.LocalPath, artifactName)
  changed := false

  if !FileExists(artifactDestination) {
    changed = true

    if err := DownloadFile(artifactUrl, artifactDestination); err != nil {
      return artifactDestination, changed, err
    }
  }

  return artifactDestination, changed, nil
}

func ProcessArtifact(artifact Artifact, repository Repository, cache string) error {
  // Download artifact
  artifactPath, artifactChanged, err := DownloadArtifact(artifact, repository)

  if err != nil {
    return err
  }

  if artifactChanged || !CompareFiles(artifactPath, cache) {
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
