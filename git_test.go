package git

import (
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strings"
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

func setDir() {
	// Set current work-dir
	os.Chdir("/Workspace/playground/demo-repo")
}

func TestHasNoChanges(t *testing.T) {

	// Create workspace
	setDir()
	resetChanges()

	var v bool
	v = HasChanges()
	if v != false {
		t.Error("Expected true, got ", v)
	}

}

func TestHasChanges(t *testing.T) {

	// Create workspace
	setDir()
	makeChanges()

	var v bool
	v = HasChanges()
	if v != true {
		t.Error("Expected false, got ", v)
	}

}

func TestCheckStatusTrue(t *testing.T) {

	// Create workspace
	setDir()
	resetChanges()

	var v bool
	v = CheckStatus()
	if v == true {
		t.Error("Expected false, got ", v)
	}
}

func TestCheckStatusFalse(t *testing.T) {

	// Create workspace
	setDir()
	makeChanges()

	var v bool
	v = CheckStatus()
	if v == false {
		t.Error("Expected false, got ", v)
	}
}

func TestTagListEmpty(t *testing.T) {

	// Create workspace
	setDir()
	Cleanup()
	resetChanges()

	var v []string

	v = TagList("")

	if len(v) > 0 {
		t.Error("Expected to have an empty tag list")
	}

}

// func TestTagListWithEntry(t *testing.T) {
//
// 	// Create workspace
// 	setDir()
// 	makeChanges()
// 	tagChanges()
//
// 	var v []string
//
// 	v = TagList("")
//
// 	if len(v) > 0 {
// 		t.Error("Expected to have an empty tag list")
// 	}
//
// }

func TestCreateTagFail(t *testing.T) {

	makeChanges()

	var v bool
	v = CreateTag("", "v0.0.1")

	if v != false {
		t.Error("Expected false, got", v)
	}

	Cleanup()
	resetChanges()

}

func TestCreateTagFailStillChanges(t *testing.T) {

	// Create workspace
	setDir()
	resetChanges()
	makeChanges()

	var v bool
	v = CreateTag("v0.0.0", "v0.0.1")

	if v != false {
		t.Error("Expected false, got", v)
	}

	Cleanup()

}

func TestCreateTagPass(t *testing.T) {

	// Create workspace
	setDir()
	commitChanges()

	var v bool
	v = CreateTag("", "v0.0.1")

	if v != true {
		t.Error("Expected true, got", v)
	}

	Cleanup()

}

func TestCreateTagFailTagExists(t *testing.T) {

	// Create workspace
	Cleanup()
	setDir()
	commitChanges()
	makeDummyTag("v0.0.1")

	var v bool
	v = CreateTag("v0.0.0", "v0.0.1")

	if v != false {
		t.Error("Expected false, got", v)
	}

}

func TestGetRevParse(t *testing.T) {

	v, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		glog.Fatal(err)
	}

	revparse := strings.TrimSpace(string(v))

	if revparse != GetRevParse() {
		t.Errorf("Expected %s. got %s", revparse, v)
	}

}

func Cleanup() {
	exec.Command("git", "tag", "-d", "v0.0.1").Run()
	exec.Command("git", "tag", "-d", "v0.0.0").Run()
	resetChanges()
}
