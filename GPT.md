# GPT.md

Este repo usa `AGENTS.md` como fuente de verdad para GPT/Codex y agentes equivalentes.

Aplican sin excepción las reglas de `AGENTS.md`.

## Resumen crítico

- Este repo es un monorepo de capacidades reutilizables, no de productos
- Los módulos raíz son: `saas`, `backend`, `postgres`, `serverless`, `governance`, `artifact`, `ai`
- Dentro de este repo NO se usan sufijos `-core`
- No crear `common/`, `shared/`, `utils/` en la raíz
- Lenguaje distinto no implica repo distinto
- Solo crear subcarpetas `go/`, `rust/`, `python/` si la misma capacidad tiene múltiples implementaciones
- No meter dominio específico de Nexus, Pymes, Ponti, Fixguard o KMA/AlphaCoding
- Hexagonal se evalúa por dirección de dependencias y aislamiento del dominio, no por el nombre de las carpetas
- La estructura `usecases/domain`, `handler/dto` y `repository/models` es convención obligatoria para módulos Go con dominio real
- Rastrear siempre consumidores y productores afectados antes de afirmar que un cambio es seguro
- No decir "listo" sin verificación
