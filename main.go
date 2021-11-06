package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/mitchellh/mapstructure"
	"os"
)

type AppJsonFile struct {
	Expo struct {
		Version string
	}
}

func main() {
	err := ParseFlags()
	HandleError(err)

	r, err := GetGitFiles()
	HandleError(err)

	file, err := os.ReadFile("app.json")
	if err != nil {
		fmt.Println("file error")
		os.Exit(1)
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(file, &m)
	HandleError(err)

	var result AppJsonFile
	err = mapstructure.Decode(m, &result)
	HandleError(err)

	currentVersionCode := result.Expo.Version
	oldVersionCode := "1.2.3"

	if currentVersionCode != oldVersionCode {
		_, err := r.CreateTag(currentVersionCode, plumbing.NewHash(*Flags.Current_sha), &git.CreateTagOptions{})
		HandleError(err)
	}
}
