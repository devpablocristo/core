package gormdb

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// HandleNotFound convierte gorm.ErrRecordNotFound a un error con formato estándar.
// Para otros errores retorna un error genérico sin exponer detalles internos.
func HandleNotFound(err error, entity string, id any) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("%s %v not found", entity, id)
	}
	return fmt.Errorf("failed to get %s", entity)
}

// IsNotFound retorna true si el error es gorm.ErrRecordNotFound.
func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

type txContextKey struct{}

// WithTx adjunta una transacción GORM al contexto para que los repositorios la reutilicen.
func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txContextKey{}, tx)
}

// TxFromContext extrae la transacción del contexto si existe.
func TxFromContext(ctx context.Context) *gorm.DB {
	tx, _ := ctx.Value(txContextKey{}).(*gorm.DB)
	return tx
}

// DBOrTx retorna la transacción del contexto si existe, o el db base.
func DBOrTx(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx := TxFromContext(ctx); tx != nil {
		return tx
	}
	return db
}
