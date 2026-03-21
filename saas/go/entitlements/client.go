package entitlements

import (
	"log/slog"

	"github.com/devpablocristo/core/saas/go/entitlement"
)

type Config = entitlement.ClientConfig
type Client = entitlement.Client

func NewEntitlementsClient(cfg Config, logger *slog.Logger) *Client {
	return entitlement.NewClient(cfg, logger)
}
