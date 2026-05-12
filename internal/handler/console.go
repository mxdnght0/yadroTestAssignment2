package handler

import (
	"fmt"
	"io"
	"namecount/internal/service"
)

type ConsoleHandler struct {
	service service.CounterService
	output  io.Writer
}

func NewConsoleHandler(svc service.CounterService, output io.Writer) *ConsoleHandler {
	return &ConsoleHandler{service: svc, output: output}
}

func (h *ConsoleHandler) Handle(path string, descending bool) error {
	entries, err := h.service.Count(path, descending)
	if err != nil {
		return fmt.Errorf("counting names: %w", err)
	}

	for _, e := range entries {
		fmt.Fprintf(h.output, "%s:%d\n", e.Name, e.Count)
	}

	return nil
}
