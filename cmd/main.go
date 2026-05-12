package main

import (
	"flag"
	"fmt"
	"namecount/internal/handler"
	"namecount/internal/repository"
	"namecount/internal/service"
	"os"
)

func main() {
	descending := flag.Bool("desc", false, "sort by frequency descending")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: namecount [-desc] <file>")
		os.Exit(1)
	}

	path := flag.Arg(0)

	reader := repository.NewFileNameReader()
	svc := service.NewCounterService(reader)
	h := handler.NewConsoleHandler(svc, os.Stdout)

	if err := h.Handle(path, *descending); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
