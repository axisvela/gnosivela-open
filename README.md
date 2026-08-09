# GNOSIVELA-open

> **GNOSIVELA-open** — 开发者入口与发行面（Apache-2.0）。

GNOSIVELA 是 AxisRobo 的 Enterprise Knowledge & Data Fabric。本仓库承载对外宣传面：
**SDK、示例、API 契约、参考 Ontology 与二进制发行**，让开发者 5 分钟内体验语义契约与 Grounding。

## 三仓库 Open Core

| 仓库 | 定位 | License |
| --- | --- | --- |
| **GNOSIVELA-open（本仓库）** | SDK、示例、API 契约、二进制 | Apache-2.0 |
| [GNOSIVELA](https://github.com/axisrobo/GNOSIVELA) | 产品核心代码（Semantic Kernel） | AGPL-3.0 |
| [GNOSIVELA-ee](https://github.com/axisrobo/GNOSIVELA-ee) | 企业版：治理、规模化、行业包 | Enterprise |

> License 合规：本仓库 SDK 为 Apache-2.0，通过 HTTP 与核心交互，不链接 AGPL 核心。

## 目录结构

```
sdk/go/gnosivela/     # Go 客户端 SDK（契约 DTO 由 DSL 编译器同步生成）
sdk/go/gnosivela/gen/ # DSL 编译器生成的领域类型（Go struct + JSON Schema）
examples/             # 端到端示例（quickstart、supplier.dsl）
api/openapi.yaml      # Semantic Control Plane API 契约
```

## DSL 编译器（SDK 类型生成）

领域类型由 core 仓库的 DSL 编译器生成（上游核心库零依赖、Apache 干净的独立 DTO）：

```bash
# 从 Semantic Contract DSL 生成 Go 类型 + JSON Schema 到 sdk/go/gnosivela/gen
cd ../GNOSIVELA
go run ./cmd/gnosivela-gen -in examples/supplier.dsl -out ../GNOSIVELA-open/sdk/go/gnosivela/gen -package gen
```

生成的 `gen.Supplier` 等类型与 core 语义对象一一对应，用于客户端侧的强类型建模。

## 快速开始

```bash
# 1) 启动核心服务（AGPL core 仓库）
cd ../GNOSIVELA && go run ./cmd/gnosivela

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

完整契约见 [api/openapi.yaml](api/openapi.yaml)（语义对象绑定，不绑定任何图/向量/搜索厂商）。

## License

Apache-2.0，见 [LICENSE](LICENSE)。
