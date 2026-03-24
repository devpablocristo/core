// Package reviewclient proporciona un cliente HTTP genérico para Nexus Review API.
// Agnóstico al producto: cualquier servicio que consuma Review (Companion, Pymes, Ponti, etc.)
// puede importar este paquete en lugar de mantener su propia copia.
package reviewclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client cliente HTTP hacia Nexus Review.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient crea el cliente con baseURL (sin slash final) y API key.
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// --- DTOs alineados con Review API ---

// SubmitRequestBody cuerpo de POST /v1/requests.
type SubmitRequestBody struct {
	RequesterType  string         `json:"requester_type"`
	RequesterID    string         `json:"requester_id"`
	RequesterName  string         `json:"requester_name,omitempty"`
	ActionType     string         `json:"action_type"`
	TargetSystem   string         `json:"target_system,omitempty"`
	TargetResource string         `json:"target_resource,omitempty"`
	Params         map[string]any `json:"params,omitempty"`
	Reason         string         `json:"reason,omitempty"`
	Context        string         `json:"context,omitempty"`
}

// SubmitResponse respuesta de POST /v1/requests.
type SubmitResponse struct {
	RequestID      string `json:"request_id"`
	Decision       string `json:"decision"`
	RiskLevel      string `json:"risk_level"`
	DecisionReason string `json:"decision_reason"`
	Status         string `json:"status"`
}

// RequestSummary respuesta de GET /v1/requests/{id}.
type RequestSummary struct {
	ID             string `json:"id"`
	RequesterType  string `json:"requester_type"`
	RequesterID    string `json:"requester_id"`
	ActionType     string `json:"action_type"`
	TargetSystem   string `json:"target_system"`
	TargetResource string `json:"target_resource"`
	Reason         string `json:"reason"`
	RiskLevel      string `json:"risk_level"`
	Decision       string `json:"decision"`
	DecisionReason string `json:"decision_reason"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// --- Requests ---

// SubmitRequest envía POST /v1/requests con idempotency key opcional.
func (c *Client) SubmitRequest(ctx context.Context, idempotencyKey string, body SubmitRequestBody) (SubmitResponse, error) {
	var out SubmitResponse
	extra := http.Header{}
	if idempotencyKey != "" {
		extra.Set("Idempotency-Key", idempotencyKey)
	}
	st, raw, err := c.doJSON(ctx, http.MethodPost, "/v1/requests", body, extra)
	if err != nil {
		return out, fmt.Errorf("review submit: %w", err)
	}
	if st != http.StatusCreated {
		return out, fmt.Errorf("review submit: status %d body %s", st, string(raw))
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return out, fmt.Errorf("decode submit response: %w", err)
	}
	return out, nil
}

// GetRequest consulta GET /v1/requests/{id}. Devuelve status HTTP para distinguir 404.
func (c *Client) GetRequest(ctx context.Context, id string) (RequestSummary, int, error) {
	var out RequestSummary
	st, raw, err := c.doJSON(ctx, http.MethodGet, "/v1/requests/"+id, nil, nil)
	if err != nil {
		return out, 0, fmt.Errorf("review get request: %w", err)
	}
	if st == http.StatusNotFound {
		return out, st, nil
	}
	if st != http.StatusOK {
		return out, st, fmt.Errorf("review get request: status %d body %s", st, string(raw))
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return out, st, fmt.Errorf("decode get response: %w", err)
	}
	return out, st, nil
}

// ReportResult reporta resultado de ejecución POST /v1/requests/{id}/result.
func (c *Client) ReportResult(ctx context.Context, requestID string, success bool, durationMS int64, details string) error {
	body := map[string]any{
		"success":     success,
		"duration_ms": durationMS,
		"details":     details,
	}
	st, raw, err := c.doJSON(ctx, http.MethodPost, "/v1/requests/"+requestID+"/result", body, nil)
	if err != nil {
		return fmt.Errorf("review report result: %w", err)
	}
	if st != http.StatusOK && st != http.StatusNoContent {
		return fmt.Errorf("review report result: status %d body %s", st, string(raw))
	}
	return nil
}

// --- Policies ---

// ListPolicies consulta GET /v1/policies. Devuelve status + body crudo.
func (c *Client) ListPolicies(ctx context.Context) (int, []byte, error) {
	return c.doJSON(ctx, http.MethodGet, "/v1/policies", nil, nil)
}

// CreatePolicy envía POST /v1/policies.
func (c *Client) CreatePolicy(ctx context.Context, body any) (int, []byte, error) {
	return c.doJSON(ctx, http.MethodPost, "/v1/policies", body, nil)
}

// UpdatePolicy envía PATCH /v1/policies/{id}.
func (c *Client) UpdatePolicy(ctx context.Context, id string, body any) (int, []byte, error) {
	return c.doJSON(ctx, http.MethodPatch, "/v1/policies/"+id, body, nil)
}

// DeletePolicy envía DELETE /v1/policies/{id}.
func (c *Client) DeletePolicy(ctx context.Context, id string) (int, error) {
	st, _, err := c.doJSON(ctx, http.MethodDelete, "/v1/policies/"+id, nil, nil)
	return st, err
}

// --- Action Types ---

// ListActionTypes consulta GET /v1/action-types.
func (c *Client) ListActionTypes(ctx context.Context) (int, []byte, error) {
	return c.doJSON(ctx, http.MethodGet, "/v1/action-types", nil, nil)
}

// --- Approvals ---

// ListPendingApprovals consulta GET /v1/approvals/pending.
func (c *Client) ListPendingApprovals(ctx context.Context) (int, []byte, error) {
	return c.doJSON(ctx, http.MethodGet, "/v1/approvals/pending", nil, nil)
}

// Approve envía POST /v1/approvals/{id}/approve.
func (c *Client) Approve(ctx context.Context, id string, body any) (int, []byte, error) {
	return c.doJSON(ctx, http.MethodPost, "/v1/approvals/"+id+"/approve", body, nil)
}

// Reject envía POST /v1/approvals/{id}/reject.
func (c *Client) Reject(ctx context.Context, id string, body any) (int, []byte, error) {
	return c.doJSON(ctx, http.MethodPost, "/v1/approvals/"+id+"/reject", body, nil)
}

// --- Helpers ---

// doJSON ejecuta una petición HTTP JSON contra Review.
func (c *Client) doJSON(ctx context.Context, method, path string, body any, extraHeaders http.Header) (int, []byte, error) {
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return 0, nil, fmt.Errorf("marshal body: %w", err)
		}
		rdr = bytes.NewReader(b)
	}

	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, rdr)
	if err != nil {
		return 0, nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("X-API-Key", c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, vv := range extraHeaders {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("read body: %w", err)
	}
	return resp.StatusCode, raw, nil
}

// ParseErrorBody intenta extraer mensaje de error de respuesta de Review.
func ParseErrorBody(raw []byte) string {
	var eb struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if json.Unmarshal(raw, &eb) == nil && eb.Message != "" {
		return eb.Message
	}
	return string(raw)
}
