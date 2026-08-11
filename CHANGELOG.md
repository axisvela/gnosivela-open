# Changelog

All notable changes to GNOSIVELA-open (SDKs, examples, API contract) are
documented here, following [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and semantic versioning. Versions are kept in sync with GNOSIVELA Core.

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
