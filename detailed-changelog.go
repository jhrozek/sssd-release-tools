package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func runGitShortlog(directory string, fromTag string, toTag string) (io.Reader, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	err = os.Chdir(directory)
	if err != nil {
		return nil, err
	}
	defer os.Chdir(cwd)

	cmd := exec.Command("git",
		"shortlog",
		fmt.Sprintf("%s..%s", fromTag, toTag))

	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func formatShortlog(r io.Reader) (string, error) {
	var out bytes.Buffer

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var prefix, suffix string

		changeLine := strings.TrimSpace(scanner.Text())
		if len(changeLine) == 0 {
			continue
		}

		if strings.HasSuffix(changeLine, ":") == false {
			prefix = "      *"
			suffix = ""
		} else {
			prefix = "\n*"
			suffix = "\n"
		}
		out.WriteString(fmt.Sprintln(prefix, changeLine, suffix))
	}

	return out.String(), nil
}

func main() {
	var projDir = flag.String("directory", "", "The location of the git repo")
	var fromTag = flag.String("from-tag", "", "The tag to print the logs since")
	var toTag = flag.String("to-tag", "", "The tag to print the logs until")

	flag.Parse()

	if *projDir == "" {
		fmt.Println("Please specify the repo path with --directory")
		os.Exit(1)
	}

	if *fromTag == "" {
		fmt.Println("Please specify the from tag with --from-tag")
		os.Exit(1)
	}

	if *toTag == "" {
		*toTag = "HEAD"
	}

	shortlogReader, err := runGitShortlog(*projDir, *fromTag, *toTag)
	if err != nil {
		log.Fatal(err)
	}

	shortlog, err := formatShortlog(shortlogReader)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(shortlog)
}
