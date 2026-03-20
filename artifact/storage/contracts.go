package storage

import (
	"context"

	"github.com/devpablocristo/core/artifact"
)

type Store interface {
	Put(context.Context, artifact.Asset) (string, error)
}

type Getter interface {
	Get(context.Context, string) (artifact.Asset, error)
}
