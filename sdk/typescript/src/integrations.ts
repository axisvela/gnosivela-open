/**
 * Dependency-free agent framework adapters for GNOSIVELA.
 *
 * The helpers return neutral tool specs (JSON Schema + invoke) that can be
 * wrapped by LangChain.js, LlamaIndex or an OpenAI-style function-calling
 * loop. They call the governed control-plane API only; semantics are never
 * copied into the agent framework.
 */
import type { Client } from "./client.ts";

export interface ToolSpec {
  name: string;
  description: string;
  parameters: Record<string, unknown>;
  invoke: (args?: Record<string, unknown>) => Promise<string>;
}

export interface ToolOptions {
  principal?: string;
  purpose?: string;
}

function obj(properties: Record<string, unknown>, required: string[]): Record<string, unknown> {
  return { type: "object", properties, required, additionalProperties: false };
}

function str(description: string): Record<string, unknown> {
  return { type: "string", description };
}

function int(description: string): Record<string, unknown> {
  return { type: "integer", description };
}

async function json(value: unknown): Promise<string> {
  return JSON.stringify(value);
}

export function toolSpecs(client: Client, opts: ToolOptions = {}): ToolSpec[] {
  const principal = opts.principal ?? "";
  const purpose = opts.purpose ?? "";
  const base = {
    query: str("Natural-language or structured question."),
    principal: str("Principal on whose behalf the tool runs."),
    purpose: str("Purpose scope for the knowledge view."),
  };
  return [
    {
      name: "gnosivela_ground",
      description:
        "Assemble a governed Grounding Bundle (knowledge, evidence, contract, conflicts, gaps) for a question.",
      parameters: obj({ ...base, budget: int("Max knowledge items (0 = unlimited).") }, ["query"]),
      invoke: async (args = {}) => {
        const q = String(args.query ?? "");
        const p = args.principal === undefined ? principal : String(args.principal);
        const pu = args.purpose === undefined ? purpose : String(args.purpose);
        const b = args.budget === undefined ? 0 : Number(args.budget);
        return json(await client.groundingAssemble(q, p, pu, b));
      },
    },
    {
      name: "gnosivela_cite",
      description: "Explain the citations behind a grounded answer.",
      parameters: obj({ ...base }, ["query"]),
      invoke: async (args = {}) => {
        const q = String(args.query ?? "");
        const p = args.principal === undefined ? principal : String(args.principal);
        const pu = args.purpose === undefined ? purpose : String(args.purpose);
        return json(await client.groundingExplain(q, p, pu));
      },
    },
    {
      name: "gnosivela_search",
      description: "Run a semantic query and return a purpose-scoped KnowledgeView.",
      parameters: obj({ ...base }, ["query"]),
      invoke: async (args = {}) => {
        const q = String(args.query ?? "");
        const p = args.principal === undefined ? principal : String(args.principal);
        const pu = args.purpose === undefined ? purpose : String(args.purpose);
        return json(await client.semanticQuery(q, p, pu));
      },
    },
    {
      name: "gnosivela_hybrid_search",
      description: "Hybrid BM25/vector/graph retrieval fused with RRF.",
      parameters: obj(
        { query: str("Search query."), limit: int("Maximum number of hits (default 10).") },
        ["query"],
      ),
      invoke: async (args = {}) => {
        const q = String(args.query ?? "");
        const limit = args.limit === undefined ? 10 : Number(args.limit);
        return json(await client.hybridQuery(q, limit));
      },
    },
    {
      name: "gnosivela_resolve_entity",
      description: "Resolve an entity hint with precision-first match evidence.",
      parameters: obj({ hint: str("Entity name, alias or canonical ID.") }, ["hint"]),
      invoke: async (args = {}) => json(await client.entityResolve(String(args.hint ?? ""))),
    },
  ];
}

/** OpenAI function-calling / LangChain.js tool descriptors (schema only). */
export function openAIToolDefinitions(specs: ToolSpec[]): Array<Record<string, unknown>> {
  return specs.map((spec) => ({
    type: "function",
    function: {
      name: spec.name,
      description: spec.description,
      parameters: spec.parameters,
    },
  }));
}

/** A name -> invoke map convenient for LlamaIndex / LangChain wrappers. */
export function toolHandlers(specs: ToolSpec[]): Record<string, (args?: Record<string, unknown>) => Promise<string>> {
  const out: Record<string, (args?: Record<string, unknown>) => Promise<string>> = {};
  for (const spec of specs) {
    out[spec.name] = spec.invoke;
  }
  return out;
}

type LangChainToolFactory = (fn: (args: Record<string, unknown>) => Promise<string>, fields: {
  name: string;
  description: string;
  schema: Record<string, unknown>;
}) => unknown;

type LlamaIndexToolFactory = {
  from: (config: {
    name: string;
    description: string;
    parameters: Record<string, unknown>;
    execute: (args: Record<string, unknown>) => Promise<string>;
  }) => unknown;
};

/** Return LangChain.js tools when `@langchain/core` is installed. */
export async function langchainTools(client: Client, opts: ToolOptions = {}): Promise<unknown[]> {
  const specs = toolSpecs(client, opts);
  try {
    const mod = (await import("@langchain/core/tools" as string)) as { tool: LangChainToolFactory };
    return specs.map((spec) =>
      mod.tool(async (args) => spec.invoke(args), {
        name: spec.name,
        description: spec.description,
        schema: spec.parameters,
      }),
    );
  } catch {
    return specs;
  }
}

/** Return LlamaIndex tools when `llamaindex` is installed. */
export async function llamaIndexTools(client: Client, opts: ToolOptions = {}): Promise<unknown[]> {
  const specs = toolSpecs(client, opts);
  try {
    const mod = (await import("llamaindex" as string)) as { FunctionTool: LlamaIndexToolFactory };
    return specs.map((spec) =>
      mod.FunctionTool.from({
        name: spec.name,
        description: spec.description,
        parameters: spec.parameters,
        execute: async (args) => spec.invoke(args),
      }),
    );
  } catch {
    return specs;
  }
}
