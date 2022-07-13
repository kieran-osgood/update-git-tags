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
	Version          bool
	RepositoryUrl    string
	Branch           string
	SshKey           string
	SshPhrase        string
	PreviousHash     string
	PropertyPath     string
	FilePath         string
	VersionTagPrefix string
	VersionTagSuffix string
	// args are the positional (non-flag) command-line arguments.
	args []string
}

func GetFlags() (*AllFlags, error) {
	flags, output, err := ParseFlags(os.Args[0], os.Args[1:])

	if err == flag.ErrHelp {
		Error("%s", output)

		os.Exit(UnknownFlag)
	} else if err != nil {
		Error("Invalid args: ", err)

		if len(output) > 0 {
			Error("output:\n", output)
		}

		os.Exit(UnknownFlag)
	}

	return flags, nil
}

func ParseFlags(programName string, args []string) (config *AllFlags, output string, err error) {
	flags := flag.NewFlagSet(programName, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	var allFlags AllFlags
	flags.StringVar(&allFlags.RepositoryUrl, "RepositoryUrl", os.Getenv("REPOSITORY_URL"), "SSH url to the repository to be cloned. Default: $REPOSITORY_URL")
	flags.StringVar(&allFlags.Branch, "Branch", "main", "Branch to check. Default: \"main\"")
	flags.StringVar(&allFlags.SshKey, "SshKey", os.Getenv("SSH_KEY"), "Base64 encoded of SSH private key. Default: $SSH_KEY")
	flags.StringVar(&allFlags.SshPhrase, "SshPhrase", os.Getenv("SSH_PHRASE"), "Base64 encoded of SSH private key. Default: $SSH_PHRASE")
	flags.StringVar(&allFlags.PreviousHash, "PreviousHash", "", "Commit hash of the previous commit to HEAD. Default: $CIRCLE_SHA1")
	flags.StringVar(&allFlags.PropertyPath, "PropertyPath", "version", "Property path to the Version code in the json file. Default: \"Version\"")
	flags.StringVar(&allFlags.FilePath, "FilePath", "package.json", "File path to the json file with the Version code. Default: \"package.json\"")
	flags.StringVar(&allFlags.VersionTagPrefix, "VersionTagPrefix", "v", "Prefix for the git tag")
	flags.StringVar(&allFlags.VersionTagSuffix, "VersionTagSuffix", "", "Suffix for the git tag")

	// Flag short circuits' app run and outputs the binary version from main.Version
	flags.BoolVar(&allFlags.Version, "Version", false, "Check Version of app binary")

	err = flags.Parse(args)
	if err != nil {
		return nil, buf.String(), err
	}

	// User checking Version || help - so we expect other flags not to be set
	skipRemainingFlagChecks := allFlags.Version

	if !skipRemainingFlagChecks {
		flagsToCheck := map[string]reflect.Value{
			"RepositoryUrl": reflect.ValueOf(allFlags.RepositoryUrl),
			"Branch":        reflect.ValueOf(allFlags.Branch),
			"SshKey":        reflect.ValueOf(allFlags.SshKey),
			"PreviousHash":  reflect.ValueOf(allFlags.PreviousHash),
			"PropertyPath":  reflect.ValueOf(allFlags.PropertyPath),
			"FilePath":      reflect.ValueOf(allFlags.FilePath),
		}

		err = ValidateRequiredFlagsSet(flagsToCheck)
		if err != nil {
			return nil, buf.String(), err
		}
	}

	allFlags.args = flags.Args()
	return &allFlags, buf.String(), nil
}

func ValidateRequiredFlagsSet(flags map[string]reflect.Value) error {
	for k, v := range flags {
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
