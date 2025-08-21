package main

import (
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

const versionFile = "version.txt"

var regex = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

func main() {
	cmd := strings.Split("git describe --tags --always --abbrev=0", " ")
	fromGit, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		log.Fatal(err)
	}

	if regex.Match(fromGit) {
		os.Stdout.Write(fromGit)
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		version, err := os.ReadFile(path.Join(cwd, versionFile))
		if err != nil {
			log.Fatal(err)
		}

		os.Stdout.WriteString("v" + strings.Trim(string(version), "\n") + "-dev")
	}
}
