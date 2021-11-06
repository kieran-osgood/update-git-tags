package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
)

//goland:noinspection GoSnakeCaseUsage
type AllFlags struct {
	Project_url  *string
	Current_sha  *string
	Previous_sha *string
}

var Flags = AllFlags{
	Project_url:  flag.String("project_url", os.Getenv("FOO"), "Git repository url"),
	Current_sha:  flag.String("current_sha", os.Getenv("FOO"), "The current commit hash to check"),
	Previous_sha: flag.String("previous_sha", os.Getenv("FOO"), "The previous commit hash to check"),
}

func ParseFlags() error {
	flag.Parse()

	v := reflect.ValueOf(Flags)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		p := v.Field(i).Elem()
		err := CheckFlag(typeOfS.Field(i).Name, &p)

		if err != nil {
			return err
		}
	}

	return nil
}

func CheckFlag(name string, value *reflect.Value) error {
	if value.IsZero() {
		return errors.New(fmt.Sprintf("%v is a required field", name))
	}
	return nil
}
