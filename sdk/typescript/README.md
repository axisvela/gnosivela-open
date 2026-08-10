# GNOSIVELA TypeScript SDK

Apache-2.0. Talks to the GNOSIVELA Semantic Control Plane over HTTP (never
links the AGPL core). Uses the global `fetch` — Node >= 22.6 (or any modern
browser).

## Run directly (no build step)

Node 22.6+ runs TypeScript with type stripping:

```bash
cd sdk/typescript
node examples/quickstart.ts
node --test tests/client.test.ts
```

## Compile with tsc

```bash
npm install
npm run build
```

## Usage

```typescript
import { Client } from "@axisrobo/gnosivela";

const c = new Client("http://localhost:8080");
const o = await c.ontologyLatest("procurement.supplier");
const view = await c.bridgeQuery("procurement.supplier", "ACME risk", "risk-officer", "onboarding");
```

See `examples/quickstart.ts` for the full walkthrough.
