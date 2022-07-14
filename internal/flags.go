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
		PrintError("%s", output)

		os.Exit(UnknownFlag)
	} else if err != nil {
		PrintError("Invalid args: %v", err)

		if len(output) > 0 {
			PrintError("output: \n", output)
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
	flags.StringVar(
		&allFlags.RepositoryUrl,
		"RepositoryUrl",
		os.Getenv("REPOSITORY_URL"),
		"SSH url to the repository to be cloned. e.g. git@github.com:kieran-osgood/update-git-tags.git",
	)
	flags.StringVar(
		&allFlags.Branch,
		"Branch",
		"main",
		"Branch to check.",
	)
	flags.StringVar(
		&allFlags.SshKey,
		"SshKey",
		os.Getenv("SSH_KEY"),
		"Base64 encoded of SSH private key.",
	)
	flags.StringVar(
		&allFlags.SshPhrase,
		"SshPhrase",
		os.Getenv("SSH_PHRASE"),
		"Passphrase to the supplied --SshKey.",
	)
	flags.StringVar(
		&allFlags.PreviousHash,
		"PreviousHash",
		"",
		"Commit hash of the previous commit to HEAD.",
	)
	flags.StringVar(
		&allFlags.PropertyPath,
		"PropertyPath",
		"version",
		"Dot separated property path to the version field in the json file. e.g. for a json file like {\"expo:\": { \"version\": \"1.0.0\" }} you would supply: --PropertyPath='expo.version'",
	)
	flags.StringVar(
		&allFlags.FilePath,
		"FilePath",
		"package.json",
		"File path to the json file storing the apps version.",
	)
	flags.StringVar(
		&allFlags.VersionTagPrefix,
		"VersionTagPrefix",
		"v",
		"Added to the start of each new git tag this CLI creates.",
	)
	flags.StringVar(
		&allFlags.VersionTagSuffix,
		"VersionTagSuffix",
		"",
		"Added to the end of each new git tag this CLI creates.",
	)

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
