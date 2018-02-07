package main

import (
	"io/ioutil"
	"log"
	"os"
)

func CompareVersion(old string, new string) bool {
	return old == new
}

func ReadVersion(file string) string {
	b, err := ioutil.ReadFile(file)
	log.Printf("Reading version from %s", file)

	if err != nil {
		return ""
	}

	return string(b)
}

func SaveVersion(file string, version string) {
	data := []byte(version)
	err := ioutil.WriteFile(file, data, os.ModePerm)

	if err != nil {
		panic(err)
	}

	log.Printf("Saving new new version %s to %s", version, file)
}
