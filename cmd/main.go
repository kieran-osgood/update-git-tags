package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"go-git-tag/internal"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

//goland:noinspection GoUnusedGlobalVariable
var (
	version string
	build string
)

func main() {
	flags, err := internal.GetFlags()
	internal.HandleError(err)

	if flags.Version {
		if version == "" {
			fmt.Printf("Version code hasn't been set\n")
		}
		fmt.Printf(version)
		os.Exit(2)
	}

	key, _ := b64.StdEncoding.DecodeString(flags.SshKey)
	publicKeys, err := ssh.NewPublicKeys("git", key, "")
	internal.HandleError(err)

	r, err := internal.GetRepository(flags.RepositoryUrl, publicKeys)
	internal.HandleError(err)

	ref, err := r.Head()
	internal.HandleError(err)

	w, err := r.Worktree()
	internal.HandleError(err)

	v1, err := GetVersion(w, *flags)
	internal.HandleError(err)

	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(flags.PreviousHash),
	})
	internal.HandleError(err)

	v2, err := GetVersion(w, *flags)
	internal.HandleError(err)

	if v1 == v2 {
		internal.Warning("Version code hasn't changed, exiting")
		os.Exit(0)
	} else {
		internal.Info("Version code changed!")
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: ref.Hash(),
	})
	internal.HandleError(err)

	_, err = internal.CreateTag(r, fmt.Sprintf("%v%v%v", flags.VersionTagPrefix, v1, flags.VersionTagSuffix))
	internal.HandleError(err)

	err = internal.PushTags(r, publicKeys)
	internal.HandleError(err)
}

func GetVersion(w *git.Worktree, flags internal.AllFlags) (string, error) {
	j, err := internal.ReadJson(w, flags.FilePath)
	if err != nil {
		return "", err
	}

	v, err := ExtractVersionFromJson(j, flags.PropertyPath)
	if err != nil {
		return "", err
	}

	if v != "" {
		internal.Info("Branch: %v \nversion: %v\n", flags.Branch, v)
	} else {
		internal.Warning("Version: %v\n was not found - check the path", v)
	}

	return v, nil
}

func ExtractVersionFromJson(jsonString *[]byte, accessor string) (string, error) {
	var result map[string]interface{}
	err := json.Unmarshal(*jsonString, &result)
	if err != nil {
		return "", err
	}
	nextLevelToCheck := result

	var s string
	kArr := strings.Split(accessor, ".")
	for i := 0; i < len(kArr); i++ {
		k := kArr[i]
		if v, ok := nextLevelToCheck[k].(string); ok {
			s = v
			break
		}

		replacement, ok := nextLevelToCheck[k].(map[string]interface{})
		if !ok {
			internal.Error("Error accessing property: %v. Check Property_path flag matches json path to version property.", k)
			break
		}
		nextLevelToCheck = replacement
	}

	return s, nil
}
