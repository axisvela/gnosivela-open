package gnosivela

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client is a thin HTTP client for the GNOSIVELA Semantic Control Plane API.
type Client struct {
	baseURL string
	http    *http.Client
}

// New returns a Client targeting the given base URL (e.g. http://localhost:8080).
func New(baseURL string) *Client {
	return &Client{baseURL: baseURL, http: &http.Client{}}
}

// WithHTTPClient overrides the underlying HTTP client (timeouts, transport).
func (c *Client) WithHTTPClient(hc *http.Client) *Client {
	c.http = hc
	return c
}

// OntologyCreate creates an ontology from a DSL document (text/plain).
func (c *Client) OntologyCreate(ctx context.Context, dsl string) (*OntologyCreateResult, error) {
	var out OntologyCreateResult
	err := c.do(ctx, http.MethodPost, "/ontologies", "text/plain", []byte(dsl), &out)
	return &out, err
}

// OntologyCreateJSON creates an ontology from a JSON ontology object.
func (c *Client) OntologyCreateJSON(ctx context.Context, o *Ontology) (*OntologyCreateResult, error) {
	body, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	var out OntologyCreateResult
	err = c.do(ctx, http.MethodPost, "/ontologies", "application/json", body, &out)
	return &out, err
}

// OntologyLatest fetches the latest version of a namespace.
func (c *Client) OntologyLatest(ctx context.Context, namespace string) (*Ontology, error) {
	var out Ontology
	err := c.do(ctx, http.MethodGet, "/ontologies/"+url.PathEscape(namespace)+"/latest", "", nil, &out)
	return &out, err
}

// OntologyGet fetches a specific version.
func (c *Client) OntologyGet(ctx context.Context, namespace, version string) (*Ontology, error) {
	var out Ontology
	err := c.do(ctx, http.MethodGet,
		"/ontologies/"+url.PathEscape(namespace)+"/versions/"+url.PathEscape(version), "", nil, &out)
	return &out, err
}

// OntologyPublish publishes an ontology version.
func (c *Client) OntologyPublish(ctx context.Context, namespace, version string) error {
	var out PublishResult
	return c.do(ctx, http.MethodPost,
		"/ontologies/"+url.PathEscape(namespace)+"/versions/"+url.PathEscape(version)+"/publish",
		"", nil, &out)
}

