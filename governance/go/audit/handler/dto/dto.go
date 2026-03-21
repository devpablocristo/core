package dto

import auditdomain "github.com/devpablocristo/core/governance/go/audit/usecases/domain"

type RecordRequest struct {
	Event auditdomain.RequestEvent `json:"event"`
}

type ReplayRequest struct {
	RequestID string `json:"request_id"`
}

type RecordResponse struct {
	Event auditdomain.RequestEvent `json:"event"`
}

type ReplayResponse struct {
	Replay auditdomain.ReplayOutput `json:"replay"`
}
