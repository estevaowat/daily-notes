package notes

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCreateNewNote(t *testing.T) {
	tmpDir := t.TempDir()
	SetVaultPath(tmpDir)
	defer SetVaultPath(DefaultVaultPath)

	now := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	err := createDailyNoteAt(now, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedFile := "2026-03-15.md"
	filePath := filepath.Join(tmpDir, DaysFolder, expectedFile)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("expected file %s to be created", expectedFile)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if len(content) != 0 {
		t.Errorf("expected empty file, got: %s", content)
	}
}

func TestFileAlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()
	SetVaultPath(tmpDir)
	defer SetVaultPath(DefaultVaultPath)

	now := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
	filename := "2026-03-10.md"
	daysPath := filepath.Join(tmpDir, DaysFolder)
	os.MkdirAll(daysPath, 0755)
	os.WriteFile(filepath.Join(daysPath, filename), []byte("# 2026-03-10\n"), 0644)

	err := createDailyNoteAt(now, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreatesDaysFolder(t *testing.T) {
	tmpDir := t.TempDir()
	SetVaultPath(tmpDir)
	defer SetVaultPath(DefaultVaultPath)

	now := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	err := createDailyNoteAt(now, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	daysPath := filepath.Join(tmpDir, DaysFolder)
	info, err := os.Stat(daysPath)
	if os.IsNotExist(err) {
		t.Error("expected days folder to be created")
	} else if err != nil {
		t.Fatalf("failed to stat days folder: %v", err)
	} else if !info.IsDir() {
		t.Error("expected days to be a directory")
	}
}

func TestExtractIncompleteTasks(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "empty content",
			content:  "",
			expected: []string{},
		},
		{
			name:     "no tasks",
			content:  "# Today's Notes\n\nSome regular text.",
			expected: []string{},
		},
		{
			name:     "only completed tasks",
			content:  "- [x] Done task 1\n- [x] Done task 2",
			expected: []string{},
		},
		{
			name:     "only incomplete tasks",
			content:  "- [ ] Incomplete task 1\n- [ ] Incomplete task 2",
			expected: []string{"- [ ] Incomplete task 1", "- [ ] Incomplete task 2"},
		},
		{
			name:     "mixed tasks",
			content:  "- [ ] Todo task\n- [x] Done task\n- [ ] Another todo\n- [-] In progress",
			expected: []string{"- [ ] Todo task", "- [ ] Another todo"},
		},
		{
			name:     "task with extra spaces after checkbox",
			content:  "- [ ]  Task with extra spaces",
			expected: []string{"- [ ]  Task with extra spaces"},
		},
		{
			name:    "real world content",
			content: "# 2026-03-03\n\n- [ ] criar um notion para a viagem do japão\n- [ ] criar um aplicativo\n- [x] Already done",
			expected: []string{
				"- [ ] criar um notion para a viagem do japão",
				"- [ ] criar um aplicativo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractIncompleteTasks(tt.content)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d tasks, got %d", len(tt.expected), len(result))
				return
			}
			for i, task := range tt.expected {
				if result[i] != task {
					t.Errorf("expected task %d to be %q, got %q", i, task, result[i])
				}
			}
		})
	}
}

func TestFindPreviousNote(t *testing.T) {
	tmpDir := t.TempDir()
	SetVaultPath(tmpDir)
	defer SetVaultPath(DefaultVaultPath)

	daysPath := filepath.Join(tmpDir, DaysFolder)
	os.MkdirAll(daysPath, 0755)

	os.WriteFile(filepath.Join(daysPath, "2026-03-01.md"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(daysPath, "2026-03-05.md"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(daysPath, "2026-03-08.md"), []byte("content"), 0644)

	path, date, err := findPreviousNoteAt(time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.HasSuffix(path, "2026-03-08.md") {
		t.Errorf("expected most recent note 2026-03-08.md, got %s", path)
	}

	expectedDate := time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC)
	if !date.Equal(expectedDate) {
		t.Errorf("expected date %v, got %v", expectedDate, date)
	}
}

func TestFindPreviousNoteNoNotes(t *testing.T) {
	tmpDir := t.TempDir()
	SetVaultPath(tmpDir)
	defer SetVaultPath(DefaultVaultPath)

	daysPath := filepath.Join(tmpDir, DaysFolder)
	os.MkdirAll(daysPath, 0755)

	path, _, err := findPreviousNoteAt(time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if path != "" {
		t.Errorf("expected empty path, got %s", path)
	}
}

func TestFindPreviousNoteGap(t *testing.T) {
	tmpDir := t.TempDir()
	SetVaultPath(tmpDir)
	defer SetVaultPath(DefaultVaultPath)

	daysPath := filepath.Join(tmpDir, DaysFolder)
	os.MkdirAll(daysPath, 0755)

	os.WriteFile(filepath.Join(daysPath, "2026-02-28.md"), []byte("content"), 0644)

	path, date, err := findPreviousNoteAt(time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedDate := time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC)
	if !date.Equal(expectedDate) {
		t.Errorf("expected date %v, got %v", expectedDate, date)
	}

	_ = path
}
