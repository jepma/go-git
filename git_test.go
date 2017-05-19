package gitwrap

import (
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/golang/glog"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func resetChanges() {
	// Will make sure it has NO changes - RESET HARD!
	exec.Command("git", "reset", "HEAD", "--hard").Run()
	exec.Command("git", "clean", "-fd").Run()
}

func makeChanges() {
	// Will make sure it has changes
	d1 := []byte("hello\ngo\n")
	ioutil.WriteFile(RandStringRunes(7), d1, 0644)
}

func makeDummyTag(tag string) {

	err := exec.Command("git", "tag", tag).Run()
	if err != nil {
		glog.Fatal("Error: ", err)
	}

}

func addRemote() {
	err := exec.Command("git", "remote", "add", "origin", "git@github.com:jepma/demo-repo.git").Run()
	if err != nil {
		glog.Fatal("Error: ", err)
	}
}

func commitChanges() {

	resetChanges()
	makeDummyTag("v0.0.0")
	makeChanges()

	err := exec.Command("git", "add", ".").Run()
	if err != nil {
		glog.Fatal("Error: ", err)
	}

	err = exec.Command("git", "commit", "-m=blaat").Run()
	if err != nil {
		glog.Fatal("Error: ", err)
	}
	//
	// duration := time.Duration(2) * time.Second // Pause for 10 seconds
	// time.Sleep(duration)

}

func Cleanup() {

	// Reset GIT to first commit
	// git rev-list --max-parents=0 --abbrev-commit HEAD
	exec.Command("git", "rev-list", "--max-parents=0", "--abbrev-commit", "HEAD").Run()

	// Remove all tags
	for _, element := range tags("") {
		exec.Command("git", "tag", "-d", element).Run()
	}

	exec.Command("git", "clean", "-fd").Run()
}

func setDir() {
	// Set current work-dir
	os.Chdir("/Workspace/playground/demo-repo")
}

func getDir() string {
	dir, _ := os.Getwd()
	return dir
}

func TestOpenRepo(t *testing.T) {
	// Create workspace
	setDir()
	resetChanges()

	repo, _ := CreateRepoObject(getDir())

	if repo.workdir != getDir() {
		t.Error("Expected workdir to be ", getDir())
	}

}

func TestOpenRepoFail(t *testing.T) {

	// Create workspace
	setDir()
	resetChanges()

	_, err := CreateRepoObject("/tmp/blaat")
	if err == nil {
		t.Error("Expected error, got", err)
	}

}

func TestHasChangesFalse(t *testing.T) {

	// Create workspace
	setDir()

	// Create repo object
	repo, _ := CreateRepoObject(getDir())

	// Check for changes
	if repo.IsDirty() != false {
		t.Error("Expected hasChanges to be false, got ", repo.IsDirty())
	}

	Cleanup()
}
func TestChangesTrue(t *testing.T) {

	// Create workspace
	setDir()
	makeChanges()
	repo, _ := CreateRepoObject(getDir())

	v := repo.IsDirty()
	if v == false {
		t.Error("Expected true, got ", v)
	}

	Cleanup()

}

func TestTagListEmpty(t *testing.T) {

	// Create workspace
	setDir()

	var v []string

	v = tags("")

	if len(v) > 0 {
		t.Error("Expected to have an empty tag list")
	}

	Cleanup()

}

func TestTagListWithOneEntry(t *testing.T) {

	// Create workspace
	setDir()
	makeChanges()
	makeDummyTag("v0.0.1")

	var v []string
	v = tags("")

	if len(v) > 1 {
		t.Error("Expected to have one entry, got ", len(v), v)
	}

	Cleanup()

}

func TestTagLatest(t *testing.T) {

	// Create workspace
	setDir()
	makeChanges()
	makeDummyTag("v0.0.1")
	makeDummyTag("v0.0.2")

	var v string
	v, err := tagLatest("v")

	if v != "v0.0.2" {
		t.Error("Expected to get v0.0.2, got ", v)
	}

	if err != nil {
		t.Error("Expected err to be nil, got", err)
	}

	Cleanup()

}

