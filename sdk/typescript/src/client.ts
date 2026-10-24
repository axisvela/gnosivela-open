/**
 * HTTP client for the GNOSIVELA Semantic Control Plane API.
 * Uses the global fetch (Node >= 18 / modern browsers). Never links the core.
 */
import type { Event, EventContract, GraphQLResponse, KnowledgeAssertion, PolicyRequest, PipelineRequest, ShaclValidationReport } from "./types.ts";

// Re-export the shared contract types from the package entry point so
// consumers can import them as `import type { EntityRef } from "@axisvela/gnosivela"`.
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
  ShaclViolation,
  ShaclValidationReport,
  GraphQLResponse,
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

  private async requestText(method: string, path: string): Promise<string> {
    const resp = await fetch(this.baseURL + path, { method, headers: { ...this.headers } });
    if (!resp.ok) {
      throw new GnosivelaError(resp.status, await resp.text(), path);
    }
    return resp.text();
  }

  private async requestTextBody(method: string, path: string, body: string, contentType: string): Promise<string> {
    const resp = await fetch(this.baseURL + path, {
      method,
      headers: { ...this.headers, "Content-Type": contentType },
      body,
    });
    if (!resp.ok) {
      throw new GnosivelaError(resp.status, await resp.text(), path);
    }
    return resp.text();
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
  hybridQuery(query: string, limit = 10, weights: Record<string, number> = {}): Promise<Record<string, unknown>> {
    return this.request("POST", "/query/hybrid", { query, limit, weights });
  }
  studioModel(namespace: string): Promise<Record<string, unknown>> {
    return this.request("GET", `/studio/model?namespace=${encodeURIComponent(namespace)}`);
  }
  studioImpact(namespace: string): Promise<Record<string, unknown>> {
    return this.request("GET", `/studio/impact?namespace=${encodeURIComponent(namespace)}`);
  }
  sparqlQuery(query: string): Promise<Record<string, unknown>> {
    return this.request("POST", "/sparql", { query });
  }
  sparqlQueryText(query: string): Promise<string> {
    return this.requestTextBody("POST", "/sparql", query, "application/sparql-query");
  }
  shaclExport(namespace: string, version = ""): Promise<string> {
    const suffix = version ? `?version=${encodeURIComponent(version)}` : "";
    return this.requestText("GET", `/interop/shacl/${encodeURIComponent(namespace)}${suffix}`);
  }
  shaclImport(namespace: string, turtle: string, version = ""): Promise<Record<string, unknown>> {
    const suffix = version ? `?version=${encodeURIComponent(version)}` : "";
    return this.requestTextBody("POST", `/interop/shacl/${encodeURIComponent(namespace)}${suffix}`, turtle, "text/turtle")
      .then((text) => JSON.parse(text) as Record<string, unknown>);
  }
  shaclValidate(namespace: string, assertions?: unknown[], version = ""): Promise<ShaclValidationReport> {
    const suffix = version ? `?version=${encodeURIComponent(version)}` : "";
    const path = `/interop/shacl/${encodeURIComponent(namespace)}/validate${suffix}`;
    if (assertions === undefined) {
      return this.request("POST", path);
    }
    return this.request("POST", path, { assertions });
  }
  shaclValidateStored(namespace: string, version = ""): Promise<ShaclValidationReport> {
    const suffix = version ? `?version=${encodeURIComponent(version)}` : "";
    return this.request("GET", `/interop/shacl/${encodeURIComponent(namespace)}/validate${suffix}`);
  }
  owlExport(namespace: string, version = ""): Promise<string> {
    const suffix = version ? `?version=${encodeURIComponent(version)}` : "";
    return this.requestText("GET", `/interop/owl/${encodeURIComponent(namespace)}${suffix}`);
  }
  owlInfer(namespace: string, version = ""): Promise<string> {
    const suffix = version ? `?version=${encodeURIComponent(version)}` : "";
    return this.requestText("GET", `/interop/owl/${encodeURIComponent(namespace)}/infer${suffix}`);
  }
  graphqlQuery(query: string, variables: Record<string, unknown> = {}): Promise<GraphQLResponse> {
    return this.request("POST", "/graphql", { query, variables });
  }
  graphqlSchema(): Promise<string> {
    return this.requestText("GET", "/graphql/schema");
  }
  actions(): Promise<{ actions: Array<Record<string, unknown>> }> {
    return this.request("GET", "/actions");
  }
  actionRegister(definition: Record<string, unknown>): Promise<Record<string, unknown>> {
    return this.request("POST", "/actions", definition);
  }
  actionRun(actionId: string, principal = "", approvalId = "", input?: unknown): Promise<Record<string, unknown>> {
    return this.request("POST", `/actions/${encodeURIComponent(actionId)}/run`, { principal, approvalId, input });
  }
  actionRuns(): Promise<{ runs: Array<Record<string, unknown>> }> {
    return this.request("GET", "/actions/runs");
  }
  connectors(): Promise<{ connectors: Array<Record<string, unknown>> }> {
    return this.request("GET", "/connectors");
  }
  connectorIntrospect(kind: string, spec: Record<string, unknown>): Promise<Record<string, unknown>> {
    return this.request("POST", "/connectors/introspect", { kind, spec });
  }
  connectorBindDocument(docId: string, content: string, uri = "", version = "", extractClaims = false, contentType = "markdown", extractor?: Record<string, unknown>): Promise<Record<string, unknown>> {
    const body: Record<string, unknown> = { docId, content, uri, version, extractClaims, contentType };
    if (extractor) body.extractor = extractor;
    return this.request("POST", "/connectors/documents", body);
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
