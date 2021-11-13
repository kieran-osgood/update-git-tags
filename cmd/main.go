package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"go-git-tag/internal"
	"os"
	"strings"
)

func main() {
	flags, err := internal.GetFlags()
	internal.HandleError(err)

	key, _ := b64.StdEncoding.DecodeString(flags.SshKey)
	publicKeys, err := ssh.NewPublicKeys("git", key, "")
	internal.HandleError(err)

	r, err := internal.GetRepository(flags.RepositoryUrl, publicKeys)
	internal.HandleError(err)

	ref, err := r.Head()
	internal.HandleError(err)

	w, err := r.Worktree()
	internal.HandleError(err)

	j1, err := internal.ReadJson(w, "app.json")
	internal.HandleError(err)
	v1, err := ExtractVersionFromJson(j1, flags.PropertyPath)
	internal.HandleError(err)
	internal.Info("Branch: %v \nversion: %v\n", flags.Branch, v1)
	if v1 == "" {
		internal.Warning("Branch: %v \nversion: %v\n", flags.Branch, v1)
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(flags.PreviousHash),
	})
	internal.HandleError(err)

	j2, err := internal.ReadJson(w, "app.json")
	internal.HandleError(err)
	v2, err := ExtractVersionFromJson(j2, flags.PropertyPath)
	internal.HandleError(err)
	if v2 == "" {
		internal.Warning("Branch: %v \nversion: %v\n", flags.Branch, v2)
	}
	internal.Info("Commit: %v \nversion: %v\n", flags.PreviousHash, v2)

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

	_, err = internal.CreateTag(r, fmt.Sprintf("v%v", v1))
	internal.HandleError(err)

	err = internal.PushTags(r, publicKeys)
	internal.HandleError(err)
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
