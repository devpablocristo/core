package delegations

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	delegationdomain "github.com/devpablocristo/core/governance/go/delegations/usecases/domain"
)

type UseCases struct {
	repo Repository
	now  func() time.Time
}

func NewUseCases(repo Repository) *UseCases {
	return &UseCases{
		repo: repo,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
}

func (u *UseCases) Create(ctx context.Context, item delegationdomain.Delegation) (delegationdomain.Delegation, error) {
	if err := validate(item); err != nil {
		return delegationdomain.Delegation{}, err
	}
	if item.MaxRiskClass == "" {
		item.MaxRiskClass = delegationdomain.RiskClassHigh
	}
	item.Enabled = true
	if item.CreatedAt.IsZero() {
		item.CreatedAt = u.now()
	}
	item.UpdatedAt = u.now()
	item.AllowedActionTypes = normalizeActions(item.AllowedActionTypes)
	return u.repo.Create(ctx, item)
}

func (u *UseCases) Update(ctx context.Context, item delegationdomain.Delegation) (delegationdomain.Delegation, error) {
	if strings.TrimSpace(item.ID) == "" {
		return delegationdomain.Delegation{}, fmt.Errorf("delegation id is required")
	}
	if err := validate(item); err != nil {
		return delegationdomain.Delegation{}, err
	}
	if item.MaxRiskClass == "" {
		item.MaxRiskClass = delegationdomain.RiskClassHigh
	}
	item.AllowedActionTypes = normalizeActions(item.AllowedActionTypes)
	item.UpdatedAt = u.now()
	return u.repo.Update(ctx, item)
}

func (u *UseCases) GetByID(ctx context.Context, id string) (delegationdomain.Delegation, error) {
	return u.repo.GetByID(ctx, strings.TrimSpace(id))
}

func (u *UseCases) List(ctx context.Context) ([]delegationdomain.Delegation, error) {
	return u.repo.List(ctx)
}

func (u *UseCases) ListByAgentID(ctx context.Context, agentID string) ([]delegationdomain.Delegation, error) {
	return u.repo.ListByAgentID(ctx, strings.TrimSpace(agentID))
}

func (u *UseCases) DeleteByID(ctx context.Context, id string) error {
	return u.repo.DeleteByID(ctx, strings.TrimSpace(id))
}

func (u *UseCases) Check(ctx context.Context, agentID, actionType string, requestedRisk delegationdomain.RiskClass) (bool, delegationdomain.Delegation, error) {
	delegations, err := u.repo.ListByAgentID(ctx, strings.TrimSpace(agentID))
	if err != nil {
		return false, delegationdomain.Delegation{}, fmt.Errorf("list delegations: %w", err)
	}
	actionType = strings.TrimSpace(actionType)
	if requestedRisk == "" {
		requestedRisk = delegationdomain.RiskClassLow
	}

	for _, item := range delegations {
		if !item.Enabled {
			continue
		}
		if !matchesAction(item, actionType) {
			continue
		}
		if riskRank(requestedRisk) > riskRank(item.MaxRiskClass) {
			continue
		}
		return true, item, nil
	}

	return false, delegationdomain.Delegation{}, nil
}

func validate(item delegationdomain.Delegation) error {
	if strings.TrimSpace(item.OwnerID) == "" {
		return fmt.Errorf("owner id is required")
	}
	if strings.TrimSpace(item.AgentID) == "" {
		return fmt.Errorf("agent id is required")
	}
	if item.MaxRiskClass != "" && riskRank(item.MaxRiskClass) == 0 {
		return fmt.Errorf("invalid max risk class")
	}
	return nil
}

func normalizeActions(values []string) []string {
	items := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if !slices.Contains(items, value) {
			items = append(items, value)
		}
	}
	return items
}

func matchesAction(item delegationdomain.Delegation, actionType string) bool {
	if len(item.AllowedActionTypes) == 0 {
		return true
	}
	return slices.Contains(item.AllowedActionTypes, actionType)
}

func riskRank(value delegationdomain.RiskClass) int {
	switch value {
	case delegationdomain.RiskClassLow:
		return 1
	case delegationdomain.RiskClassMedium:
		return 2
	case delegationdomain.RiskClassHigh:
		return 3
	default:
		return 0
	}
}
