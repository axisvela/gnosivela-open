# Compatibility Matrix

This matrix records the versions exercised for the v1.7.0 release line. It is a
tested baseline, not a promise that older or newer versions are unsupported.

Machine-readable source: [`compatibility.json`](compatibility.json). Validate it
with `./scripts/ecosystem-gate.sh`. After a version bump, refresh it with
`python scripts/update-compatibility.py` (`--check` reports drift without
writing).

## Product versions

| Component | Version | Notes |
| --- | --- | --- |
| GNOSIVELA core | 1.7.0 | Go module `github.com/axisvela/GNOSIVELA/v3` |
| GNOSIVELA-open | 1.7.0 | Same tag as core |
| GNOSIVELA-EE | 1.36.1 | Independent tag; `CORE_COMPAT >=1.4.0,<2.0.0` |
| OpenAPI contract | 1.7.0 | `contracts/openapi.yaml`, byte-identical in core/open |

## SDKs and runtimes

| SDK | Version | Runtime tested |
| --- | --- | --- |
| Go | 1.7.0 | Go 1.26 toolchain (`go.mod`; local toolchain 1.26.9) |
| Python | 1.7.0 | CPython 3.14.6 (SDK tests), 3.12.8 (native adapters); requires >=3.8 |
| TypeScript | 1.7.0 | Node 24.18.0; `engines.node >=22.6` |
| Java | 1.7.0 | OpenJDK 21.0.8, Maven 3.9.10; compiler target 17 |

## Agent framework adapters

| Framework | Tested version | Runtime | Native type |
| --- | --- | --- | --- |
| `langchain-core` (Python) | 1.7.0 | CPython 3.12.8 | `StructuredTool` |
| `llama-index-core` (Python) | 0.14.25 | CPython 3.12.8 | `FunctionTool` |
| `@langchain/core` (TypeScript) | 1.3.0 | Node 24.18.0 | `DynamicStructuredTool` |
| `llamaindex` (TypeScript) | 0.12.1 | Node 24.18.0 | `FunctionTool` |

Neutral tool specs are always available; native adapters are optional.

## Connectors

| Connector | Kind | Repository | Notes |
| --- | --- | --- | --- |
| Relational | `relational` | core | PostgreSQL-flavored information_schema |
| Catalog | `catalog` | core | Data-product references |
| HTTP API | `api` | core | OpenAPI 3 / JSON Schema |
| S3 object store | `objectstore` | core | SigV4 or anonymous; no data copy |
| Iceberg REST | `iceberg` | core | Table schema only |
| AWS Glue | `glue` | core | SigV4 `GetTable` |
| Hive Metastore | `hive` | core | Thrift `get_table`; no Kerberos/SASL in core |
| Streaming / Schema Registry | `streaming` | EE | Advanced connector |
| Snowflake | `snowflake` | EE | Advanced connector |
| Salesforce | `salesforce` | EE | Advanced connector |
| Databricks Unity Catalog | `warehouse` | EE | Advanced connector |
| SharePoint | `sharepoint` | EE | Advanced connector |

## Known limitations

- Hive enterprise authentication (Kerberos/SASL) and private-network packaging
  are EE scope.
- OCR and IE require an external HTTP provider; no bundled model is shipped.
- Container image publication is deferred for v1.7.0.
- PyPI / npm / Maven Central publication is pending credentials; use the GitHub
  Release artifacts.

## Reporting a compatibility issue

Open an issue with: component and version, runtime/framework version, command,
expected vs actual behavior, and a redacted log. Do not include credentials or
source data.
