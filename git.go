/*
Package git implements a simple library for GIT interactions.
*/
package git

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"

	"github.com/golang/glog"
)
import "strings"

// Error codes returned by failures to parse an expression.
var (
	ErrNoRepo = errors.New("GIT: No repository found in working directory")
)

// Pull fetches the latest version from the upstream branch
func Pull() {

}

// HasChanges checks if there are changes on the GIT repository
func HasChanges() (status bool) {

	output, err := exec.Command("git", "status", "-s", ".").Output()
	if err != nil {
		glog.Fatal(err)
	}

	if string(output) != "" {
		status = true
	} else {
		status = false
	}

	return status
}

// CheckStatus will check if there are changes pending
func CheckStatus() (status bool) {

	status = HasChanges()

	if status == true {
		// glog.Error("Still changes pending, please commit before going furthur.")
		fmt.Printf("Still changes pending, please commit before going furthur.\n")
	}

	return status
}

// TagList will retrieve a list of all tags for the GIT repository
func TagList(searchTag string) (tags []string) {

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
		tags = strings.Split(out.String(), "\n")
	}

	return
}

// GetTagLatest gets the latest tag from GIT
// - if no tag exists, it will return ""
func GetTagLatest() (tag string) {

	// Get tags (default string = v*)
	tags := TagList("v*")
	tag = tags[0]

	return tag
}

// TagExists will verify if the tag already exists
func TagExists(tag string) (status bool) {

	glog.Info("Check if tag exists: ", tag)
	tagResults := TagList(tag)

	if len(tagResults) > 0 {
		if tagResults[0] == tag {
			return true
		}
	}

	return false
}

// HasRemote will check if there is a remote origin configured
func HasRemote() (status bool) {

	output, err := exec.Command("git", "remote", "-v").Output()
	if err != nil {
		glog.Fatal(err)
	}

	if string(output) == "" {
		status = false
	} else {
		status = true
	}

	return
}

// Push will push all new tags onto the upstream branch
func Push() (status bool) {

	output, err := exec.Command("git", "push", "--tags").Output()
	if err != nil {
		glog.Fatal(err)
	}

	if string(output) == "Everything up-to-date" {
		status = true
	} else {
		status = false
	}

	return
}

// CreateTag will tag the current state of repository with the new version
func CreateTag(previouseTagName string, newTagName string) (status bool) {

	if HasChanges() == true {
		fmt.Printf("Changes have been found, please commit first.\n")
		// glog.Error("Changes have been found, please commit first.")
		// glog.Error("@TODO: Ask if we have to commit the files")
		return false
	}

	if previouseTagName != "" {
		if DiffSinceRelease(previouseTagName) == false {
			// glog.Error("No changes since latest release have been found, TAG aborted")
			fmt.Printf("No changes since latest release have been found, TAG aborted\n")
			return false
		}
	}

	if TagExists(newTagName) == true {
		fmt.Printf("Tag already exists in GIT: %s \n", newTagName)
		// glog.Error("Tag already exists in GIT: ", newTagName)
		return false
	}

	// git tag $(TAG)
	output, err := exec.Command("git", "tag", newTagName).Output()
	if err != nil {
		glog.Fatal(err)
		return false
	}
	glog.Info("Output: ", string(output))

	// Check if we have a remote
	if HasRemote() == true {
		glog.Info("We found a remote!")
		Push()
	}

	return true
}

// GetLatestVersion gets the latest tag from GIT and determines the version, verify if this matches our latest version.
func GetLatestVersion() (latestVersionRaw string) {
	latestTag := GetTagLatest()
	latestVersionRaw = strings.Split(latestTag, "v")[1]
	glog.Info("Latest version via GIT: ", latestVersionRaw)

	return
}

// GetRevParse will get the rev-parse from GIT repository - used for --dirty tags.
func GetRevParse() (revparse string) {

	// git rev-parse --short HEAD
	output, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		glog.Fatal(err)
	}

	revparse = strings.TrimSpace(string(output))

	return
}

// DiffSinceRelease will check if the current state differs from latest release
func DiffSinceRelease(releaseTag string) (status bool) {

	glog.Info("DiffSinceRelease: ", releaseTag)

	// git diff --shortstat -r $tag .
	output, err := exec.Command("git", "diff", "--shortstat", "-r", releaseTag, ".").Output()
	if err != nil {
		glog.Error(err)
		return false
	}

	if string(output) != "" {
		// glog.Fatal(string(output))
		status = true
	} else {
		status = false
	}
	glog.Info("Do we have a diff: ", status)

	return status
}
