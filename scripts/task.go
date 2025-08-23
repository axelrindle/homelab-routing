package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

const versionFile = "version.txt"

var regex = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

func command(cmd string) []byte {
	args := strings.Split(cmd, " ")
	stdout, err := exec.Command(args[0], args[1:]...).Output()

	if err != nil {
		log.Fatal(err)
	}

	return bytes.TrimSpace(stdout)
}

func main() {
	fromGit := command("git describe --tags --always --abbrev=0")

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

		manual := fmt.Sprintf("%s+dev%s",
			bytes.TrimSpace(version),
			os.Getenv("GITHUB_RUN_NUMBER"))
		os.Stdout.WriteString(manual)
	}
}
