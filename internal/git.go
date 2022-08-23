package internal

import (
	"bytes"
	"github.com/go-git/go-git/v5/config"
	"io"
	"os"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
)

func GetRepository(url string, publicKeys *ssh.PublicKeys) (*git.Repository, error) {
	r, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL:      url,
		Auth:     publicKeys,
		Progress: os.Stdout,
	})
	if err != nil {
		return nil, err
	}

	return r, nil
}

func ReadJson(w *git.Worktree, path string) (*[]byte, error) {
	f, err := w.Filesystem.Open(path)
	if err != nil {
		return nil, err
	}

	var b []byte
	for {
		buffer := make([]byte, 8)
		_, err := f.Read(buffer)
		b = append(b, buffer...)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
	}
	b = bytes.Trim(b, "\x00")

	return &b, nil
}

func CreateTag(r *git.Repository, tagName string) (Success bool, Error error) {
	h, err := r.Head()
	if err != nil {
		PrintWarning("r.HEAD: %w", err)
		return false, err
	}
	PrintInfo("Set tagName %s", tagName)
	PrintInfo("git tagName -a %s %s -m \"%s\"", tagName, h.Hash(), tagName)
	_, err = r.CreateTag(
		tagName,
		h.Hash(),
		&git.CreateTagOptions{
			Message: tagName,
		},
	)

	if err != nil {
		PrintError("CreateTag: %w", err)
		return false, err
	}

	return true, nil
}

func PushTags(r *git.Repository, publicKeys *ssh.PublicKeys) error {
	po := &git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		RefSpecs:   []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")},
		Auth:       publicKeys,
	}
	PrintInfo("git push --tags")
	err := r.Push(po)

	if err != nil {
		switch err {
		case git.NoErrAlreadyUpToDate:
			PrintInfo("origin remote was up to date, no push done")
			return nil
		default:
			PrintInfo("push to remote origin error: %w", err)
			return err
		}
	}

	return nil
}
