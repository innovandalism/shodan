// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	bindata()
	writeGithash()
}

func bindata() {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		panic("GOPATH not set")
	}
	bindataPath := path.Join(gopath, "/bin/go-bindata")
	cmd := exec.Command(bindataPath, "-nometadata", "-o", "bindata/bindata.go", "-pkg", "bindata", "assets/...")
	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH") + ":" + path.Join(gopath, "bin"),
	}
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	cmd.Wait()
}

func writeGithash() {
	hash := getGithash()
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
	cmd := exec.Command(gitpath, "status", "-s")
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
	return len(haystack) > 0
}
