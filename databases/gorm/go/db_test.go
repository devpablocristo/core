package gormdb

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	t.Parallel()
	cfg := DefaultConfig()
	if cfg.MaxOpenConns != 10 {
		t.Errorf("expected MaxOpenConns=10, got %d", cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns != 5 {
		t.Errorf("expected MaxIdleConns=5, got %d", cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime != time.Hour {
		t.Errorf("expected ConnMaxLifetime=1h, got %v", cfg.ConnMaxLifetime)
	}
}

func TestOpenEmptyURL(t *testing.T) {
	t.Parallel()
	_, err := Open("", DefaultConfig())
	if err == nil {
		t.Fatal("expected error on empty URL")
	}
}

func TestPingNilDB(t *testing.T) {
	t.Parallel()
	err := Ping(nil, nil)
	if err == nil {
		t.Fatal("expected error on nil DB")
	}
}

func TestCloseNilDB(t *testing.T) {
	t.Parallel()
	err := Close(nil)
	if err != nil {
		t.Fatalf("expected no error on nil DB, got %v", err)
	}
}
