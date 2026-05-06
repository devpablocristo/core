package governanceclient

import (
	"github.com/devpablocristo/core/http/go/httpclient"
)

// RequestOption modifica una llamada individual sin tocar la config global
// del Client. Es la forma de pasar headers per-call (X-Org-ID,
// Idempotency-Key) sin acoplar al consumer al package httpclient.
type RequestOption func(*requestOptions)

type requestOptions struct {
	orgID          string
	idempotencyKey string
}

// WithOrgID setea el header X-Org-ID en la llamada. Sirve para callers
// multi-tenant (ej: Pymes service que opera para varios tenants con una
// única API key con scope nexus:cross_org). Sin esta opción, Nexus usa el
// org del principal de la API key.
func WithOrgID(orgID string) RequestOption {
	return func(o *requestOptions) { o.orgID = orgID }
}

// WithIdempotencyKey setea el header Idempotency-Key en la llamada.
func WithIdempotencyKey(key string) RequestOption {
	return func(o *requestOptions) { o.idempotencyKey = key }
}

func collectOptions(opts ...RequestOption) []httpclient.RequestOption {
	var ro requestOptions
	for _, o := range opts {
		if o != nil {
			o(&ro)
		}
	}
	var out []httpclient.RequestOption
	if ro.orgID != "" {
		out = append(out, httpclient.WithHeader("X-Org-ID", ro.orgID))
	}
	if ro.idempotencyKey != "" {
		out = append(out, httpclient.WithIdempotencyKey(ro.idempotencyKey))
	}
	return out
}
