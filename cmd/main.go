package main

import (
	"encoding/json"
	"fmt"
	"go-git-tag/internal"
	"strings"
)

func main() {
	test()
	//err := ParseFlags()
	//HandleError(err)
	//
	//key, _ := b64.StdEncoding.DecodeString(*Flags.Ssh_Key)
	//publicKeys, err := ssh.NewPublicKeys("git", key, "")
	//HandleError(err)
	//
	//r, err := GetRepository(publicKeys)
	//HandleError(err)
	//
	//ref, err := r.Head()
	//HandleError(err)
	//
	//w, err := r.Worktree()
	//HandleError(err)
	//
	//v1, err := ReadAppJson(w)
	//HandleError(err)
	//
	//Info("Branch: %v \nversion: %v\n", *Flags.Branch, v1.Expo.Version)
	//
	//err = w.Checkout(&git.CheckoutOptions{
	//	Hash: plumbing.NewHash(*Flags.Previous_Hash),
	//})
	//HandleError(err)
	//
	//v2, err := ReadAppJson(w)
	//HandleError(err)
	//Info("Commit: %v \nversion: %v\n",*Flags.Previous_Hash, v2.Expo.Version)
	//
	//if v1.Expo.Version == v2.Expo.Version {
	//	Info("Version code hasn't changed, exiting")
	//	os.Exit(0)
	//}
	//
	//err = w.Checkout(&git.CheckoutOptions{
	//	Hash: ref.Hash(),
	//})
	//HandleError(err)
	//
	//_, err = CreateTag(r, fmt.Sprintf("v%v", v1.Expo.Version))
	//HandleError(err)
	//
	//err = PushTags(r, publicKeys)
	//HandleError(err)
}

type PackageJson struct {
	Expo ExpoContent
}
type ExpoContent struct {
	Version string
}

func test() {
	//js := `{"version": "1.0.0"}`
	js := `{"expo":{"version":"1.0.0"}}`
	var result map[string]interface{}
	err := json.Unmarshal([]byte(js), &result)
	if err != nil {
		return
	}
	nextLevelToCheck := result

	var s string
	//kArr := strings.Split("version", ".")
	kArr := strings.Split("expo.version", ".")
	//kArr := strings.Split(*internal.Flags.Property_path, ".")
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
	fmt.Printf("s: %v", s)
}
