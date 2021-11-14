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
	Version bool

	RepositoryUrl string
	Branch        string
	SshKey        string
	PreviousHash  string
	PropertyPath  string
	FilePath      string

	VersionPrefix string
	VersionSuffix string

	// args are the positional (non-flag) command-line arguments.
	args []string
}

func GetFlags() (*AllFlags, error) {
	flags, output, err := ParseFlags(os.Args[0], os.Args[1:])
	if err == flag.ErrHelp {
		fmt.Println(output)
		os.Exit(2)
	} else if err != nil {
		fmt.Println("got error:", err)
		fmt.Println("output:\n", output)
		os.Exit(1)
	}

	return flags, nil
}

func ParseFlags(programName string, args []string) (config *AllFlags, output string, err error) {
	flags := flag.NewFlagSet(programName, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	var allFlags AllFlags
	flags.StringVar(&allFlags.RepositoryUrl, "RepositoryUrl", os.Getenv("REPOSITORY_URL"), "Repository URL to check. Default: $REPOSITORY_URL")
	flags.StringVar(&allFlags.Branch, "Branch", "main", "Branch to check. Default: \"main\"")
	flags.StringVar(&allFlags.SshKey, "SshKey", os.Getenv("SSH_KEY"), "Base64 encoded of SSH private key. Default: $SSH_KEY")
	flags.StringVar(&allFlags.PreviousHash, "PreviousHash", os.Getenv("CIRCLE_SHA1"), "Commit hash of the previous commit to HEAD. Default: $CIRCLE_SHA1")
	flags.StringVar(&allFlags.PropertyPath, "PropertyPath", "Version", "Property path to the Version code in the json file. Default: \"Version\"")
	flags.StringVar(&allFlags.FilePath, "FilePath", "package.json", "File path to the json file with the Version code. Default: \"package.json\"")

	flags.StringVar(&allFlags.VersionPrefix, "VersionPrefix", "v", "Prefix for the git tag")
	flags.StringVar(&allFlags.VersionSuffix, "VersionSuffix", "", "Suffix for the git tag")

	// Flag short circuits' app run and outputs the binary version from main.Version
	flags.BoolVar(&allFlags.Version, "Version", false, "Check Version of app binary")

	err = flags.Parse(args)
	if err != nil {
		return nil, buf.String(), err
	}

	skipFlagChecks := false
	if allFlags.Version {
		// User checking Version || help - so we expect other flags not to be set
		skipFlagChecks = true
	}

	if !skipFlagChecks {
		flagsToCheck := map[string]reflect.Value{
			"RepositoryUrl": reflect.ValueOf(allFlags.RepositoryUrl),
			"Branch": reflect.ValueOf(allFlags.Branch),
			"SshKey": reflect.ValueOf(allFlags.SshKey),
			"PreviousHash": reflect.ValueOf(allFlags.PreviousHash),
			"PropertyPath": reflect.ValueOf(allFlags.PropertyPath),
			"FilePath": reflect.ValueOf(allFlags.FilePath),
		}

		err = CheckFlags(flagsToCheck)
		if err != nil {
			return nil, buf.String(), err
		}
	}

	allFlags.args = flags.Args()
	return &allFlags, buf.String(), nil
}

func CheckFlags(flags map[string]reflect.Value) error {
	for k, v := range flags  {
		err := CheckFlag(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func CheckFlag(name string, value reflect.Value) error {
	if value.IsZero() {
		return errors.New(fmt.Sprintf("%v is a required field", name))
	}
	return nil
}
