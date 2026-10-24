# Changelog

All notable changes to GNOSIVELA-open (SDKs, examples, API contract) are
documented here, following [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and semantic versioning. Versions are kept in sync with GNOSIVELA Core.

## [Unreleased]

### Added

- Ecosystem scaffolding: partner program, connector certification checklist,
  machine-readable compatibility matrix (`compatibility.json`), partner
  registry (`partners.json` + schema), benchmark reproduction guide and issue
  templates under `docs/ecosystem/` and `.github/ISSUE_TEMPLATE/`.
- `scripts/ecosystem-gate.sh` validates ecosystem docs, core/open versions,
  SDK versions, OpenAPI contract parity and the partner registry.
- `scripts/publish-registry.sh` for PyPI / npm / Maven publication with
  `--check` and `--dry-run` modes. PyPI and npm are ready; Maven Central still
  requires `distributionManagement` and credentials.

## [1.6.0] - 2026-10-24

### Added

- `DocumentBindingRequest` supports `contentEncoding` and an `ocr` provider
  (`kind=http`) for binary PDF/image document bindings. The Go/Python/
  TypeScript SDKs expose the new fields/parameters; `artifactDigest` remains
  the digest of the original bytes.
- Verified the Python native adapters on Python 3.12.8 with `langchain-core
  1.7.0` (`StructuredTool`) and `llama-index-core 0.14.25` (`FunctionTool`);
  native tests skip when a framework is absent.
- TypeScript `langchainTools` / `llamaIndexTools` return native
  `DynamicStructuredTool` / `FunctionTool` objects when `@langchain/core` or
  `llamaindex` is installed, and neutral specs otherwise. Verified on Node
  24.18.0 with `@langchain/core 1.3.0` and `llamaindex 0.12.1`.
- `ColumnLineage` in the Go / Python / TypeScript / Java SDKs, backed by
  `POST /lineage/columns` (source column → ontology property, unmapped columns
  and orphan mapping rules).

## [1.5.0] - 2026-10-24

### Added

- Document bindings expose a SHA-256 `artifactDigest` of the original input
  content; claim candidates carry the same digest through `EvidenceRef`.
- `DocumentBindingRequest` accepts an optional `extractor` (`keyword` or `http`)
  to select a claim extraction provider. The Go/Python/TypeScript SDKs expose
  the new field/parameter.
- Connector documentation includes the `hive` Metastore connector and the new
  IE provider.

## [1.4.2] - 2026-10-23

### Changed

- Version aligned with GNOSIVELA Core v1.4.2 (streaming/Snowflake/Databricks/
  SharePoint connectors moved to EE). No SDK API changes.

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