func TestTagExistsTrue(t *testing.T) {

	// Create workspace
	setDir()
	makeChanges()
	makeDummyTag("v0.0.1")

	var v bool
	v = tagExists("v0.0.1")

	if v == false {
		t.Error("Expected true from tagExists, got ", v)
	}

	Cleanup()

}

func TestTagExistsFalse(t *testing.T) {

	// Create workspace
	setDir()
	makeChanges()

	var v bool
	v = tagExists("v0.0.1")

	if v != false {
		t.Error("Expected false from tagExists, got ", v)
	}

	Cleanup()

}

func TestHasRemoteTrue(t *testing.T) {

	// Create workspace
	setDir()

	// Add remote
	exec.Command("git", "remote", "add", "origin", "git@github.com:jepma/demo-repo.git").Output()

	var v bool
	v, err := hasRemote()

	if v != true {
		t.Error("Expected true from hasRemote, got ", v)
	}
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	// Delete remote
	exec.Command("git", "remote", "remove", "origin").Output()

	Cleanup()

}

func TestHasRemoteFalse(t *testing.T) {

	// Create workspace
	setDir()

	var v bool
	v, err := hasRemote()

	if v != false {
		t.Error("Expected false from hasRemote, got", v)
	}
	if err != ErrNoRemote {
		t.Error("Expected ErrNoRemote error, got ", err)
	}

	Cleanup()

}

func TestGetRevParse(t *testing.T) {

	// Create workspace
	setDir()

	var v string
	v, err := getRevParse()

	if err != nil || v == "" {
		t.Error("Expected string and no error from getRevParse, got", v, err)
	}

	Cleanup()

}

func TestDifSinceReleaseTrue(t *testing.T) {

	// Create workspace
	setDir()
	makeDummyTag("v0.0.1")
	makeChanges()
	commitChanges()
	var v bool
	v, err := diffSinceRelease("v0.0.1")

	if err != nil {

		t.Error("Expected no error from diffSinceRelease, got", err)
	}

	if v != true {

		t.Error("Expected true from diffSinceRelease, got", v)
	}

	Cleanup()

}

func TestDifSinceReleaseFalse(t *testing.T) {

	// Create workspace
	setDir()
	makeDummyTag("v0.0.1")
	makeChanges()

	var v bool
	v, err := diffSinceRelease("v0.0.1")
	if err != nil || v != false {
		t.Error("Expected false and no error from diffSinceRelease, got", v, err)
	}

	Cleanup()

}

func TestDifSinceReleaseFalseTagDoesNotExist(t *testing.T) {

	// Create workspace
	setDir()
	makeDummyTag("v0.0.1")
	makeChanges()

	var v bool
	v, err := diffSinceRelease("v1.0.0")

	if err == nil {

		t.Error("Expected error from diffSinceRelease, got", err)
	}

	if v != false {

		t.Error("Expected false from diffSinceRelease, got", v)
	}

	Cleanup()

}

func TestCreateTagSuccess(t *testing.T) {

	// Create workspace
	setDir()
	makeDummyTag("v0.0.1")
	makeChanges()
	commitChanges()

	var v bool
	var err error
	v, err = tagCreate("v0.0.2", "")

	if v != true && err != nil {
		t.Error("Expected true and no error, got", v, err)
	}

	Cleanup()
	resetChanges()

}

func TestCreateTagErrorAlreadyExists(t *testing.T) {

	// Create workspace
	setDir()
	makeDummyTag("v0.0.1")
	makeChanges()
	commitChanges()

	var v bool
	var err error
	v, err = tagCreate("v0.0.1", "")

	if v != false && err != ErrTagAlreadyExists {
		t.Error("Expected true and no error, got", v, err)
	}

	Cleanup()
	// resetChanges()

}

