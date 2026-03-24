package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Caller ejecuta peticiones HTTP con cuerpo JSON opcional contra un baseURL fijo
// (patrón típico service-to-service: API key u otras cabeceras en Header).
type Caller struct {
	HTTP    *http.Client
	BaseURL string
	Header  http.Header
}

func (c *Caller) join(path string) string {
	b := strings.TrimSuffix(strings.TrimSpace(c.BaseURL), "/")
	p := path
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return b + p
}

// DoJSON ejecuta method en path (ej. "/v1/requests"). body nil para GET/DELETE sin cuerpo.
// Devuelve status HTTP, cuerpo crudo y error de red / construcción.
func (c *Caller) DoJSON(ctx context.Context, method, path string, body any) (status int, raw []byte, err error) {
	var rdr io.Reader
	if body != nil {
		b, mErr := json.Marshal(body)
		if mErr != nil {
			return 0, nil, fmt.Errorf("marshal body: %w", mErr)
		}
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.join(path), rdr)
	if err != nil {
		return 0, nil, fmt.Errorf("build request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.Header != nil {
		for k, vv := range c.Header {
			for _, v := range vv {
				req.Header.Add(k, v)
			}
		}
	}
	hc := c.HTTP
	if hc == nil {
		hc = http.DefaultClient
	}
	resp, err := hc.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()
	raw, err = io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("read body: %w", err)
	}
	return resp.StatusCode, raw, nil
}
