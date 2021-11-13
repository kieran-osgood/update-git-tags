package internal

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
)

//goland:noinspection GoSnakeCaseUsage
type AllFlags struct {
	Project_url   *string
	Branch        *string
	Ssh_Key       *string
	Previous_hash *string
	Property_path *string
}

var Flags = AllFlags{
	Project_url: flag.String("project_url", os.Getenv("REPOSITORY_URL"), "Git repository url"),
	Branch:      flag.String("branch", "main", "The default branch to checkout and tag"),
	Ssh_Key:     flag.String("ssh_key", os.Getenv("SSH_KEY"), "The ssh key file to use to clone"),
	Previous_hash:     flag.String("previous_hash", "", "The previous commit hash to compare file against"),
	Property_path:     flag.String("property_path", "", "Dot notation form to access the version field in JSON"),
}

func ParseFlags() error {
	flag.Parse()

	v := reflect.ValueOf(Flags)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		p := v.Field(i).Elem()
		err := CheckFlag(typeOfS.Field(i).Name, &p)

		HandleError(err)
	}

	return nil
}

func CheckFlag(name string, value *reflect.Value) error {
	if value.IsZero() {
		return errors.New(fmt.Sprintf("%v is a required field", name))
	}
	return nil
}
