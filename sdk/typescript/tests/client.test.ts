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
    const sendText = (body: string, status = 200) => {
      res.writeHead(status, { "Content-Type": "text/turtle" });
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
    if (req.method === "POST" && path === "/ontologies") {
      return send({ namespace: "acme", version: "2.0", status: "draft" }, 201);
    }
    if (req.method === "GET" && path === "/audit/verify") {
      return send({ intact: true, brokenAt: -1 });
    }
    if (req.method === "GET" && path === "/audit") {
      return send({ entries: [{ id: "AUD-1", action: "publish", actor: "alice" }] });
    }
    if (req.method === "GET" && path.startsWith("/approval/requests/")) {
      return send({ id: path.split("/").pop(), status: "pending", steps: [] });
    }
    if (req.method === "POST" && path.endsWith("/reject")) {
      return send({ id: path.split("/")[3], status: "rejected" });
    }
    if (req.method === "POST" && path === "/entities/merge-candidate") {
      return send({ relation: { id: "rel:1", authority: "candidate" }, evidence: { status: "auto", score: 1 } }, 201);
    }
    if (req.method === "POST" && path === "/grounding/explain") {
      return send({ intent: { query: "x" }, citations: [] });
    }
    if (req.method === "POST" && path === "/query/hybrid") {
      return send({ query: "supplier risk", provider: "hashing/256", fusion: "rrf", hits: [{ id: "ka:risk-1", score: 0.5 }] });
    }
    if (req.method === "GET" && path === "/connectors") {
      return send({ connectors: [{ kind: "api" }, { kind: "streaming" }] });
    }
    if (req.method === "GET" && path === "/studio/model") {
      return send({ namespace: "acme", version: "1.2", status: "published", concepts: [{ id: "Supplier", assertionCount: 2 }], relations: [] });
    }
    if (req.method === "GET" && path === "/studio/impact") {
      return send({ namespace: "acme", fromVersion: "1.0", toVersion: "1.1", additive: [{ kind: "added", type: "concept", id: "Contract" }], blocking: false });
    }
    if (req.method === "GET" && path.startsWith("/interop/shacl/")) {
      return sendText("@prefix sh: <http://www.w3.org/ns/shacl#> .\nex:SupplierShape a sh:NodeShape .\n");
    }
    if (req.method === "GET" && path === "/actions") {
      return send({ actions: [{ id: "notify", kind: "webhook" }] });
    }
    if (req.method === "GET" && path === "/actions/runs") {
      return send({ runs: [{ id: "ACT-1", status: "succeeded" }] });
    }
    if (req.method === "POST" && path === "/actions") {
      return send({ id: "notify", kind: "webhook" }, 201);
    }
    if (req.method === "POST" && path.startsWith("/actions/") && path.endsWith("/run")) {
      return send({ id: "ACT-1", actionId: "notify", status: "succeeded" });
    }
    if (req.method === "POST" && path === "/connectors/introspect") {
      return send({ kind: "streaming", connector: "streaming", fields: [{ name: "eventId", dataType: "string", required: true }] });
    }
    if (req.method === "POST" && path === "/connectors/documents") {
      return send({ id: "doc:p1", docId: "p1", bindingState: "bound", chunks: [{ id: "c1", docId: "p1", locator: "p#c" }] }, 201);
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

test("hybrid query", async () => {
  const c = new Client(base);
  const res = await c.hybridQuery("supplier risk", 5, { lexical: 1, vector: 1, graph: 0.5 });
  assert.equal(res.fusion, "rrf");
  assert.equal((res.hits as Array<Record<string, unknown>>)[0].id, "ka:risk-1");
});

test("studio model", async () => {
  const c = new Client(base);
  const model = await c.studioModel("acme");
  assert.equal(model.namespace, "acme");
  assert.equal((model.concepts as Array<Record<string, unknown>>)[0].assertionCount, 2);
});

test("studio impact", async () => {
  const c = new Client(base);
  const report = await c.studioImpact("acme");
  assert.equal(report.fromVersion, "1.0");
  assert.equal((report.additive as Array<Record<string, unknown>>)[0].id, "Contract");
});

test("actions", async () => {
  const c = new Client(base);
  const list = await c.actions();
  assert.equal(list.actions[0].id, "notify");
  const created = await c.actionRegister({ id: "notify", kind: "webhook", url: "https://example.com" });
  assert.equal(created.id, "notify");
  const run = await c.actionRun("notify", "alice", "", { x: 1 });
  assert.equal(run.status, "succeeded");
  const runs = await c.actionRuns();
  assert.equal(runs.runs[0].id, "ACT-1");
});

test("shacl export", async () => {
  const c = new Client(base);
  const text = await c.shaclExport("acme", "1.2");
  assert.ok(text.includes("ex:SupplierShape"));
});

test("connectors", async () => {
  const c = new Client(base);
  const list = await c.connectors();
  assert.equal(list.connectors.length, 2);
  const schema = await c.connectorIntrospect("streaming", { options: { inlineSchema: "{}" } });
  assert.equal(schema.kind, "streaming");
  const binding = await c.connectorBindDocument("p1", "Suppliers must be reviewed.", "", "", true);
  assert.equal(binding.docId, "p1");
});

test("error surfaces status", async () => {
  const c = new Client(base);
  await assert.rejects(() => c.ontologyLatest("missing"), (err: unknown) => {
    assert.ok(err instanceof GnosivelaError);
    assert.equal((err as GnosivelaError).status, 404);
    return true;
  });
});

test("api-surface parity methods", async () => {
  const c = new Client(base);
  const o = await c.ontologyCreateJSON({ namespace: "acme", version: "2.0" });
  assert.equal(o.version, "2.0");
  const approval = await c.approvalGet("AR-001");
  assert.equal(approval.status, "pending");
  const rejected = await c.approvalReject("AR-001", "alice", "steward", "nope");
  assert.equal(rejected.status, "rejected");
  const entries = await c.auditListFiltered("alice", "publish", "");
  assert.equal(entries.entries[0].action, "publish");
  const verify = await c.auditVerify();
  assert.equal(verify.intact, true);
  const mc = await c.entityMergeCandidate({ namespace: "crm", canonicalId: "ACME" }, { namespace: "mdm", canonicalId: "C-1042" }, "exactSameAs", "rule");
  assert.equal(mc.relation.authority, "candidate");
  const expl = await c.groundingExplain("ACME risk", "alice", "onboarding");
  assert.ok(expl.intent);
});
