# GNOSIVELA-open

> **GNOSIVELA-open** — 开发者入口与发行面（Apache-2.0）。

GNOSIVELA 是 AxisRobo 的 **Enterprise Semantic & Knowledge Fabric**：把企业"概念如何定义、实体如何统一引用、关系如何解释、主张如何证据化、冲突如何治理"沉淀为可消费的语义层。本仓库承载对外发行面：**SDK、示例、API 契约、参考 Ontology 与二进制发行**。

## 解决的问题

企业知识分散在 ERP / CRM / 文档 / 数据湖 / 政策库中，每个 Agent、Planner 或 BI 都各自猜测含义。结果：

- **同名不同义** — 三种 "Active Customer" 定义并存，无人知道用的是哪一个
- **身份分裂** — 同一供应商在 12 个系统里有 12 个 ID，跨源合并靠运气
- **无证据回答** — Agent 输出"看起来对"的结论，却拿不出来源、时间和权限链
- **冲突被静默覆盖** — 新数据覆盖旧数据，而不是保留为竞争性主张
- **厂商锁定** — 语义被绑死在某个图数据库 / 向量库的内部 Schema

GNOSIVELA 的答案是：先定义语义对象（Ontology / Assertion / Identity / Relation / Grounding），再决定存储。**语义控制平面**负责一致性与治理，**Polyglot 投影**（SQL / Graph / Search / Vector）可替换且不反向污染语义 API。

## 产品特性

| 能力 | 说明 |
| --- | --- |
| Semantic Contract DSL | 声明概念、属性、约束与关系；DSL 编译器同步生成 SDK 类型 + JSON Schema |
| Ontology Registry | 版本化 / 语义 Diff / 校验 / 发布 / 弃用 |
| Knowledge Assertion | 每个主张带来源、时间、上下文、状态、置信度与证据链 |
| Entity Identity | 跨源实体解析：精确匹配与证据优先，模糊仅作候选、绝不静默合并 |
| Polyglot Projection | 同一语义对象投影到 Relational / Graph / Search / Vector，可替换、不暴露厂商类型 |
| Semantic Query | Intent 解析 → 实体/关系/文档/指标联合检索 → 目的限定的 KnowledgeView + 冲突/Gap |
| Grounding（Phase 1 进行中） | 面向目的组装 Bundle：intent + assertions + evidence + policy + conflict/gap |

## 三仓库 Open Core

| 仓库 | 定位 | License |
| --- | --- | --- |
| **GNOSIVELA-open（本仓库）** | SDK、示例、API 契约、二进制 | Apache-2.0 |
| [GNOSIVELA](https://github.com/axisrobo/GNOSIVELA) | 产品核心代码（Semantic Kernel） | AGPL-3.0 |
| [GNOSIVELA-ee](https://github.com/axisrobo/GNOSIVELA-ee) | 企业版：治理、规模化、行业包 | Enterprise |

> License 合规：本仓库 SDK 为 Apache-2.0，通过 HTTP 与核心交互，不链接 AGPL 核心。

## 目录结构

```
contracts/             # 语义 API 契约（openapi.yaml）+ DSL 样例（examples/）
sdk/go/gnosivela/     # Go 客户端 SDK（契约 DTO 由 DSL 编译器同步生成）
sdk/go/gnosivela/gen/ # DSL 编译器生成的领域类型（Go struct + JSON Schema）
examples/             # 端到端示例（quickstart）
```

## DSL 编译器（SDK 类型生成）

领域类型由 core 仓库的 DSL 编译器生成（上游核心库零依赖、Apache 干净的独立 DTO）：

```bash
# 从 Semantic Contract DSL 生成 Go 类型 + JSON Schema 到 sdk/go/gnosivela/gen
cd ../GNOSIVELA/backend
go run ./cmd/gnosivela-gen -in ../../GNOSIVELA/contracts/examples/supplier.dsl -out ../../GNOSIVELA-open/sdk/go/gnosivela/gen -package gen
```

生成的 `gen.Supplier` 等类型与 core 语义对象一一对应，用于客户端侧的强类型建模。

## 快速开始

```bash
# 1) 启动核心服务（AGPL core 仓库，模块位于 backend/）
cd ../GNOSIVELA/backend && go run ./cmd/gnosivela

# 2) 运行 Go SDK 示例
go run ./examples/quickstart
```

## 使用 SDK

```go
import gnosivela "github.com/axisrobo/GNOSIVELA-open/sdk/go/gnosivela"

c := gnosivela.New("http://localhost:8080")

res, err := c.OntologyCreate(ctx, dslDoc)      // DSL 创建 Ontology
matches, err := c.EntityResolve(ctx, "ACME")   // 实体解析
a, err := c.AssertionPropose(ctx, assertion)   // 知识主张
```

## API 契约

完整契约见 [contracts/openapi.yaml](contracts/openapi.yaml)（语义对象绑定，不绑定任何图/向量/搜索厂商）。

## License

Apache-2.0，见 [LICENSE](LICENSE)。
