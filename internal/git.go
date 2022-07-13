package internal

import (
	"bytes"
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
)

func GetRepository(url string, publicKeys *ssh.PublicKeys) (*git.Repository, error) {
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
	//PrintInfo(encoded)

	//err = os.RemoveAll("bin/repo")
	//r, err := git.PlainClone("bin/repo", false, &git.CloneOptions{
	//	URL:      *Flags.Project_url,
	//	Auth:     publicKeys,
	//	Progress: os.Stdout,
	//})
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

func GetTags(r *git.Repository) (Tags []string, Error error) {
	t, err := r.TagObjects()
	if err != nil {
		return nil, err
	}

	// Retrieves all tags for current Branch
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
	//	PrintInfo("tag.Name: %v\n", tag.Name())
	//	PrintInfo("tag.Type: %v \n\n", tag.Type())
	//	return nil
	//})

	return tags, nil
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

func tagExists(tag string, r *git.Repository) bool {
	tagFoundErr := "tag was found"
	PrintInfo("git show-ref --tag")
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
		PrintInfo("iterate tags error: %s", err)
		return false
	}
	return res
}

func CreateTag(r *git.Repository, tagName string) (Success bool, Error error) {
	if tagExists(tagName, r) {
		PrintWarning("tagName %s already exists", tagName)
		return false, nil
	}
	h, err := r.Head()
	if err != nil {
		PrintWarning("get HEAD error: %s", err)
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
		PrintError("CreateTag error: %s", err)
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
		if err == git.NoErrAlreadyUpToDate {
			PrintInfo("origin remote was up to date, no push done")
			return nil
		}
		PrintInfo("push to remote origin error: %s", err)
		return err
	}

	return nil
}
