package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5/config"
	"io"
	"log"
	"os"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/mitchellh/mapstructure"
)

func GetRepository(publicKeys *ssh.PublicKeys) (*git.Repository, error) {
	//home, _ := os.UserHomeDir()
	//// Open file on disk.
	//f, _ := os.Open(home + "/.ssh/update-git-tags_circle-ci")
	//
	//// Read entire JPG into byte slice.
	//reader := bufio.NewReader(f)
	//content, _ := ioutil.ReadAll(reader)
	//
	//// Encode as base64.
	//encoded := b64.StdEncoding.EncodeToString(content)
	//Info(encoded)



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


func tagExists(tag string, r *git.Repository) bool {
	tagFoundErr := "tag was found"
	Info("git show-ref --tag")
	tags, err := r.TagObjects()
	if err != nil {
		log.Printf("get tags error: %s", err)
		return false
	}
	res := false
	err = tags.ForEach(func(t *object.Tag) error {
		if t.Name == tag {
			res = true
			return fmt.Errorf(tagFoundErr)
		}
		return nil
	})
	if err != nil && err.Error() != tagFoundErr {
		Info("iterate tags error: %s", err)
		return false
	}
	return res
}

func CreateTag(r *git.Repository, tag string) (bool, error) {
	if tagExists(tag, r) {
		Warning("tag %s already exists", tag)
		return false, nil
	}
	Warning("Set tag %s", tag)
	h, err := r.Head()
	if err != nil {
		Warning("get HEAD error: %s", err)
		return false, err
	}
	Info("git tag -a %s %s -m \"%s\"", tag, h.Hash(), tag)
	_, err = r.CreateTag(tag, h.Hash(), &git.CreateTagOptions{
		Message: tag,
	})

	if err != nil {
		Warning("create tag error: %s", err)
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
	Info("git push --tags")
	err := r.Push(po)

	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			Info("origin remote was up to date, no push done")
			return nil
		}
		Info("push to remote origin error: %s", err)
		return err
	}

	return nil
}