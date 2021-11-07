package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func main() {
	err := ParseFlags()
	HandleError(err)

	r, err := GetRepository()
	HandleError(err)

	//ref, err := r.Head()
	//HandleError(err)
	
	w, err := r.Worktree()
	HandleError(err)

	v1, err := ReadAppJson(w)
	HandleError(err)
	Info(v1.Expo.Version)

	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash("6b3887536d1c17d49b305a83b8bc2693681cdc22"),
	})
	HandleError(err)

	v2, err := ReadAppJson(w)
	HandleError(err)
	Info(v2.Expo.Version)

	//Info(result.Expo.Version)
	//currentVersionCode := result.Expo.Version
	//oldVersionCode := "1.2.3"

	//if currentVersionCode != oldVersionCode {
	//	_, err := r.CreateTag(currentVersionCode, plumbing.NewHash(*Flags.Current_sha), &git.CreateTagOptions{})
	//	HandleError(err)
	//}
}
