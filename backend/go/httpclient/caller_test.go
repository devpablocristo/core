package httpclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCaller_DoJSON_Get(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/ping" {
			t.Fatalf("path %s", r.URL.Path)
		}
		if r.Header.Get("X-API-Key") != "k" {
			t.Fatalf("missing api key")
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"ok": "1"})
	}))
	t.Cleanup(srv.Close)

	h := make(http.Header)
	h.Set("X-API-Key", "k")
	c := Caller{HTTP: srv.Client(), BaseURL: srv.URL, Header: h}
	st, raw, err := c.DoJSON(context.Background(), http.MethodGet, "/v1/ping", nil)
	if err != nil {
		t.Fatal(err)
	}
	if st != http.StatusOK {
		t.Fatalf("status %d", st)
	}
	var m map[string]string
	if err := json.Unmarshal(raw, &m); err != nil || m["ok"] != "1" {
		t.Fatalf("body %s err %v", raw, err)
	}
}

func TestCaller_DoJSON_Post(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var in map[string]int
		_ = json.NewDecoder(r.Body).Decode(&in)
		if in["n"] != 42 {
			t.Fatalf("body %+v", in)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"x"}`))
	}))
	t.Cleanup(srv.Close)

	c := Caller{HTTP: srv.Client(), BaseURL: srv.URL}
	st, raw, err := c.DoJSON(context.Background(), http.MethodPost, "/v1/x", map[string]int{"n": 42})
	if err != nil {
		t.Fatal(err)
	}
	if st != http.StatusCreated {
		t.Fatalf("status %d", st)
	}
	if string(raw) != `{"id":"x"}` {
		t.Fatalf("raw %s", raw)
	}
}
