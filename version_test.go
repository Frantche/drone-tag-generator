package main

import "testing"

func TestBumpVersion_Master(t *testing.T) {
	got, err := bumpVersion("1.2.3", true, false, "main", "42")
	if err != nil || got != "1.2.4" {
		t.Errorf("expected 1.2.4, got %q, err=%v", got, err)
	}

	got, err = bumpVersion("1.2.3", true, true, "main", "42")
	if err != nil || got != "1.2.3" {
		t.Errorf("expected 1.2.3, got %q, err=%v", got, err)
	}
}

func TestBumpVersion_OtherBranch(t *testing.T) {
	got, err := bumpVersion("1.2.3", false, false, "feature/xyz", "7")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "1.2.3-feature-xyz.7"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestSanitizePrerelease(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"feature/X", "feature-X"},
		{"BUG_fix@", "BUG-fix"},
		{"???", "unknown"},
		{"clean-name", "clean-name"},
	}
	for _, c := range cases {
		if got := sanitizePrerelease(c.in); got != c.want {
			t.Errorf("sanitizePrerelease(%q) = %q; want %q", c.in, got, c.want)
		}
	}
}
