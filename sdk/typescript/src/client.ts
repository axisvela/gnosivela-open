/**
 * HTTP client for the GNOSIVELA Semantic Control Plane API.
 * Uses the global fetch (Node >= 18 / modern browsers). Never links the core.
 */
import type { Event, EventContract, KnowledgeAssertion, PolicyRequest, PipelineRequest } from "./types.ts";

// Re-export the shared contract types from the package entry point so
// consumers can import them as `import type { EntityRef } from "@axisrobo/gnosivela"`.
export type {
  EntityRef,
  AssertionContext,
  Value,
  EvidenceRef,
  KnowledgeAssertion,
  Conflict,
  PolicyRequest,
  PolicyDecision,
  PipelineRequest,
  EventContract,
  EventTemplate,
  Event,
} from "./types.ts";

export class GnosivelaError extends Error {
  status: number;
  body: string;
  path: string;

  constructor(status: number, body: string, path: string) {
    super(`gnosivela: ${path} -> ${status}: ${body}`);
    this.status = status;
    this.body = body;
    this.path = path;
  }
}

export interface ClientOptions {
  headers?: Record<string, string>;
}

export class Client {
  private baseURL: string;
  private headers: Record<string, string>;

  constructor(baseURL = "http://localhost:8080", opts: ClientOptions = {}) {
    this.baseURL = baseURL.replace(/\/$/, "");
    this.headers = { "Content-Type": "application/json", ...(opts.headers ?? {}) };
  }

  private async request<T>(method: string, path: string, body?: unknown): Promise<T> {
    const url = this.baseURL + path;
    const init: RequestInit = { method, headers: { ...this.headers } };
    if (body !== undefined) {
      init.body = typeof body === "string" ? body : JSON.stringify(body);
    }
    const resp = await fetch(url, init);
    if (!resp.ok) {
      throw new GnosivelaError(resp.status, await resp.text(), path);
    }
    if (resp.status === 204) {
      return undefined as T;
    }
    return (await resp.json()) as T;
  }

  // ---- ontology ----
  ontologyCreate(dsl: string): Promise<Record<string, unknown>> {
    return this.request("POST", "/ontologies", dsl);
  }
  ontologyCreateJSON(ontology: Record<string, unknown>): Promise<Record<string, unknown>> {
    return this.request("POST", "/ontologies", ontology);
  }
  ontologyLatest(namespace: string): Promise<Record<string, unknown>> {
    return this.request("GET", `/ontologies/${encodeURIComponent(namespace)}/latest`);
  }
  ontologyGet(namespace: string, version: string): Promise<Record<string, unknown>> {
    return this.request("GET", `/ontologies/${encodeURIComponent(namespace)}/versions/${encodeURIComponent(version)}`);
  }
  ontologyPublish(namespace: string, version: string, approval = ""): Promise<Record<string, unknown>> {
    const q = approval ? `?approval=${encodeURIComponent(approval)}` : "";
    return this.request("POST", `/ontologies/${encodeURIComponent(namespace)}/versions/${encodeURIComponent(version)}/publish${q}`);
  }
  ontologyImpact(namespace: string, version: string): Promise<Record<string, unknown>> {
    return this.request("GET", `/ontologies/${encodeURIComponent(namespace)}/versions/${encodeURIComponent(version)}/impact`);
  }
  ontologyRollback(namespace: string, version: string): Promise<Record<string, unknown>> {
    return this.request("POST", `/ontologies/${encodeURIComponent(namespace)}/versions/${encodeURIComponent(version)}/rollback`);
  }
  ontologyDiff(namespace: string, version: string, other: string): Promise<Record<string, unknown>> {
    return this.request("GET", `/ontologies/${encodeURIComponent(namespace)}/versions/${encodeURIComponent(version)}/diff?other=${encodeURIComponent(other)}`);
  }

