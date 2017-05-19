// Package gitwrap is a versioning library that helps generate proper version numbers.
package gitwrap

import "os"

// Repo struct
type Repo struct {
	workdir string
}

// IsDirty will return the dirty status of the repository
func (v *Repo) IsDirty() bool {
	os.Chdir(v.workdir)
	return isDirty()
}

// Tags retrieves all known tags
func (v *Repo) Tags(query string) (strTags []string) {
	strTags = tags(query)
	return strTags
}

// HasRemote verifies if we have upstream activity
func (v *Repo) HasRemote() (status bool) {
	os.Chdir(v.workdir)
	remote, _ := hasRemote()
	return remote
}

// GetRevParse retrieves the revparse from the repository. Used in formatting the versionstring
func (v *Repo) GetRevParse() string {
	strRevParse, _ := getRevParse()
	return strRevParse
}

// DiffSinceRelease retrieves the revparse from the repository. Used in formatting the versionstring
func (v *Repo) DiffSinceRelease(releaseTag string) bool {
	diffBool, _ := diffSinceRelease(releaseTag)
	return diffBool
}

// Pull retrieves latest state from remote
func (v *Repo) Pull() {

}

// TagCreate retrieves latest state from remote
func (v *Repo) TagCreate(newTagName string, message string) (status bool, err error) {
	return tagCreate(newTagName, message)
}

// TagLatest retrieves latest tag from Repo
func (v *Repo) TagLatest(tagPrefix string) (tag string, err error) {
	return tagLatest(tagPrefix)
}

// TagExists retrieves latest tag from Repo
func (v *Repo) TagExists(tag string) bool {
	return tagExists(tag)
}
