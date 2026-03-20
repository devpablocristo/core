package evidence

import "testing"

func TestHMACSignerSignsPack(t *testing.T) {
	t.Parallel()

	signer, err := NewHMACSigner("secret", "kid-1")
	if err != nil {
		t.Fatalf("NewHMACSigner returned error: %v", err)
	}
	pack := Pack{}
	if err := signer.Sign(&pack); err != nil {
		t.Fatalf("Sign returned error: %v", err)
	}
	if pack.Signature == "" {
		t.Fatal("expected signature")
	}
}