  // ---- assertion / entity ----
  assertionPropose(a: KnowledgeAssertion): Promise<Record<string, unknown>> {
    return this.request("POST", "/assertions", a);
  }
  assertionList(subjectNs: string, subjectId: string): Promise<Record<string, unknown>> {
    return this.request("GET", `/assertions?subjectNs=${encodeURIComponent(subjectNs)}&subjectId=${encodeURIComponent(subjectId)}`);
  }
  entitySave(e: Record<string, unknown>): Promise<Record<string, unknown>> {
    return this.request("POST", "/entities", e);
  }
  entityResolve(hint: string): Promise<Record<string, unknown>> {
    return this.request("POST", "/entities/resolve", { hint });
  }
  entityExplain(hint: string): Promise<Record<string, unknown>> {
    return this.request("POST", "/entities/explain", { hint });
  }
  entityMergeCandidate(left: Record<string, unknown>, right: Record<string, unknown>, kind: string, generatedBy: string): Promise<Record<string, unknown>> {
    return this.request("POST", "/entities/merge-candidate", { left, right, kind, generatedBy });
  }

  // ---- query / grounding ----
  semanticQuery(query: string, principal = "", purpose = ""): Promise<Record<string, unknown>> {
    return this.request("POST", "/query/semantic", { query, principal, purpose });
  }
  pathQuery(from: string, to: string): Promise<Record<string, unknown>> {
    return this.request("POST", "/query/path", { from, to });
  }
  subgraphQuery(node: string): Promise<Record<string, unknown>> {
    return this.request("POST", "/query/subgraph", { node });
  }
  groundingAssemble(query: string, principal: string, purpose: string, budget = 0): Promise<Record<string, unknown>> {
    return this.request("POST", "/grounding/assemble", { query, principal, purpose, budget });
  }
  groundingExplain(query: string, principal: string, purpose: string): Promise<Record<string, unknown>> {
    return this.request("POST", "/grounding/explain", { query, principal, purpose });
  }
  groundingRedact(query: string, principal: string, purpose: string, hideTags: string[] = []): Promise<Record<string, unknown>> {
    return this.request("POST", "/grounding/redact", { query, principal, purpose, hideTags });
  }

  // ---- governance ----
  consistencyReport(ontologyNamespace = ""): Promise<Record<string, unknown>> {
    const q = ontologyNamespace ? `?ontology=${encodeURIComponent(ontologyNamespace)}` : "";
    return this.request("GET", `/consistency/report${q}`);
  }
  consistencyConflicts(): Promise<{ conflicts: unknown[] }> {
    return this.request("GET", "/consistency/conflicts");
  }
  consistencyResolve(): Promise<{ resolutions: unknown[] }> {
    return this.request("POST", "/consistency/resolve");
  }
  consistencyAudit(): Promise<{ resolutions: unknown[] }> {
    return this.request("GET", "/consistency/audit");
  }
  policyList(): Promise<Record<string, unknown>> {
    return this.request("GET", "/policy/policies");
  }
  policyEvaluate(req: PolicyRequest): Promise<{ allowed: boolean; reason: string }> {
    return this.request("POST", "/policy/evaluate", req);
  }
  approvalCreate(action: string, resource: string, requester = "", reason = ""): Promise<Record<string, unknown>> {
    return this.request("POST", "/approval/requests", { action, resource, requester, reason });
  }
  approvalList(): Promise<{ requests: unknown[] }> {
    return this.request("GET", "/approval/requests");
  }
  approvalGet(id: string): Promise<Record<string, unknown>> {
    return this.request("GET", `/approval/requests/${encodeURIComponent(id)}`);
  }
  approvalApprove(id: string, approver: string, role: string, comment = ""): Promise<Record<string, unknown>> {
    return this.request("POST", `/approval/requests/${encodeURIComponent(id)}/approve`, { approver, role, comment });
  }
  approvalReject(id: string, approver: string, role: string, comment = ""): Promise<Record<string, unknown>> {
    return this.request("POST", `/approval/requests/${encodeURIComponent(id)}/reject`, { approver, role, comment });
  }
  auditList(): Promise<{ entries: unknown[] }> {
    return this.request("GET", "/audit");
  }
  auditListFiltered(actor = "", action = "", resource = ""): Promise<{ entries: unknown[] }> {
    const q = new URLSearchParams();
    if (actor) q.set("actor", actor);
    if (action) q.set("action", action);
    if (resource) q.set("resource", resource);
    return this.request("GET", `/audit?${q.toString()}`);
  }
  auditVerify(): Promise<{ intact: boolean; brokenAt: number }> {
    return this.request("GET", "/audit/verify");
  }
  auditAttest(entryId: string, ref: string, content: string, by: string): Promise<Record<string, unknown>> {
    return this.request("POST", "/audit/attest", { entryId, ref, content, by });
  }

