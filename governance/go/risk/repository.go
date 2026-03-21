package risk

import "context"

// Repository define el puerto para cargar historia operativa del sujeto evaluado.
type Repository interface {
	GetHistory(context.Context, string) (History, error)
}
