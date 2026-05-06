package governanceclient

import "encoding/json"

// Decision son los valores que Nexus retorna en el campo `decision` de
// SubmitResponse / SimulateResponse. Constantes públicas para que los
// consumers no hardcodeen strings.
const (
	DecisionAllow           = "allow"
	DecisionDeny            = "deny"
	DecisionRequireApproval = "require_approval"
)

// Status son los valores que Nexus retorna en el campo `status` del request.
// Reflejan el estado del lifecycle: Allowed/Denied finalizan; PendingApproval
// queda esperando humano. Constantes para evitar string-typing.
const (
	StatusAllowed         = "allowed"
	StatusDenied          = "denied"
	StatusPendingApproval = "pending_approval"
	StatusExecuted        = "executed"
	StatusFailed          = "failed"
	StatusExpired         = "expired"
)

// PolicyMode determina si una policy actúa o solo observa. Wire format
// alineado con Nexus (`enforced` / `shadow`).
const (
	PolicyModeEnforced = "enforced"
	PolicyModeShadow   = "shadow"
)

// PolicyEffect son los valores válidos para CreatePolicyRequest.Effect.
const (
	PolicyEffectAllow           = "allow"
	PolicyEffectDeny            = "deny"
	PolicyEffectRequireApproval = "require_approval"
)

// SimulateRequestBody es el body de POST /v1/requests/simulate. Idéntico al
// de Submit (alias) — Nexus expone los mismos campos pero no persiste nada.
// Hot-path para callers que necesitan decisión síncrona sin auditoría.
type SimulateRequestBody = SubmitRequestBody

// SimulateResponse es la respuesta de POST /v1/requests/simulate.
//
// RiskAssessment se deja como json.RawMessage para no acoplar el cliente al
// shape interno de risk de Nexus; los consumers pueden Unmarshal-earlo si
// necesitan inspeccionar factores específicos.
type SimulateResponse struct {
	Decision             string          `json:"decision"`
	RiskLevel            string          `json:"risk_level"`
	DecisionReason       string          `json:"decision_reason"`
	Status               string          `json:"status"`
	PolicyMatched        *string         `json:"policy_matched,omitempty"`
	RiskAssessment       json.RawMessage `json:"risk_assessment,omitempty"`
	WouldRequireApproval bool            `json:"would_require_approval"`
	AISummary            string          `json:"ai_summary,omitempty"`
}

// CreatePolicyRequest es el body de POST /v1/policies.
type CreatePolicyRequest struct {
	Name         string  `json:"name"`
	Description  string  `json:"description,omitempty"`
	ActionType   *string `json:"action_type,omitempty"`
	TargetSystem *string `json:"target_system,omitempty"`
	Expression   string  `json:"expression"`
	Effect       string  `json:"effect"`
	RiskOverride *string `json:"risk_override,omitempty"`
	Priority     int     `json:"priority"`
	Mode         string  `json:"mode,omitempty"`
	Enabled      bool    `json:"enabled"`
}

// UpdatePolicyRequest es el body de PATCH /v1/policies/{id}. Todos los
// campos son punteros para distinguir "no tocar" de "set a zero value".
type UpdatePolicyRequest struct {
	Name         *string `json:"name,omitempty"`
	Description  *string `json:"description,omitempty"`
	ActionType   *string `json:"action_type,omitempty"`
	TargetSystem *string `json:"target_system,omitempty"`
	Expression   *string `json:"expression,omitempty"`
	Effect       *string `json:"effect,omitempty"`
	RiskOverride *string `json:"risk_override,omitempty"`
	Priority     *int    `json:"priority,omitempty"`
	Mode         *string `json:"mode,omitempty"`
	Enabled      *bool   `json:"enabled,omitempty"`
}

// PolicyResponse es el shape que retorna GET /v1/policies y los CRUD.
type PolicyResponse struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description,omitempty"`
	ActionType   *string `json:"action_type,omitempty"`
	TargetSystem *string `json:"target_system,omitempty"`
	Expression   string  `json:"expression"`
	Effect       string  `json:"effect"`
	RiskOverride *string `json:"risk_override,omitempty"`
	Priority     int     `json:"priority"`
	Origin       string  `json:"origin"`
	Mode         string  `json:"mode"`
	Enabled      bool    `json:"enabled"`
	ShadowHits   int     `json:"shadow_hits"`
	ArchivedAt   *string `json:"archived_at,omitempty"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}
