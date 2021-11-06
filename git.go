package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
)

func GetGitFiles() (*git.Repository, error) {
	publicKeys, err := ssh.NewPublicKeysFromFile("git", "", "")
	if err != nil {
		return nil, err
	}
	/**
	  Next steps:
	  	create a new ssh key - test it hard coded in here first
	  		then add ssh key to circleci
	  	pull repo, copy version value of master-head
	  	copy version value of master head~1

	  compare values
	*/
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:           *Flags.Project_url,
		Auth:          publicKeys,
		ReferenceName: "master",
	})
	if err != nil {
		return nil, err
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(*Flags.Previous_sha),
	})
	if err != nil {
		return nil, err
	}

	return r, nil
}

