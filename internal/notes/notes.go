package notes

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	DefaultVaultPath = "/Users/estevaowatanabe/Library/Mobile Documents/iCloud~md~obsidian/Documents/second_brain"
	DaysFolder       = "days"
)

var vaultPath = DefaultVaultPath

func SetVaultPath(path string) {
	vaultPath = path
}

func CreateDailyNote() error {
	now := time.Now()

	var content string
	notePath, _, err := findPreviousNoteAt(now)
	if err != nil {
		return err
	}

	if notePath != "" {
		prevContent, err := os.ReadFile(notePath)
		if err != nil {
			return fmt.Errorf("failed to read previous note: %w", err)
		}

		tasks := ExtractIncompleteTasks(string(prevContent))
		if len(tasks) > 0 {
			var output bytes.Buffer
			output.WriteString("## Previous's Tasks\n")
			for _, task := range tasks {
				output.WriteString(task + "\n")
			}
			content = output.String()
		}
	}

	return createDailyNoteAt(now, content)
}

func createDailyNoteAt(t time.Time, content string) error {
	dateStr := t.Format("2006-01-02")
	filename := fmt.Sprintf("%s.md", dateStr)
	daysPath := filepath.Join(vaultPath, DaysFolder)

	if err := os.MkdirAll(daysPath, 0755); err != nil {
		return fmt.Errorf("failed to create days directory: %w", err)
	}

	notePath := filepath.Join(daysPath, filename)

	_, err := os.Stat(notePath)
	if err == nil {
		fmt.Printf("the file %s already exists\n", filename)
		return nil
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check file existence: %w", err)
	}

	if err := os.WriteFile(notePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}

	fmt.Printf("Created: %s\n", notePath)
	return nil
}

func FindPreviousNote() (string, time.Time, error) {
	return findPreviousNoteAt(time.Now())
}

func findPreviousNoteAt(t time.Time) (string, time.Time, error) {
	daysPath := filepath.Join(vaultPath, DaysFolder)

	for i := 1; i <= 365; i++ {
		searchDate := t.AddDate(0, 0, -i)
		dateStr := searchDate.Format("2006-01-02")
		filename := fmt.Sprintf("%s.md", dateStr)
		notePath := filepath.Join(daysPath, filename)

		if _, err := os.Stat(notePath); err == nil {
			return notePath, searchDate, nil
		}
	}

	return "", time.Time{}, nil
}

func ExtractIncompleteTasks(content string) []string {
	var tasks []string
	lines := bytes.Split([]byte(content), []byte("\n"))

	for _, line := range lines {
		text := string(bytes.TrimSpace(line))
		if len(text) >= 6 && text[:6] == "- [ ] " {
			tasks = append(tasks, string(line))
		}
	}

	return tasks
}
