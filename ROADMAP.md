# GNOSIVELA-open Roadmap

This is the **open distribution surface** roadmap (SDKs, examples, API contract,
distribution links). Product phases, core capabilities and EE scope live in the
core repository's `docs/ROADMAP.md` and `docs/OPEN-CORE.md`.

## Current release

- **v1.6.0** (2026-10-24) — same tag as GNOSIVELA Core v1.6.0.
- SDKs expose document `artifactDigest`, IE/OCR providers, column lineage, and
  native Python/TypeScript LangChain/LlamaIndex adapters.
- Verified adapter matrix: Python 3.12.8 (`langchain-core 1.7.0`,
  `llama-index-core 0.14.25`); Node 22 (`@langchain/core 1.3.0`,
  `llamaindex 0.12.1`).
- Distribution: GitHub Release artifacts (Python wheel/sdist, npm tgz, Java
  jar, OpenAPI, checksums); PyPI / npm / Maven Central publication is pending
  registry credentials. Container image publication is deferred.

## Repository division

| Repository | Owns | Does not own |
| --- | --- | --- |
| **GNOSIVELA-open** | Go / Python / TypeScript / Java SDKs, examples, API contract, agent framework adapters, distribution links | Semantic Kernel, connectors, governance, HA/SLA |
| **GNOSIVELA core** | Semantic Kernel, core connectors (SQL / document / catalog / S3 / Iceberg / Glue / Hive), document `artifactDigest` + IE provider SPI, Studio, governance, Docker/Helm runnability baseline | Enterprise auth, HA/sharding/10B, SLO/SLA |
| **GNOSIVELA-EE** | Advanced connectors (streaming/Snowflake/Salesforce/Databricks/SharePoint), CDC, Kerberos/SASL and private-network auth, multi-tenancy, federation governance, HA/sharding/10B, multi-region, SLO/SLA, commercial industry profiles, Managed/Appliance | Public SDK surface, core semantic contracts |

## Next milestones

| Version | Open-surface work | Gate |
| --- | --- | --- |
| v1.5.x | PyPI / npm / Maven Central publication; LangChain / LlamaIndex matrix (Python: langchain-core 1.7.0 / llama-index-core 0.14.25; TS: @langchain/core 1.3.0 / llamaindex 0.12.1); examples and MCP packaging; SDK parity checks | Registry artifacts installable; adapter matrix documented |
| v1.6 | Studio SDK surface; OCR/IE compatibility hardening; compatibility matrix; distribution hardening | SDK parity with core APIs; reproducible examples |
| v1.7 | Partner adapters and ecosystem packages (aligned with core P4-M4) | 2 partner packages published |

## Notes

- open and core always share the same version and tag.
- EE keeps an independent version and tag (`CORE_COMPAT` declares the supported
  core range).
- open never implements the Semantic Kernel or enterprise scale guarantees.
