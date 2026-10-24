import assert from "node:assert";
import test from "node:test";

import type { Client } from "../src/client.ts";
import { langchainTools, llamaIndexTools } from "../src/integrations.ts";

async function hasModule(name: string): Promise<boolean> {
  try {
    await import(name as string);
    return true;
  } catch {
    return false;
  }
}

const hasLangChain = await hasModule("@langchain/core/tools");
const hasLlamaIndex = await hasModule("llamaindex");
const client = {} as Client;

test("langchain tools are native when installed", { skip: !hasLangChain }, async () => {
  const tools = (await langchainTools(client, { principal: "alice" })) as Array<Record<string, unknown>>;
  assert.equal(tools.length, 5);
  assert.equal(tools[0].name, "gnosivela_ground");
  assert.ok("schema" in tools[0]);
});

test("langchain tools fall back to neutral specs", { skip: hasLangChain }, async () => {
  const tools = (await langchainTools(client, { principal: "alice" })) as Array<Record<string, unknown>>;
  assert.equal(tools.length, 5);
  assert.equal(tools[0].name, "gnosivela_ground");
  assert.ok("parameters" in tools[0]);
});

test("llama index tools are native when installed", { skip: !hasLlamaIndex }, async () => {
  const tools = (await llamaIndexTools(client, { principal: "alice" })) as Array<{
    metadata?: { name?: string; parameters?: unknown };
  }>;
  assert.equal(tools.length, 5);
  assert.equal(tools[0].metadata?.name, "gnosivela_ground");
  assert.ok(tools[0].metadata?.parameters);
});

test("llama index tools fall back to neutral specs", { skip: hasLlamaIndex }, async () => {
  const tools = (await llamaIndexTools(client, { principal: "alice" })) as Array<Record<string, unknown>>;
  assert.equal(tools.length, 5);
  assert.equal(tools[0].name, "gnosivela_ground");
  assert.ok("parameters" in tools[0]);
});