  // ---- pipeline / federation / bridge / events ----
  pipelineRun(req: PipelineRequest): Promise<Record<string, unknown>> {
    return this.request("POST", "/pipeline/run", req);
  }
  federationQuery(query: string, principal = "", purpose = ""): Promise<Record<string, unknown>> {
    return this.request("POST", "/federation/query", { query, principal, purpose });
  }
  federationDomains(): Promise<{ domains: unknown[] }> {
    return this.request("GET", "/federation/domains");
  }
  federationDomainAdd(name: string, baseUrl: string): Promise<Record<string, unknown>> {
    return this.request("POST", "/federation/domains", { name, baseUrl });
  }
  bridgeContractExport(namespace: string): Promise<Record<string, unknown>> {
    return this.request("GET", `/bridge/${encodeURIComponent(namespace)}/contract`);
  }
  bridgeQuery(namespace: string, query: string, principal = "", purpose = ""): Promise<Record<string, unknown>> {
    return this.request("POST", "/bridge/query", { namespace, query, principal, purpose });
  }
  eventContractRegister(contract: EventContract): Promise<Record<string, unknown>> {
    return this.request("POST", "/events/contracts", contract);
  }
  eventContractList(): Promise<{ contracts: unknown[] }> {
    return this.request("GET", "/events/contracts");
  }
  eventIngest(contractId: string, ev: Event): Promise<{ assertions: unknown[]; resolved: string[]; gaps: string[] }> {
    return this.request("POST", "/events/ingest", { contractId, event: ev });
  }
  metrics(): Promise<{ counts: Record<string, number> }> {
    return this.request("GET", "/metrics");
  }
  quality(ontologyNamespace = ""): Promise<Record<string, unknown>> {
    const q = ontologyNamespace ? `?ontology=${encodeURIComponent(ontologyNamespace)}` : "";
    return this.request("GET", `/quality${q}`);
  }
  incidentRuleAdd(rule: Record<string, unknown>): Promise<Record<string, unknown>> {
    return this.request("POST", "/incidents/rules", rule);
  }
  incidentCheck(): Promise<{ opened: unknown[]; quality: Record<string, unknown> }> {
    return this.request("POST", "/incidents/check");
  }
  incidentList(): Promise<{ incidents: unknown[] }> {
    return this.request("GET", "/incidents");
  }
  incidentResolve(id: string): Promise<Record<string, unknown>> {
    return this.request("POST", `/incidents/${encodeURIComponent(id)}/resolve`);
  }
  metricDefinitions(): Promise<{ definitions: unknown[] }> {
    return this.request("GET", "/metrics/definitions");
  }
  metricDefinitionRegister(def: Record<string, unknown>): Promise<Record<string, unknown>> {
    return this.request("POST", "/metrics/definitions", def);
  }
  industryPacks(): Promise<{ packs: unknown[] }> {
    return this.request("GET", "/industry/packs");
  }
  industryPack(id: string): Promise<Record<string, unknown>> {
    return this.request("GET", `/industry/packs/${encodeURIComponent(id)}`);
  }
}
