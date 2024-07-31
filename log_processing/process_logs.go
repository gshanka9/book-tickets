package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type LogEntry struct {
	Level string `json:"level"`
	Msg   string `json:"msg"`
	File  string `json:"file"`
	Line  int    `json:"line"`
}

func main() {
	file, err := os.Open("app.log")
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var logEntries []LogEntry

	for scanner.Scan() {
		var entry LogEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			log.Printf("Failed to unmarshal log entry: %v", err)
			continue
		}
		if entry.Level == "error" {
			logEntries = append(logEntries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to scan log file: %v", err)
	}

	if len(logEntries) == 0 {
		log.Println("No errors found in log file.")
		return
	}

	for _, entry := range logEntries {
		baseFile := filepath.Base(entry.File)
		finalPath := "hotel_booking/" + baseFile
		author := getAuthor(finalPath, entry.Line)
		fmt.Printf("Error in file %s at line %d by %s\n", finalPath, entry.Line, author)
	}
}

func getAuthor(file string, line int) string {
	cmd := exec.Command("git", "blame", "-L", fmt.Sprintf("%d,%d", line, line), file)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to execute git blame: %v", err)
		return "unknown1"
	}

	re := regexp.MustCompile(`\((.+?)\s+\d{4}`)
	match := re.FindStringSubmatch(string(output))
	if len(match) < 2 {
		return "unknown2"
	}
	return match[1]
}

func newGitHubClient() *github.Client {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatalf("GITHUB_TOKEN environment variable not set")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}
