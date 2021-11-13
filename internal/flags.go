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
	flags.StringVar(&allFlags.SshKey, "sshKey", os.Getenv("SSH_KEY"), "Base64 encoded of SSH private key. Default: $SSH_KEY")
	flags.StringVar(&allFlags.PreviousHash, "PreviousHash", os.Getenv("CIRCLE_SHA1"), "Commit hash of the previous commit to HEAD. Default: $CIRCLE_SHA1")
	flags.StringVar(&allFlags.PropertyPath, "PropertyPath", "Version", "Property path to the Version code in the json file. Default: \"Version\"")
	flags.StringVar(&allFlags.FilePath, "FilePath", "package.json", "File path to the json file with the Version code. Default: \"package.json\"")

	flags.BoolVar(&allFlags.Version, "Version", false, "Check Version of app binary")

	err = flags.Parse(args)
	if err != nil {
		return nil, buf.String(), err
	}

	skipFlagChecks := false
	if allFlags.Version  {
		// User checking Version || help - so we expect other flags not to be set
		skipFlagChecks = true
	}

	if !skipFlagChecks {
		err = CheckFlag("RepositoryUrl", reflect.ValueOf(allFlags.RepositoryUrl))
		if err != nil {
			return nil, buf.String(), err
		}
		err = CheckFlag("Branch", reflect.ValueOf(allFlags.Branch))
		if err != nil {
			return nil, buf.String(), err
		}
		err = CheckFlag("sshKey", reflect.ValueOf(allFlags.SshKey))
		if err != nil {
			return nil, buf.String(), err
		}
		err = CheckFlag("PreviousHash", reflect.ValueOf(allFlags.PreviousHash))
		if err != nil {
			return nil, buf.String(), err
		}
		err = CheckFlag("PropertyPath", reflect.ValueOf(allFlags.PropertyPath))
		if err != nil {
			return nil, buf.String(), err
		}
		err = CheckFlag("FilePath", reflect.ValueOf(allFlags.PropertyPath))
		if err != nil {
			return nil, buf.String(), err
		}
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
