package main

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"os"
)

type AppJsonFile struct {
	Expo struct {
		Version string
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
	return
}

func main() {
	file, err := os.ReadFile("app.json")
	if err != nil {
		fmt.Println("file error")
		os.Exit(1)
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(file, &m)
	handleError(err)

	var result AppJsonFile
	err = mapstructure.Decode(m, &result)
	handleError(err)

	fmt.Printf("result: %v \n", result.Expo.Version)

	/*
	startup parameters:
		- git repo url
		- sha_1
		- sha_2

	variables needed:
		! need to parse json for this
		current commits [app.json].expo.version nd last commit

	// https://raw.githubusercontent.com/kieran-osgood/punchline/d32ad0be7dc4b187c751ce9378de2155dafa1a84/app.json?token=AEG6IAD5RGV4TNNSFMVRU5TBQVLES

	if last_sha != current_sha then update git commit
	*/
}