func TestCreateTagErrTagCreateFalseName(t *testing.T) {

	// Create workspace
	setDir()
	makeDummyTag("v0.0.1")

	var v bool
	var err error
	v, err = tagCreate("!#!@$$", "")

	if v != false {
		t.Error("Expected false, got", v)
	}

	if err != ErrTagCreateNoChanges {
		t.Error("Expected ErrTagCreateNoChanges, got", err)
	}

	Cleanup()
	// resetChanges()

}

/*
	Will not do this now, slows down testing.
*/
func DisabledTestPushTrue(t *testing.T) {

	// Create workspace
	setDir()
	makeDummyTag("v0.0.1")

	// Add remote
	exec.Command("git", "remote", "add", "origin", "git@github.com:jepma/demo-repo.git").Output()

	v, err := pushRemote()

	if v != true {
		t.Error("Expected true, got", v)
	}

	if err != nil {
		t.Error("Expected nil, got", err)
	}

	Cleanup()

	// Delete remote
	exec.Command("git", "remote", "remove", "origin").Output()

}

func TestPushFalse(t *testing.T) {

	// Create workspace
	setDir()
	makeDummyTag("v0.0.1")

	v, err := pushRemote()

	if v != false {
		t.Error("Expected false, got", v)
	}

	if err != ErrNoRemote {
		t.Error("Expected ErrNoRemote, got", err)
	}

	Cleanup()

}

/*
	Will not do this now, slows down testing.
*/
func DisabledTestPullTrue(t *testing.T) {

	// Create workspace
	setDir()
	makeDummyTag("v0.0.1")

	// Add remote
	exec.Command("git", "remote", "add", "origin", "git@github.com:jepma/demo-repo.git").Output()

	v, err := pushRemote()

	if v != true {
		t.Error("Expected true, got", v)
	}

	if err != nil {
		t.Error("Expected nil, got", err)
	}

	Cleanup()

	// Delete remote
	exec.Command("git", "remote", "remove", "origin").Output()

}

func TestPullFalse(t *testing.T) {

	// Create workspace
	setDir()
	makeDummyTag("v0.0.1")

	v, err := pullRemote()

	if v != false {
		t.Error("Expected false, got", v)
	}

	if err != ErrNoRemote {
		t.Error("Expected ErrNoRemote, got", err)
	}

	Cleanup()

}

//
// func TestCreateTagFailStillChanges(t *testing.T) {
//
// 	// Create workspace
// 	setDir()
// 	resetChanges()
// 	makeChanges()
//
// 	var v bool
// 	v = CreateTag("v0.0.0", "v0.0.1")
//
// 	if v != false {
// 		t.Error("Expected false, got", v)
// 	}
//
// 	Cleanup()
//
// }
//
// func TestCreateTagPass(t *testing.T) {
//
// 	// Create workspace
// 	setDir()
// 	commitChanges()
//
// 	var v bool
// 	v = CreateTag("", "v0.0.1")
//
// 	if v != true {
// 		t.Error("Expected true, got", v)
// 	}
//
// 	Cleanup()
//
// }
//
// func TestCreateTagFailTagExists(t *testing.T) {
//
// 	// Create workspace
// 	Cleanup()
// 	setDir()
// 	commitChanges()
// 	makeDummyTag("v0.0.1")
//
// 	var v bool
// 	v = CreateTag("v0.0.0", "v0.0.1")
//
// 	if v != false {
// 		t.Error("Expected false, got", v)
// 	}
//
// }
//
// func TestGetRevParse(t *testing.T) {
//
// 	v, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
// 	if err != nil {
// 		glog.Fatal(err)
// 	}
//
// 	revparse := strings.TrimSpace(string(v))
//
// 	if revparse != GetRevParse() {
// 		t.Errorf("Expected %s. got %s", revparse, v)
// 	}
//
// }
//
