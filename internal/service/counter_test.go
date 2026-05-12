package service_test

import (
	"fmt"
	"namecount/internal/service"
	"testing"
)

type mockReader struct {
	lines []string
	err   error
}

func (m *mockReader) ReadNames(_ string) (<-chan string, <-chan error) {
	names := make(chan string, len(m.lines))
	errs := make(chan error, 1)

	for _, l := range m.lines {
		names <- l
	}
	close(names)

	if m.err != nil {
		errs <- m.err
	}
	close(errs)

	return names, errs
}

func TestCount_FrequencyIsCorrect(t *testing.T) {
	reader := &mockReader{lines: []string{"Алёна", "Миша", "Алёна", "Дима"}}
	svc := service.NewCounterService(reader)

	entries, err := svc.Count("any", false)
	if err != nil {
		t.Fatal(err)
	}

	freq := make(map[string]int, len(entries))
	for _, e := range entries {
		freq[e.Name] = e.Count
	}

	if freq["Алёна"] != 2 {
		t.Errorf("Алёна: want 2, got %d", freq["Алёна"])
	}
	if freq["Миша"] != 1 {
		t.Errorf("Миша: want 1, got %d", freq["Миша"])
	}
	if freq["Дима"] != 1 {
		t.Errorf("Дима: want 1, got %d", freq["Дима"])
	}
}

func TestCount_SortAscending(t *testing.T) {
	reader := &mockReader{lines: []string{"Алёна", "Миша", "Алёна", "Алёна"}}
	svc := service.NewCounterService(reader)

	entries, err := svc.Count("any", false)
	if err != nil {
		t.Fatal(err)
	}

	if entries[0].Count > entries[len(entries)-1].Count {
		t.Error("expected ascending order")
	}
}

func TestCount_SortDescending(t *testing.T) {
	reader := &mockReader{lines: []string{"Алёна", "Миша", "Алёна", "Алёна"}}
	svc := service.NewCounterService(reader)

	entries, err := svc.Count("any", true)
	if err != nil {
		t.Fatal(err)
	}

	if entries[0].Count < entries[len(entries)-1].Count {
		t.Error("expected descending order")
	}
}

func TestCount_EmptyInput(t *testing.T) {
	reader := &mockReader{lines: []string{}}
	svc := service.NewCounterService(reader)

	entries, err := svc.Count("any", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestCount_ReaderError(t *testing.T) {
	reader := &mockReader{err: fmt.Errorf("disk error")}
	svc := service.NewCounterService(reader)

	_, err := svc.Count("any", false)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
