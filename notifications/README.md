# notifications

Capacidad reusable para envío de notificaciones por email y adapters de transporte.

Implementación actual: `notifications/go/`

## Pertenece

- contratos de email sender
- adaptadores `noop`, `smtp` y `ses`
- config reusable por env
- bootstrap reusable de senders

## No pertenece

- preferencias de notificación de una app
- plantillas/copy de producto
- colas, workers o handlers específicos de un producto

## Fuente inicial esperada

- `/home/pablo/Projects/Pablo/core/saas/go/notifications` como puerto de dominio, no como implementación
- `/home/pablo/Projects/Pablo/pymes/pymes-core/backend/internal/notifications/sender.go`
