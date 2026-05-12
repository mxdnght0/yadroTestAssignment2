package handler_test

import (
	"bytes"
	"fmt"
	"namecount/internal/handler"
	"namecount/internal/service"
	"strings"
	"testing"
)

type mockService struct {
	entries []service.NameEntry
	err     error
}

func (m *mockService) Count(_ string, _ bool) ([]service.NameEntry, error) {
	return m.entries, m.err
}

func TestHandle_Output(t *testing.T) {
	svc := &mockService{
		entries: []service.NameEntry{
			{Name: "Алёна", Count: 2},
			{Name: "Миша", Count: 1},
		},
	}

	var buf bytes.Buffer
	h := handler.NewConsoleHandler(svc, &buf)

	if err := h.Handle("any", false); err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	if !strings.Contains(output, "Алёна:2") {
		t.Errorf("expected Алёна:2 in output, got: %s", output)
	}
	if !strings.Contains(output, "Миша:1") {
		t.Errorf("expected Миша:1 in output, got: %s", output)
	}
}

func TestHandle_ServiceError(t *testing.T) {
	svc := &mockService{err: fmt.Errorf("some error")}
	var buf bytes.Buffer
	h := handler.NewConsoleHandler(svc, &buf)

	if err := h.Handle("any", false); err == nil {
		t.Fatal("expected error, got nil")
	}
}
