import assert from "node:assert";
import { createServer, type Server } from "node:http";
import test from "node:test";

import { Client, GnosivelaError } from "../src/client.ts";

let server: Server;
let base = "";

test.before(async () => {
  server = createServer((req, res) => {
    const send = (payload: unknown, status = 200) => {
      const body = JSON.stringify(payload);
      res.writeHead(status, { "Content-Type": "application/json" });
      res.end(body);
    };
    const path = (req.url ?? "").split("?")[0];
    if (req.method === "GET" && path === "/ontologies/acme/latest") {
      return send({ namespace: "acme", version: "1.2", status: "published" });
    }
    if (req.method === "GET" && path === "/consistency/conflicts") {
      return send({ conflicts: [] });
    }
    if (req.method === "GET" && path === "/events/contracts") {
      return send({ contracts: [{ id: "c-1", type: "price.updated" }] });
    }
    if (req.method === "GET" && path === "/quality") {
      return send({ citationCompleteness: 0.98, unresolvedRate: 0, conflicts: 0 });
    }
    if (req.method === "GET" && path === "/metrics") {
      return send({ counts: { "query.semantic": 3 } });
    }
    if (req.method === "GET" && path === "/bridge/procurement.supplier/contract") {
      return send({ id: "procurement.supplier@1.2", signature: "a".repeat(64), concepts: [] });
    }
    if (req.method === "GET" && path === "/metrics/definitions") {
      return send({ definitions: [{ id: "m:1", formula: "sum(x)" }] });
    }
    if (req.method === "POST" && path === "/incidents/check") {
      return send({ opened: [{ id: "INC-001", status: "open" }], quality: {} });
    }
    if (req.method === "POST" && path === "/metrics/definitions") {
      return send({ id: "m:1", formula: "sum(x)" }, 201);
    }
    if (req.method === "POST" && path === "/policy/evaluate") {
      return send({ allowed: true, reason: "open default" });
    }
    if (req.method === "POST" && path === "/events/ingest") {
      return send({ assertions: [{ assertionId: "ev:e-1:p" }], resolved: ["ev:e-1:p"], gaps: [] }, 201);
    }
    if (req.method === "POST" && path === "/entities") {
      return send({ namespace: "mdm", canonicalId: "C-1" }, 201);
    }
    return send({ error: "not found" }, 404);
  });
  await new Promise<void>((resolve) => server.listen(0, "127.0.0.1", resolve));
  const addr = server.address();
  if (addr && typeof addr === "object") {
    base = `http://127.0.0.1:${addr.port}`;
  }
});

test.after(() => {
  server.close();
});

test("ontology latest", async () => {
  const c = new Client(base);
  const o = await c.ontologyLatest("acme");
  assert.equal(o.version, "1.2");
});

test("governance endpoints", async () => {
  const c = new Client(base);
  const conflicts = await c.consistencyConflicts();
  assert.deepEqual(conflicts.conflicts, []);
  const d = await c.policyEvaluate({ action: "ontology.publish", resource: "ontology:x" });
  assert.equal(d.allowed, true);
});

test("events roundtrip", async () => {
  const c = new Client(base);
  const list = await c.eventContractList();
  assert.equal(list.contracts.length, 1);
  const ing = await c.eventIngest("c-1", { id: "e-1", type: "price.updated", source: "s", payload: {} });
  assert.equal(ing.resolved.length, 1);
});

test("entity save", async () => {
  const c = new Client(base);
  const e = await c.entitySave({ namespace: "mdm", canonicalId: "C-1", type: "Company" });
  assert.equal(e.canonicalId, "C-1");
});

test("operations endpoints", async () => {
  const c = new Client(base);
  const q = await c.quality("");
  assert.ok(q.citationCompleteness > 0.9);
  const counts = await c.metrics();
  assert.equal(counts.counts["query.semantic"], 3);
  const opened = await c.incidentCheck();
  assert.equal(opened.opened[0].status, "open");
  const contract = await c.bridgeContractExport("procurement.supplier");
  assert.equal(contract.signature.length, 64);
  await c.metricDefinitionRegister({ id: "m:1", formula: "sum(x)" });
  const defs = await c.metricDefinitions();
  assert.equal(defs.definitions[0].id, "m:1");
});

test("error surfaces status", async () => {
  const c = new Client(base);
  await assert.rejects(() => c.ontologyLatest("missing"), (err: unknown) => {
    assert.ok(err instanceof GnosivelaError);
    assert.equal((err as GnosivelaError).status, 404);
    return true;
  });
});
