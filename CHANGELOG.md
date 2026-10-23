# Changelog

All notable changes to GNOSIVELA-open (SDKs, examples, API contract) are
documented here, following [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and semantic versioning. Versions are kept in sync with GNOSIVELA Core.

## [1.4.1] - 2026-10-23

### Changed

- Version aligned with GNOSIVELA Core v1.4.1 (Salesforce connector moved to
  EE). No SDK API changes.

## [1.4.0] - 2026-10-23

### Changed

- Version aligned with GNOSIVELA Core v1.4.0 (`pkg/kernel` facade and public
  `Connector` SPI). No SDK API changes in this release.

## [1.3.0] - 2026-10-23

### Added

- `HybridQuery` in the Go / Python / TypeScript / Java SDKs, backed by Core's
  new `POST /query/hybrid` endpoint (BM25 + vector + graph, RRF fusion).
- OpenAPI contract updated with `/query/hybrid` and `/mcp` (Model Context
  Protocol JSON-RPC endpoint), plus `HybridWeights` / `HybridHit` /
  `HybridResult` / `MCPRequest` / `MCPResponse` schemas.
- Connector methods (`Connectors`, `ConnectorIntrospect`,
  `ConnectorBindDocument`) in all four SDKs, backed by `/connectors`,
  `/connectors/introspect` and `/connectors/documents`. Introspection covers
  relational, catalog, OpenAPI/JSON Schema, Schema Registry (Avro), S3 object
  stores, Databricks Unity Catalog and Snowflake. Document binding accepts
  markdown, plain text or HTML via `contentType`.
- `StudioModel` and `StudioImpact` in all four SDKs, backed by the read-only
  `GET /studio/model` and `GET /studio/impact` endpoints, including an optional
  change-impact summary for unpublished versions.
- `ShaclExport` in all four SDKs, backed by `GET /interop/shacl/{namespace}`.
- `SparqlQuery` in all four SDKs, backed by the read-only `POST /sparql`
  endpoint (SPARQL 1.1 subset).
- Action methods (`Actions`, `ActionRegister`, `ActionRun`, `ActionRuns`) in all
  four SDKs, backed by the governed webhook runtime (`/actions*`). Action
  definitions accept `schedule`, `schedulePrincipal` and `scheduleInput`.
- Interop methods in all four SDKs: `ShaclImport` / `ShaclValidate`
  (`POST /interop/shacl/{namespace}` and `/validate`), `OwlExport` /
  `OwlInfer` (`/interop/owl/*`), `GraphQLQuery` / `GraphQLSchema`
  (`/graphql`), and `SparqlQueryText` for CONSTRUCT results. Action runs now
  expose the `output` field for `kind=function`.
- Dependency-free agent framework adapters: Python `gnosivela.integrations`
  (`tool_specs`, `langchain_tools`, `llama_index_tools`) and TypeScript
  `@axisvela/gnosivela/integrations` (`toolSpecs`, `openAIToolDefinitions`,
  `toolHandlers`). Native framework objects are returned when available;
  neutral specs otherwise.

## [1.0.0] — 2026-08-21

General availability SDK baseline, aligned with the Core 1.0.0 GA release.

### Changed
- All SDK package versions unified to `1.0.0` (no pre-release suffix; see `VERSION`).
- Container image tag reference updated to `ghcr.io/axisvela/GNOSIVELA:1.0.0`.

## [1.0.0-rc.1] — 2026-08-11

Release candidate SDK baseline, aligned with the Core 1.0.0-rc.1 milestone.

### Changed
- All SDK package versions unified to `1.0.0-rc.1` (see `VERSION`).

## [1.0.0-beta.1] — 2026-08-10

RC-candidate SDK baseline, aligned with the Core 1.0.0-beta.1 milestone.

### Added
- SDKs for Go (reference), Python (standard-library only), TypeScript
  (fetch-based, runs without a build step) and Java (JDK HttpClient, zero
  dependencies).
- DSL compiler multi-language output: Go / Python / TypeScript / JSON Schema.
- Full API surface across all SDKs: ontology lifecycle, assertions, entities,
  semantic query, grounding, consistency, policy, approval, audit
  (incl. `/audit/verify`), pipeline, federation, bridge, events, metrics,
  quality, incidents, metric definitions, industry packs.
- Generated reference DTOs (`gen_supplier.py`, `gen_supplier.ts`,
  `sdk/go/gnosivela/gen/`).

### Changed
- All SDK package versions unified to `1.0.0-beta.1` (see `VERSION`).
