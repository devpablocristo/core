// Package reviewclient proporciona un cliente HTTP genérico para Nexus Review API.
// Agnóstico al producto: cualquier servicio que consuma Review puede importar este paquete.
package reviewclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/devpablocristo/core/backend/go/httpclient"
)

// Client cliente HTTP hacia Nexus Review.
type Client struct {
	caller *httpclient.Caller
}

// NewClient crea el cliente con baseURL y API key.
func NewClient(baseURL, apiKey string) *Client {
	h := make(http.Header)
	h.Set("X-API-Key", apiKey)
	return &Client{
		caller: &httpclient.Caller{
			BaseURL:     baseURL,
			Header:      h,
			HTTP:        &http.Client{Timeout: 30 * time.Second},
			MaxBodySize: 1 << 20, // 1MB
		},
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
	var opts []httpclient.RequestOption
	if idempotencyKey != "" {
		opts = append(opts, httpclient.WithIdempotencyKey(idempotencyKey))
	}

	var out SubmitResponse
	st, raw, err := c.caller.DoJSON(ctx, http.MethodPost, "/v1/requests", body, opts...)
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
	st, raw, err := c.caller.DoJSON(ctx, http.MethodGet, "/v1/requests/"+id, nil)
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
	body := map[string]any{"success": success, "duration_ms": durationMS, "details": details}
	st, raw, err := c.caller.DoJSON(ctx, http.MethodPost, "/v1/requests/"+requestID+"/result", body)
	if err != nil {
		return fmt.Errorf("review report result: %w", err)
	}
	if st != http.StatusOK && st != http.StatusNoContent {
		return fmt.Errorf("review report result: status %d body %s", st, string(raw))
	}
	return nil
}

// --- Policies ---

func (c *Client) ListPolicies(ctx context.Context) (int, []byte, error) {
	return c.caller.DoJSON(ctx, http.MethodGet, "/v1/policies", nil)
}

func (c *Client) CreatePolicy(ctx context.Context, body any) (int, []byte, error) {
	return c.caller.DoJSON(ctx, http.MethodPost, "/v1/policies", body)
}

func (c *Client) UpdatePolicy(ctx context.Context, id string, body any) (int, []byte, error) {
	return c.caller.DoJSON(ctx, http.MethodPatch, "/v1/policies/"+id, body)
}

func (c *Client) DeletePolicy(ctx context.Context, id string) (int, error) {
	st, _, err := c.caller.DoJSON(ctx, http.MethodDelete, "/v1/policies/"+id, nil)
	return st, err
}

// --- Action Types ---

func (c *Client) ListActionTypes(ctx context.Context) (int, []byte, error) {
	return c.caller.DoJSON(ctx, http.MethodGet, "/v1/action-types", nil)
}

// --- Approvals ---

func (c *Client) ListPendingApprovals(ctx context.Context) (int, []byte, error) {
	return c.caller.DoJSON(ctx, http.MethodGet, "/v1/approvals/pending", nil)
}

func (c *Client) Approve(ctx context.Context, id string, body any) (int, []byte, error) {
	return c.caller.DoJSON(ctx, http.MethodPost, "/v1/approvals/"+id+"/approve", body)
}

func (c *Client) Reject(ctx context.Context, id string, body any) (int, []byte, error) {
	return c.caller.DoJSON(ctx, http.MethodPost, "/v1/approvals/"+id+"/reject", body)
}

// --- Helpers ---

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
