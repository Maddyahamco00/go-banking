package ledger

import "fmt"

type ReferenceStore interface {
	MarkReferenceUsed(reference string) error
}

type inMemoryReferenceStore struct {
	used map[string]struct{}
}

func NewInMemoryReferenceStore() ReferenceStore {
	return &inMemoryReferenceStore{used: map[string]struct{}{}}
}

func (s *inMemoryReferenceStore) MarkReferenceUsed(reference string) error {
	if reference == "" {
		return fmt.Errorf("reference_required")
	}
	if _, ok := s.used[reference]; ok {
		return fmt.Errorf("reference_duplicate")
	}
	s.used[reference] = struct{}{}
	return nil
}
