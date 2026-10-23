import assert from "node:assert";
import test from "node:test";

import type { Client } from "../src/client.ts";
import { openAIToolDefinitions, toolHandlers, toolSpecs } from "../src/integrations.ts";

class FakeClient {
  calls: Array<unknown[]> = [];
  async groundingAssemble(query: string, principal: string, purpose: string, budget: number) {
    this.calls.push(["ground", query, principal, purpose, budget]);
    return { resolvedIntent: query, knowledge: [{ assertionId: "a1" }] };
  }
  async groundingExplain(query: string, principal: string, purpose: string) {
    this.calls.push(["cite", query, principal, purpose]);
    return { intent: query, citations: [{ assertionId: "a1" }], citationCompleteness: 1 };
  }
  async semanticQuery(query: string, principal: string, purpose: string) {
    this.calls.push(["search", query, principal, purpose]);
    return { intent: { query }, view: { assertions: [] } };
  }
  async hybridQuery(query: string, limit: number) {
    this.calls.push(["hybrid", query, limit]);
    return { query, fusion: "rrf", hits: [] };
  }
  async entityResolve(hint: string) {
    this.calls.push(["resolve", hint]);
    return { hint, canonical: { canonicalId: "C-1" } };
  }
}

function fake(): { client: Client; inner: FakeClient } {
  const inner = new FakeClient();
  return { client: inner as unknown as Client, inner };
}

test("tool specs expose five tools", () => {
  const { client } = fake();
  const specs = toolSpecs(client, { principal: "alice", purpose: "onboarding" });
  assert.equal(specs.length, 5);
  assert.ok(specs.some((s) => s.name === "gnosivela_ground"));
  assert.ok(specs.some((s) => s.name === "gnosivela_hybrid_search"));
  assert.equal((specs[0].parameters as Record<string, unknown>).type, "object");
});

test("ground invocation forwards principal, purpose and budget", async () => {
  const { client, inner } = fake();
  const ground = toolSpecs(client, { principal: "alice", purpose: "onboarding" }).find(
    (s) => s.name === "gnosivela_ground",
  );
  assert.ok(ground);
  const out = JSON.parse(await ground.invoke({ query: "ACME risk", budget: 3 }));
  assert.equal(out.resolvedIntent, "ACME risk");
  assert.deepEqual(inner.calls[0], ["ground", "ACME risk", "alice", "onboarding", 3]);
});

test("hybrid and resolve tools", async () => {
  const { client } = fake();
  const byName = toolHandlers(toolSpecs(client, { principal: "bob" }));
  const hybrid = JSON.parse(await byName.gnosivela_hybrid_search({ query: "q", limit: 3 }));
  assert.equal(hybrid.fusion, "rrf");
  const resolved = JSON.parse(await byName.gnosivela_resolve_entity({ hint: "ACME" }));
  assert.equal(resolved.canonical.canonicalId, "C-1");
});

test("OpenAI tool definitions mirror the specs", () => {
  const { client } = fake();
  const defs = openAIToolDefinitions(toolSpecs(client));
  assert.equal(defs.length, 5);
  const first = defs[0] as { type: string; function: { name: string; parameters: unknown } };
  assert.equal(first.type, "function");
  assert.ok(first.function.name.startsWith("gnosivela_"));
  assert.ok(first.function.parameters);
});
