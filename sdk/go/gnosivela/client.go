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
