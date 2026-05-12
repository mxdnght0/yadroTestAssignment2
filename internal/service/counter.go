package service

import (
	"fmt"
	"namecount/internal/repository"
	"sort"
)

type NameEntry struct {
	Name  string
	Count int
}

type CounterService interface {
	Count(path string, descending bool) ([]NameEntry, error)
}

type counterService struct {
	reader repository.NameReader
}

func NewCounterService(reader repository.NameReader) CounterService {
	return &counterService{reader: reader}
}

func (s *counterService) Count(path string, descending bool) ([]NameEntry, error) {
	names, errs := s.reader.ReadNames(path)

	freq := make(map[string]int)
	for name := range names {
		freq[name]++
	}

	if err := <-errs; err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	entries := make([]NameEntry, 0, len(freq))
	for name, count := range freq {
		entries = append(entries, NameEntry{Name: name, Count: count})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			if descending {
				return entries[i].Count > entries[j].Count
			}
			return entries[i].Count < entries[j].Count
		}
		return entries[i].Name < entries[j].Name
	})

	return entries, nil
}
