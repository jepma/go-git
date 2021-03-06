/*
Package gitwrap implements a simple library for GIT interactions.
*/
package gitwrap

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/golang/glog"
)

// CreateRepoObject will set the path for the repository and will create go-git repo object.
func CreateRepoObject(path string) (RepoStruct Repo, err error) {

	_, err = os.Stat(path)
	if err != nil {
		return RepoStruct, err
	}

	RepoStruct.workdir = path

	return RepoStruct, err
}

//
//
//
//
//

// isDirty will trigger the GIT status --porcelain command
func isDirty() (status bool) {

	// git status --porcelain
	output, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		glog.Fatal(err)
	}

	if string(output) == "" {
		status = false
	} else {
		status = true
	}

	return status
}

// tags will retrieve a list of all tags for the GIT repository
func tags(searchTag string) (tags []string) {

	// If empty searchTag is given, use wildcard.
	if searchTag == "" {
		searchTag = "v*"
	}

	// git tag -l aws-cli --sort=version:refname
	cmd := exec.Command("git", "tag", "-l", searchTag, "--sort=-version:refname")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		glog.Fatal(err)
	}

	if out.String() != "" {
		tags = strings.Split(strings.TrimSpace(out.String()), "\n")
	}

	return
}

// GetTagLatest gets the latest tag from GIT. If no tag exists, it will return ""
func tagLatest(tagPrefix string) (tag string, err error) {

	// Get tags (default string = v*)
	// tags := tags("v*")
	tags := tags(fmt.Sprintf("%s*", tagPrefix))

	// Check if we have actual results
	if len(tags) == 0 {
		return tag, ErrTagEmpty
	}

	return tags[0], nil
}

// TagExists will verify if the tag already exists
func tagExists(tag string) (status bool) {

	tagResults := tags(tag)

	if len(tagResults) > 0 {
		if tagResults[0] == tag {
			return true
		}
	}

	return false
}

func tagCreate(newTagName string, message string) (status bool, err error) {

	// git tag $(TAG)
	output, err := exec.Command("git", "tag", newTagName).CombinedOutput()

	if err != nil {
		glog.Errorf("Error: %s", output)
		return false, ErrTagCreating
	}

	if string(output) != "" {
		glog.Errorf("Output: %s", string(output))
		return false, err
	}

	return true, err
}

// HasRemote will check if there is a remote origin configured
func hasRemote() (status bool, err error) {

	output, err := exec.Command("git", "remote", "-v").Output()
	if err != nil {
		glog.Error(err)
		return false, ErrNoRemote
	}

	if string(output) == "" {
		return false, ErrNoRemote
	}
	return true, nil
}

func pushRemote() (status bool, err error) {

	// Check if we have upstream
	if statRemote, errRemote := hasRemote(); statRemote == false {
		return false, errRemote
	}

	output, err := exec.Command("git", "push", "--tags").Output()
	if err != nil {
		glog.Error(err)
		return false, err
	}

	if string(output) == "Everything up-to-date" {
		status = true
	} else {
		status = false
	}

	return status, err
}

// Pull remote repo
func pullRemote() (status bool, err error) {

	// Check if we have upstream
	if statRemote, errRemote := hasRemote(); statRemote == false {
		return false, errRemote
	}

	output, err := exec.Command("git", "pull", "--tags").Output()
	if err != nil {
		glog.Error(err)
		return false, err
	}

	if string(output) == "Everything up-to-date" {
		status = true
	} else {
		status = false
	}

	return status, err
}

// getRevParse will get the rev-parse from GIT repository - used for --dirty tags.
func getRevParse() (revparse string, er error) {

	// git rev-parse --short HEAD
	output, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		glog.Error(err)
		return "", err
	}

	revparse = strings.TrimSpace(string(output))

	return
}

// DiffSinceRelease will check if the current state differs from latest release
func diffSinceRelease(releaseTag string) (status bool, err error) {

	var outbuf, errbuf bytes.Buffer

	// git diff --shortstat -r $tag .
	cmd := exec.Command("git", "diff", "--shortstat", "-r", releaseTag, ".")
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	// Start command
	if err = cmd.Run(); err != nil {

		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if _, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				// glog.Error("DiffSinceRelease Exit Status: ", status.ExitStatus())
				return false, ErrDiff
			}
		} else {
			// log.Fatalf("DiffSinceRelease cmd.Wait: %v", err)
			return false, ErrDiff
		}

		glog.Fatalf("Error starting command: %s", err)
	}

	// Fetch output
	stdout := outbuf.String()

	// If there are no changes, we get an empty response
	if stdout == "" {
		return false, err
	}

	return true, err
}
