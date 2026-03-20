# serverless

Capacidad reusable para runtimes serverless y adapters AWS.

## Pertenece

- Lambda runtime
- API Gateway envelopes y responses
- SQS
- S3
- DynamoDB
- bootstraps AWS
- routing y validaciones reutilizables en serverless

## No pertenece

- dominio específico de KMA/AlphaCoding
- workflows de una app
- features de negocio concretas como audits, facilities o projects

## Fuente inicial esperada

- `/home/pablo/Projects/AlphaCoding/kma-backend/lambdas/shared`

## Nota

`serverless` es independiente de `backend`, salvo piezas técnicas realmente genéricas.

## Bootstrap inicial creado

Paquetes ya iniciados en este módulo:

- `lambdahttp/`
- `event/`
- `s3store/`
- `sqsqueue/`
- `dynamodbtable/`

El recorte actual ya cubre el runtime HTTP, envelopes de eventos y los adapters AWS más repetidos del stack.
