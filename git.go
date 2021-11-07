package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"io"
	"os"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/mitchellh/mapstructure"
)

func GetRepository() (*git.Repository, error) {
	key, _ := b64.StdEncoding.DecodeString(*Flags.Ssh_Key)
	publicKeys, err := ssh.NewPublicKeys("git", key, "")
	if err != nil {
		return nil, err
	}

	//err = os.RemoveAll("bin/repo")
	//r, err := git.PlainClone("bin/repo", false, &git.CloneOptions{
	//	URL:      *Flags.Project_url,
	//	Auth:     publicKeys,
	//	Progress: os.Stdout,
	//})
	r, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL:      *Flags.Project_url,
		Auth:     publicKeys,
		Progress: os.Stdout,
	})
	if err != nil {
		return nil, err
	}

	return r, nil
}

func GetTags(r *git.Repository) ([]string, error) {
	t, err := r.TagObjects()
	if err != nil {
		return nil, err
	}

	// Retrieves all tags for current branch
	var tags []string
	err = t.ForEach(func(t *object.Tag) error {
		tags = append(tags, t.Name)
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Retrieves *all* tags
	//tags, err := r.Tags()
	//err = tags.ForEach(func(tag *plumbing.Reference) error {
	//	Info("tag.Name: %v\n", tag.Name())
	//	Info("tag.Type: %v \n\n", tag.Type())
	//	return nil
	//})

	return tags, nil
}

func CreateTag(r *git.Repository, name string) error {
	head, err := r.Head()
	if err != nil {
		return err
	}

	_, err = r.CreateTag(name, head.Hash(), &git.CreateTagOptions{})
	if err != nil {
		return err
	}

	return nil
}

type AppJsonFile struct {
	Expo struct {
		Version string
	}
}

func ReadAppJson(w *git.Worktree) (*AppJsonFile, error) {
	f, err := w.Filesystem.Open("app.json")
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

	m := make(map[string]interface{})
	err = json.Unmarshal(bytes.Trim(b, "\x00"), &m)
	if err != nil {
		return nil, err
	}

	var result AppJsonFile
	err = mapstructure.Decode(m, &result)
	HandleError(err)

	return &result, nil
}
