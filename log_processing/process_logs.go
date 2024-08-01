package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type LogEntry struct {
	Level  string `json:"level"`
	Msg    string `json:"msg"`
	File   string `json:"file"`
	Line   int    `json:"line"`
	Author string `json:"author"`
}

func main() {
	http.HandleFunc("/errors", errorLogHandler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func errorLogHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}

	log.Println("Processing error log request")

	file, err := os.Open("app.log")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to open log file: %v", err), http.StatusInternalServerError)
		return
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
			baseFile := filepath.Base(entry.File)
			finalPath := "hotel_booking/" + baseFile
			entry.Author = getAuthor(finalPath, entry.Line)
			logEntries = append(logEntries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to scan log file: %v", err), http.StatusInternalServerError)
		return
	}

	if len(logEntries) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	json.NewEncoder(w).Encode(logEntries)
}

func getAuthor(file string, line int) string {
	cmd := exec.Command("git", "blame", "-L", fmt.Sprintf("%d,%d", line, line), file)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to execute git blame: %v", err)
		return "unknown"
	}

	re := regexp.MustCompile(`\((.+?)\s+\d{4}`)
	match := re.FindStringSubmatch(string(output))
	if len(match) < 2 {
		return "unknown"
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
