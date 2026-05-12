package repository_test

import (
	"namecount/internal/repository"
	"os"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "names-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func collectNames(t *testing.T, names <-chan string, errs <-chan error) []string {
	t.Helper()
	var result []string
	for name := range names {
		result = append(result, name)
	}
	if err := <-errs; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return result
}

func TestReadNames_Normal(t *testing.T) {
	path := writeTempFile(t, "Алёна\nМиша\nДима\n")
	reader := repository.NewFileNameReader()
	names, errs := reader.ReadNames(path)
	got := collectNames(t, names, errs)

	if len(got) != 3 {
		t.Fatalf("expected 3 names, got %d", len(got))
	}
}

func TestReadNames_SkipsBlankLines(t *testing.T) {
	path := writeTempFile(t, "Алёна\n\n  \nМиша\n")
	reader := repository.NewFileNameReader()
	names, errs := reader.ReadNames(path)
	got := collectNames(t, names, errs)

	if len(got) != 2 {
		t.Fatalf("expected 2 names, got %d", len(got))
	}
}

func TestReadNames_EmptyFile(t *testing.T) {
	path := writeTempFile(t, "")
	reader := repository.NewFileNameReader()
	names, errs := reader.ReadNames(path)
	got := collectNames(t, names, errs)

	if len(got) != 0 {
		t.Fatalf("expected 0 names, got %d", len(got))
	}
}

func TestReadNames_FileNotFound(t *testing.T) {
	reader := repository.NewFileNameReader()
	names, errs := reader.ReadNames("/nonexistent/path/file.txt")

	for range names {
	}
	if err := <-errs; err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
