# scripts

Acá viven scripts transversales del monorepo `core`.

## Uso esperado

- bootstrap del workspace
- checks de calidad
- generación de código o contratos
- tareas repetitivas de mantenimiento

## Scripts actuales

- `validate-runtime-layout.sh`: valida que todo runtime viva bajo `go/`, `python/` o `rust/`
- `validate-module-versions.sh`: valida `VERSION`, semver y consistencia con manifests por módulo
- `list-module-versions.sh`: lista módulos versionados y el tag esperado de release
- `bump-module-version.sh`: sube la versión de un módulo concreto
- `test-go-modules.sh`: corre `go test ./...` en cada módulo Go independiente
- `test-rust-modules.sh`: corre `cargo test` en cada módulo Rust independiente
- `test-ai.sh`: crea un `.venv` local para `ai/python`, instala dependencias de test y corre `compileall` + `unittest`
- `test-all.sh`: ejecuta validación de layout, versionado y toda la suite del repo

## Convención

- nombrar por acción y contexto
- no meter scripts específicos de una app externa
- si un script solo sirve para un módulo, evaluar primero ubicarlo dentro de ese módulo
- no asumir una versión global del repo; el versionado siempre es por implementación (`saas/go`, `ai/python`, etc.)
