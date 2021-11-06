package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"os"
	"reflect"
)

type AppJsonFile struct {
	Expo struct {
		Version string
	}
}

func main() {
	_, _, _, err := parseFlags()
	handleError(err)

	//fmt.Printf("ProjectUrlPtr: %v \n", *projectUrlPtr)
	//fmt.Printf("CHashPtr: %v", *cHashPtr)
	//fmt.Printf("PHashPtr: %v", *pHashPtr)

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

	currentVersionCode := result.Expo.Version
	oldVersionCode := "1.2.3"

	if currentVersionCode != oldVersionCode {
		fmt.Printf("TRUE")
	}

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

func handleError(err error) {
	if err != nil {
		panic(err)
	}
	return
}

//goland:noinspection GoSnakeCaseUsage
type AllFlags struct {
	Project_url *string
	Current_sha      *string
	Previous_sha      *string
}

var Flags = AllFlags{
	Project_url: flag.String("project_url", "", "Git repository url"),
	Current_sha:      flag.String("current_sha", "", "The current commit hash to check"),
	Previous_sha:      flag.String("previous_sha", "", "The previous commit hash to check"),
}

func parseFlags() (*string, *string, *string, error) {
	flag.Parse()

	v := reflect.ValueOf(Flags)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		p := v.Field(i).Elem()
		err := checkFlag(typeOfS.Field(i).Name, &p)

		if err != nil {
			return nil, nil, nil, err
		}
	}

	return Flags.Project_url, Flags.Current_sha, Flags.Previous_sha, nil
}

func checkFlag(name string, value *reflect.Value) error {
	if value.IsZero() {
		message := fmt.Sprintf("%v is a required field", name)
		return errors.New(message)
	}
	return nil
}
