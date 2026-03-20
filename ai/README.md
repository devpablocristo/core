# ai

Capacidad reusable para runtime AI, providers y orquestación.

## Pertenece

- orchestration
- provider factory
- auth/rate limit AI
- runtime helpers
- contexts
- resilience/circuit breaker
- adapters AI reutilizables

## No pertenece

- prompts específicos de negocio
- tools específicas de una app
- contratos cerrados a un solo producto

## Fuente inicial esperada

- `/home/pablo/Projects/Pablo/ai-core`

## Nota

Aunque hoy probablemente siga siendo Python, el módulo se nombra por capacidad y no por lenguaje.

## Estructura actual

El módulo ya sigue una estructura Python/FastAPI estándar:

- `src/core_ai/domain/`
- `src/core_ai/providers/`
- `src/core_ai/services/`
- `src/core_ai/registry/`
- `src/core_ai/config/`
- `src/core_ai/api/`
- wrappers de compatibilidad en `types.py`, `orchestrator.py` y `provider_factory.py`
- `tests/`
