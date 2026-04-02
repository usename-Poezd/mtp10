package store

import "testing"

func TestStoreCRUD(t *testing.T) {
	s := NewStore()

	all := s.GetAll()
	if all == nil {
		t.Fatal("expected empty slice, got nil")
	}
	if len(all) != 0 {
		t.Fatalf("expected 0 todos, got %d", len(all))
	}

	created := s.Create("Test Task")
	if created.ID == "" {
		t.Fatal("expected created todo id to be non-empty")
	}
	if created.Title != "Test Task" {
		t.Fatalf("expected title Test Task, got %s", created.Title)
	}
	if created.Completed {
		t.Fatal("expected completed=false for new todo")
	}

	got, ok := s.GetByID(created.ID)
	if !ok {
		t.Fatal("expected todo to exist")
	}
	if got.ID != created.ID {
		t.Fatalf("expected id %s, got %s", created.ID, got.ID)
	}

	updated, ok := s.Update(created.ID, "Updated", true)
	if !ok {
		t.Fatal("expected update to succeed")
	}
	if updated.Title != "Updated" || !updated.Completed {
		t.Fatalf("unexpected updated todo: %+v", updated)
	}

	if !s.Delete(created.ID) {
		t.Fatal("expected delete to succeed")
	}
	if s.Delete(created.ID) {
		t.Fatal("expected delete to fail for missing todo")
	}
}

func TestStoreNotFound(t *testing.T) {
	s := NewStore()

	if _, ok := s.GetByID("missing"); ok {
		t.Fatal("expected missing todo")
	}

	if _, ok := s.Update("missing", "Updated", true); ok {
		t.Fatal("expected update to fail for missing todo")
	}

	if s.Delete("missing") {
		t.Fatal("expected delete to fail for missing todo")
	}
}