// OntologyPublishWithApproval publishes an ontology version that has breaking
// changes, using an approval token.
func (c *Client) OntologyPublishWithApproval(ctx context.Context, namespace, version, approval string) (*PublishResult, error) {
	var out PublishResult
	err := c.do(ctx, http.MethodPost,
		"/ontologies/"+url.PathEscape(namespace)+"/versions/"+url.PathEscape(version)+"/publish?approval="+url.QueryEscape(approval),
		"", nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// OntologyImpact analyses the impact of publishing version against the latest
// published version of the namespace.
func (c *Client) OntologyImpact(ctx context.Context, namespace, version string) (*ImpactReport, error) {
	var out ImpactReport
	err := c.do(ctx, http.MethodGet,
		"/ontologies/"+url.PathEscape(namespace)+"/versions/"+url.PathEscape(version)+"/impact",
		"", nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// OntologyRollback rolls back a published version, deprecating it and
// returning the version the namespace resolves to from now on.
func (c *Client) OntologyRollback(ctx context.Context, namespace, version string) (string, error) {
	var out struct {
		Status    string `json:"status"`
		ToVersion string `json:"toVersion"`
	}
	err := c.do(ctx, http.MethodPost,
		"/ontologies/"+url.PathEscape(namespace)+"/versions/"+url.PathEscape(version)+"/rollback",
		"", nil, &out)
	return out.ToVersion, err
}

// OntologyDiff computes the semantic diff between two versions.
func (c *Client) OntologyDiff(ctx context.Context, namespace, from, to string) ([]Change, error) {
	var out struct {
		Changes []Change `json:"changes"`
	}
	err := c.do(ctx, http.MethodGet,
		"/ontologies/"+url.PathEscape(namespace)+"/versions/"+url.PathEscape(from)+"/diff?other="+url.QueryEscape(to),
		"", nil, &out)
	return out.Changes, err
}

// ConsistencyReport runs rule validation and the seven consistency checks.
// An empty ontologyNamespace runs the store-wide checks only; pass a namespace
// to also enable structural and rule validation.
func (c *Client) ConsistencyReport(ctx context.Context, ontologyNamespace string) (*ConsistencyReport, error) {
	path := "/consistency/report"
	if ontologyNamespace != "" {
		path += "?ontology=" + url.QueryEscape(ontologyNamespace)
	}
	var out ConsistencyReport
	err := c.do(ctx, http.MethodGet, path, "", nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// ConsistencyConflicts lists the detected groups of competing assertions.
func (c *Client) ConsistencyConflicts(ctx context.Context) ([]*Conflict, error) {
	var out struct {
		Conflicts []*Conflict `json:"conflicts"`
	}
	err := c.do(ctx, http.MethodGet, "/consistency/conflicts", "", nil, &out)
	return out.Conflicts, err
}

// ConsistencyResolve applies the resolution policy to all detected conflicts,
// marking losers contested and escalating ties. Returns the recorded decisions.
func (c *Client) ConsistencyResolve(ctx context.Context) ([]*ConsistencyResolution, error) {
	var out struct {
		Resolutions []*ConsistencyResolution `json:"resolutions"`
	}
	err := c.do(ctx, http.MethodPost, "/consistency/resolve", "", nil, &out)
	return out.Resolutions, err
}

// ConsistencyAudit lists the recorded resolution audit trail.
func (c *Client) ConsistencyAudit(ctx context.Context) ([]*ConsistencyResolution, error) {
	var out struct {
		Resolutions []*ConsistencyResolution `json:"resolutions"`
	}
	err := c.do(ctx, http.MethodGet, "/consistency/audit", "", nil, &out)
	return out.Resolutions, err
}

// PolicyList lists the registered fine-grained policies and engine status.
func (c *Client) PolicyList(ctx context.Context) ([]*Policy, bool, error) {
	var out struct {
		Enabled  bool      `json:"enabled"`
		Policies []*Policy `json:"policies"`
	}
	err := c.do(ctx, http.MethodGet, "/policy/policies", "", nil, &out)
	return out.Policies, out.Enabled, err
}

// PolicyEvaluate evaluates a policy request and returns the decision.
func (c *Client) PolicyEvaluate(ctx context.Context, req PolicyRequest) (*PolicyDecision, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var out PolicyDecision
	err = c.do(ctx, http.MethodPost, "/policy/evaluate", "application/json", body, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// ApprovalCreate opens an approval request for a governed action.
func (c *Client) ApprovalCreate(ctx context.Context, create ApprovalCreate) (*ApprovalRequest, error) {
	body, err := json.Marshal(create)
	if err != nil {
		return nil, err
	}
	var out ApprovalRequest
	err = c.do(ctx, http.MethodPost, "/approval/requests", "application/json", body, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// ApprovalList lists approval requests in creation order.
func (c *Client) ApprovalList(ctx context.Context) ([]*ApprovalRequest, error) {
	var out struct {
		Requests []*ApprovalRequest `json:"requests"`
	}
	err := c.do(ctx, http.MethodGet, "/approval/requests", "", nil, &out)
	return out.Requests, err
}

// ApprovalGet fetches one approval request.
func (c *Client) ApprovalGet(ctx context.Context, id string) (*ApprovalRequest, error) {
	var out ApprovalRequest
	err := c.do(ctx, http.MethodGet, "/approval/requests/"+url.PathEscape(id), "", nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// ApprovalApprove approves a pending step for the given role.
func (c *Client) ApprovalApprove(ctx context.Context, id, approver, comment string, role ApprovalRole) (*ApprovalRequest, error) {
	body, _ := json.Marshal(map[string]any{"approver": approver, "role": role, "comment": comment})
	var out ApprovalRequest
	err := c.do(ctx, http.MethodPost, "/approval/requests/"+url.PathEscape(id)+"/approve", "application/json", body, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// ApprovalReject rejects a pending step, settling the request as rejected.
func (c *Client) ApprovalReject(ctx context.Context, id, approver, comment string, role ApprovalRole) (*ApprovalRequest, error) {
	body, _ := json.Marshal(map[string]any{"approver": approver, "role": role, "comment": comment})
	var out ApprovalRequest
	err := c.do(ctx, http.MethodPost, "/approval/requests/"+url.PathEscape(id)+"/reject", "application/json", body, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// AuditList returns the governance audit trail in chronological order.
func (c *Client) AuditList(ctx context.Context) ([]*AuditEntry, error) {
	var out struct {
		Entries []*AuditEntry `json:"entries"`
	}
	err := c.do(ctx, http.MethodGet, "/audit", "", nil, &out)
	return out.Entries, err
}

// AuditAttest pins a sha256 digest over content against an audit entry.
func (c *Client) AuditAttest(ctx context.Context, entryID, ref, content, by string) (*AuditEntry, error) {
	body, _ := json.Marshal(map[string]any{"entryId": entryID, "ref": ref, "content": content, "by": by})
	var out AuditEntry
	err := c.do(ctx, http.MethodPost, "/audit/attest", "application/json", body, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// AuditVerify re-computes the tamper-evident audit hash chain. intact is false
// and brokenAt >= 0 when an entry was tampered with.
func (c *Client) AuditVerify(ctx context.Context) (intact bool, brokenAt int, err error) {
	var out struct {
		Intact   bool `json:"intact"`
		BrokenAt int  `json:"brokenAt"`
	}
	err = c.do(ctx, http.MethodGet, "/audit/verify", "", nil, &out)
	return out.Intact, out.BrokenAt, err
}

// PipelineRun executes the mapping CI/CD closed loop and returns the stage
// report. A blocked run (approval/impact gate) surfaces as a non-nil error.
func (c *Client) PipelineRun(ctx context.Context, req PipelineRequest) (*PipelineReport, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var out PipelineReport
	err = c.do(ctx, http.MethodPost, "/pipeline/run", "application/json", body, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// FederationQuery runs a federated semantic query across the registered
// domains and returns the merged, conflict-preserving view.
func (c *Client) FederationQuery(ctx context.Context, query, principal, purpose string) (*FederatedView, error) {
	body, _ := json.Marshal(map[string]any{"query": query, "principal": principal, "purpose": purpose})
	var out FederatedView
	err := c.do(ctx, http.MethodPost, "/federation/query", "application/json", body, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// FederationDomains lists the registered federated domains.
func (c *Client) FederationDomains(ctx context.Context) ([]map[string]any, error) {
	var out struct {
		Domains []map[string]any `json:"domains"`
	}
	err := c.do(ctx, http.MethodGet, "/federation/domains", "", nil, &out)
	return out.Domains, err
}

// FederationDomainAdd registers a remote federated domain (HTTP endpoint).
func (c *Client) FederationDomainAdd(ctx context.Context, name, baseURL string) error {
	body, _ := json.Marshal(map[string]any{"name": name, "baseUrl": baseURL})
	return c.do(ctx, http.MethodPost, "/federation/domains", "application/json", body, nil)
}

// BridgeContractExport fetches the signed semantic contract of the latest
// published ontology version of a namespace.
func (c *Client) BridgeContractExport(ctx context.Context, namespace string) (*BridgeContract, error) {
	var out BridgeContract
	err := c.do(ctx, http.MethodGet, "/bridge/"+url.PathEscape(namespace)+"/contract", "", nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// BridgeQuery runs a contract-driven query filtered to the contract's declared
// predicates, pinning the returned view to the governing contract.
func (c *Client) BridgeQuery(ctx context.Context, namespace, query, principal, purpose string) (*ContractView, error) {
	body, _ := json.Marshal(map[string]any{"namespace": namespace, "query": query, "principal": principal, "purpose": purpose})
	var out ContractView
	err := c.do(ctx, http.MethodPost, "/bridge/query", "application/json", body, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// EventContractRegister registers an event contract.
func (c *Client) EventContractRegister(ctx context.Context, contract EventContract) error {
	body, err := json.Marshal(contract)
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodPost, "/events/contracts", "application/json", body, nil)
}

// EventContractList lists the registered event contracts.
func (c *Client) EventContractList(ctx context.Context) ([]*EventContract, error) {
	var out struct {
		Contracts []*EventContract `json:"contracts"`
	}
	err := c.do(ctx, http.MethodGet, "/events/contracts", "", nil, &out)
	return out.Contracts, err
}

// EventIngest ingests an event and persists the governed assertions it produces.
func (c *Client) EventIngest(ctx context.Context, contractID string, ev *Event) (*IngestResult, error) {
	body, _ := json.Marshal(map[string]any{"contractId": contractID, "event": ev})
	var out IngestResult
	err := c.do(ctx, http.MethodPost, "/events/ingest", "application/json", body, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Metrics returns the operational counters snapshot.
func (c *Client) Metrics(ctx context.Context) (map[string]int64, error) {
	var out struct {
		Counts map[string]int64 `json:"counts"`
	}
	err := c.do(ctx, http.MethodGet, "/metrics", "", nil, &out)
	return out.Counts, err
}

// Quality returns the aggregated quality/SLO snapshot. An empty
// ontologyNamespace runs store-wide checks; pass one to also include
// structural/rule findings.
func (c *Client) Quality(ctx context.Context, ontologyNamespace string) (*QualityReport, error) {
	path := "/quality"
	if ontologyNamespace != "" {
		path += "?ontology=" + url.QueryEscape(ontologyNamespace)
	}
	var out QualityReport
	err := c.do(ctx, http.MethodGet, path, "", nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// IncidentRuleAdd registers an SLO rule.
func (c *Client) IncidentRuleAdd(ctx context.Context, rule IncidentRule) error {
	body, err := json.Marshal(rule)
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodPost, "/incidents/rules", "application/json", body, nil)
}

// IncidentCheck evaluates SLO rules against the current quality snapshot and
// returns the incidents opened on the failing edge.
func (c *Client) IncidentCheck(ctx context.Context) ([]*Incident, *QualityReport, error) {
	var out struct {
		Opened  []*Incident   `json:"opened"`
		Quality *QualityReport `json:"quality"`
	}
	err := c.do(ctx, http.MethodPost, "/incidents/check", "", nil, &out)
	return out.Opened, out.Quality, err
}

// IncidentList lists incidents, newest first.
func (c *Client) IncidentList(ctx context.Context) ([]*Incident, error) {
	var out struct {
		Incidents []*Incident `json:"incidents"`
	}
	err := c.do(ctx, http.MethodGet, "/incidents", "", nil, &out)
	return out.Incidents, err
}

// IncidentResolve resolves an open incident explicitly.
func (c *Client) IncidentResolve(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodPost, "/incidents/"+url.PathEscape(id)+"/resolve", "", nil, nil)
}

// MetricDefinitions lists the governed metric definitions.
func (c *Client) MetricDefinitions(ctx context.Context) ([]MetricDefinition, error) {
	var out struct {
		Definitions []MetricDefinition `json:"definitions"`
	}
	err := c.do(ctx, http.MethodGet, "/metrics/definitions", "", nil, &out)
	return out.Definitions, err
}

// MetricDefinitionRegister registers a governed metric definition.
func (c *Client) MetricDefinitionRegister(ctx context.Context, def MetricDefinition) error {
	body, err := json.Marshal(def)
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodPost, "/metrics/definitions", "application/json", body, nil)
}

// IndustryPacks lists the industry reference packs.
func (c *Client) IndustryPacks(ctx context.Context) ([]IndustryPackSummary, error) {
	var out struct {
		Packs []IndustryPackSummary `json:"packs"`
	}
	err := c.do(ctx, http.MethodGet, "/industry/packs", "", nil, &out)
	return out.Packs, err
}

// IndustryPack returns a full industry reference pack including its ontology.
func (c *Client) IndustryPack(ctx context.Context, id string) (map[string]any, error) {
	var out map[string]any
	err := c.do(ctx, http.MethodGet, "/industry/packs/"+url.PathEscape(id), "", nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AssertionPropose proposes a knowledge assertion.
func (c *Client) AssertionPropose(ctx context.Context, a *KnowledgeAssertion) (*KnowledgeAssertion, error) {
	body, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	var out KnowledgeAssertion
	err = c.do(ctx, http.MethodPost, "/assertions", "application/json", body, &out)
	return &out, err
}

// AssertionList returns assertions about a subject entity.
func (c *Client) AssertionList(ctx context.Context, subject EntityRef) ([]*KnowledgeAssertion, error) {
	path := "/assertions?subjectNs=" + url.QueryEscape(subject.Namespace) + "&subjectId=" + url.QueryEscape(subject.CanonicalID)
	var out []*KnowledgeAssertion
	err := c.do(ctx, http.MethodGet, path, "", nil, &out)
	return out, err
}

// EntitySave registers a canonical entity reference.
func (c *Client) EntitySave(ctx context.Context, e *EntityRef) (*EntityRef, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	var out EntityRef
	err = c.do(ctx, http.MethodPost, "/entities", "application/json", body, &out)
	return &out, err
}

// EntityResolve resolves a hint to canonical entity references.
func (c *Client) EntityResolve(ctx context.Context, hint string) ([]*EntityRef, error) {
	body, _ := json.Marshal(map[string]string{"hint": hint})
	var out []*EntityRef
	err := c.do(ctx, http.MethodPost, "/entities/resolve", "application/json", body, &out)
	return out, err
}

// EntityExplain resolves a hint into a full Resolution (canonical ref, match
// evidence, related identity relations, gaps).
func (c *Client) EntityExplain(ctx context.Context, hint string) (*Resolution, error) {
	body, _ := json.Marshal(map[string]string{"hint": hint})
	var out Resolution
	err := c.do(ctx, http.MethodPost, "/entities/explain", "application/json", body, &out)
	return &out, err
}

// EntityMergeCandidate proposes a merge candidate relation with evidence.
// The relation is stored with authority "candidate" until confirmed.
func (c *Client) EntityMergeCandidate(ctx context.Context, left, right EntityRef, kind RelationKind, generatedBy string) (*Relation, *MatchEvidence, error) {
	body, _ := json.Marshal(map[string]any{
		"left": left, "right": right, "kind": kind, "generatedBy": generatedBy,
	})
	var out struct {
		Relation *Relation      `json:"relation"`
		Evidence *MatchEvidence `json:"evidence"`
	}
	err := c.do(ctx, http.MethodPost, "/entities/merge-candidate", "application/json", body, &out)
	return out.Relation, out.Evidence, err
}

// SemanticQuery runs a semantic query and returns the KnowledgeView, intent,
// relations, paths and evidence chains.
func (c *Client) SemanticQuery(ctx context.Context, query, principal, purpose string) (*SemanticResult, error) {
	body, _ := json.Marshal(map[string]any{
		"query": query, "principal": principal, "purpose": purpose,
	})
	var out SemanticResult
	err := c.do(ctx, http.MethodPost, "/query/semantic", "application/json", body, &out)
	return &out, err
}

// PathQuery finds a multi-hop path between two node keys.
func (c *Client) PathQuery(ctx context.Context, from, to string) (*Path, error) {
	body, _ := json.Marshal(map[string]string{"from": from, "to": to})
	var out Path
	err := c.do(ctx, http.MethodPost, "/query/path", "application/json", body, &out)
	return &out, err
}

// SubgraphQuery returns the 1-hop neighborhood of a node.
func (c *Client) SubgraphQuery(ctx context.Context, node string) ([]*Relation, error) {
	body, _ := json.Marshal(map[string]string{"node": node})
	var out []*Relation
	err := c.do(ctx, http.MethodPost, "/query/subgraph", "application/json", body, &out)
	return out, err
}

// GroundingAssemble assembles an evidence-backed grounding bundle.
func (c *Client) GroundingAssemble(ctx context.Context, query, principal, purpose string, budget int) (*GroundingBundle, error) {
	body, _ := json.Marshal(map[string]any{
		"query": query, "principal": principal, "purpose": purpose, "budget": budget,
	})
	var out GroundingBundle
	err := c.do(ctx, http.MethodPost, "/grounding/assemble", "application/json", body, &out)
	return &out, err
}

// GroundingExplain returns the citation-level explanation of an assembly.
func (c *Client) GroundingExplain(ctx context.Context, query, principal, purpose string) (*GroundingExplanation, error) {
	body, _ := json.Marshal(map[string]string{
		"query": query, "principal": principal, "purpose": purpose,
	})
	var out GroundingExplanation
	err := c.do(ctx, http.MethodPost, "/grounding/explain", "application/json", body, &out)
	return &out, err
}

// GroundingRedact re-assembles with sensitive evidence sources hidden.
func (c *Client) GroundingRedact(ctx context.Context, query, principal, purpose string, hideTags []string) (*GroundingBundle, error) {
	body, _ := json.Marshal(map[string]any{
		"query": query, "principal": principal, "purpose": purpose, "hideTags": hideTags,
	})
	var out GroundingBundle
	err := c.do(ctx, http.MethodPost, "/grounding/redact", "application/json", body, &out)
	return &out, err
}

func (c *Client) do(ctx context.Context, method, path, contentType string, body []byte, out any) error {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("gnosivela: %s %s -> %d: %s", method, path, resp.StatusCode, string(data))
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
