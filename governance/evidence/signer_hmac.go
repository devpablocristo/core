package evidence

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type HMACSigner struct {
	key   []byte
	keyID string
}

func NewHMACSigner(signingKey, keyID string) (*HMACSigner, error) {
	signingKey = strings.TrimSpace(signingKey)
	if signingKey == "" {
		return nil, fmt.Errorf("signing key is required")
	}
	keyID = strings.TrimSpace(keyID)
	if keyID == "" {
		keyID = "default"
	}
	return &HMACSigner{key: []byte(signingKey), keyID: keyID}, nil
}

func (s *HMACSigner) Sign(pack *Pack) error {
	pack.Signature = ""
	payload, err := json.Marshal(pack)
	if err != nil {
		return fmt.Errorf("marshal evidence pack: %w", err)
	}
	mac := hmac.New(sha256.New, s.key)
	_, _ = mac.Write(payload)
	pack.Signature = s.keyID + ":" + time.Now().UTC().Format(time.RFC3339) + ":" + hex.EncodeToString(mac.Sum(nil))
	return nil
}
