package requestid

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	id1 := New()
	id2 := New()

	if id1 == "" {
		t.Error("New() returned empty string")
	}

	if id1 == id2 {
		t.Errorf("New() returned duplicate IDs: %s", id1)
	}
}

func TestNewMonotonic(t *testing.T) {
	ids := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		id := New()
		if ids[id] {
			t.Fatalf("duplicate ID at iteration %d: %s", i, id)
		}
		ids[id] = true
	}
}

func TestWithContextAndFromContext(t *testing.T) {
	ctx := context.Background()
	id := New()

	ctx = WithContext(ctx, id)
	got := FromContext(ctx)

	if got != id {
		t.Errorf("FromContext() = %q, want %q", got, id)
	}
}

func TestFromContextEmpty(t *testing.T) {
	ctx := context.Background()
	got := FromContext(ctx)

	if got != "" {
		t.Errorf("FromContext() on empty context = %q, want empty", got)
	}
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New()
	}
}
