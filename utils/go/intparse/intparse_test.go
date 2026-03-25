package intparse

import "testing"

func TestParsePositiveInt(t *testing.T) {
	t.Parallel()
	n, err := ParsePositiveInt("", 7)
	if err != nil || n != 7 {
		t.Fatalf("empty: got %d err=%v", n, err)
	}
	n, err = ParsePositiveInt("3", 7)
	if err != nil || n != 3 {
		t.Fatalf("3: got %d err=%v", n, err)
	}
	_, err = ParsePositiveInt("0", 7)
	if err == nil {
		t.Fatal("expected error for zero")
	}
	_, err = ParsePositiveInt("-1", 7)
	if err == nil {
		t.Fatal("expected error for negative")
	}
}
