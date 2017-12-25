// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	writeGithash()
}

func writeGithash() {
	hash := getGithash()
	if getGitIsDirty() {
		hash = hash + "~HEAD"
	}
	gitConfigFile := `package config

const VersionGitHash string = "%s"`
	fileContents := fmt.Sprintf(gitConfigFile, hash)
	ioutil.WriteFile("./config/githash.go", []byte(fileContents), 0700)
}

func getGithash() string {
	gitpath, err := exec.LookPath("git")
	if err != nil {
		panic(err)
	}
	cmd := exec.Command(gitpath, "log", "--format=%H", "-n", "1")
	cmd.Dir, _ = filepath.Abs(".")
	r, w, err := os.Pipe()
	cmd.Stdout = w
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	cmd.Wait()
	w.Close()

	x, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	commit := fmt.Sprintf("%s", x)
	commit = strings.Trim(commit, "\n")
	return commit
}

func getGitIsDirty() bool {
	gitpath, err := exec.LookPath("git")
	if err != nil {
		panic(err)
	}
	cmd := exec.Command(gitpath, "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir, _ = filepath.Abs(".")
	r, w, err := os.Pipe()
	cmd.Stdout = w
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	cmd.Wait()
	w.Close()

	x, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	haystack := fmt.Sprintf("%s", x)
	haystack = strings.Trim(haystack, "\t\n ")
	return haystack == "master"
}
