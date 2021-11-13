package internal

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
)

type AllFlags struct {
	repositoryUrl string
	branch        string
	sshKey        string
	previousHash  string
	propertyPath  string

	// args are the positional (non-flag) command-line arguments.
	args []string
}

func ParseFlags(programName string, args []string) (config *AllFlags, output string, err error) {
	flags := flag.NewFlagSet(programName, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	var allFlags AllFlags
	flags.StringVar(&allFlags.repositoryUrl, "repositoryUrl", os.Getenv("REPOSITORY_URL"), "set repo project_url")
	flags.StringVar(&allFlags.branch, "branch", "main", "set repo project_url")
	flags.StringVar(&allFlags.sshKey, "sshKey", os.Getenv("SSH_KEY"), "set repo project_url")
	flags.StringVar(&allFlags.previousHash, "previousHash", os.Getenv("CIRCLE_SHA1"), "set repo project_url")
	flags.StringVar(&allFlags.propertyPath, "propertyPath", "version", "set repo project_url")

	err = flags.Parse(args)
	if err != nil {
		return nil, buf.String(), err
	}

	err = CheckFlag("repositoryUrl", reflect.ValueOf(allFlags.repositoryUrl))
	if err != nil {
		return nil, buf.String(), err
	}
	err = CheckFlag("branch", reflect.ValueOf(allFlags.branch))
	if err != nil {
		return nil, buf.String(), err
	}
	err = CheckFlag("sshKey", reflect.ValueOf(allFlags.sshKey))
	if err != nil {
		return nil, buf.String(), err
	}
	err = CheckFlag("previousHash", reflect.ValueOf(allFlags.previousHash))
	if err != nil {
		return nil, buf.String(), err
	}
	err = CheckFlag("propertyPath", reflect.ValueOf(allFlags.propertyPath))
	if err != nil {
		return nil, buf.String(), err
	}

	allFlags.args = flags.Args()
	return &allFlags, buf.String(), nil
}

func CheckFlag(name string, value reflect.Value) error {
	if value.IsZero() {
		return errors.New(fmt.Sprintf("%v is a required field", name))
	}
	return nil
}
