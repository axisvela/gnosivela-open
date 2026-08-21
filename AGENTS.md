# GNOSIVELA-open — AGENTS.md

## 仓库定位

GNOSIVELA-open（Apache-2.0）是对外发行面：SDK / 示例 / API 契约 / 二进制发行。三仓库结构见 `../GNOSIVELA/docs/OPEN-CORE.md`：
- **GNOSIVELA-open**（Apache-2.0）：SDK / 示例 / API 契约 / 二进制发行面（本仓库）
- **GNOSIVELA (core)**（AGPL-3.0）：Semantic Kernel 与全部核心模块
- **GNOSIVELA-EE**（Enterprise License）：企业版高价值与规模化能力

## 版本与发布规则

- **版本格式**：`major.minor.patch`（SemVer），**不带后缀**（禁止 `-rc.N`、`-beta.N` 等预发布后缀）。
- **版本同步**：本仓库（`open`）与 `core` 使用**相同版本号、相同 git tag**（两仓库在同一发布点打同名 tag）。
- **EE 独立版本**：`EE` 可有**独立 tag**，不与 open/core 强制同步；EE 以 core 为基底扩展。
- 版本唯一来源：core 仓库 `VERSION` 与 `backend/pkg/version/version.go`；本仓库 `VERSION` 与 SDK 三件套（`sdk/python/pyproject.toml`、`sdk/typescript/package.json`、`sdk/java/pom.xml`）必须与之一致，core 的 `scripts/release-gate.sh` 强制校验。
- 发布流程详见 core 仓库 `docs/RELEASE.md` 与 `docs/plans/GA-RELEASE-PLAN.md`。
