package gnosivela

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOntologyCreateAndDiff(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/ontologies":
			writeJSON(w, http.StatusCreated, OntologyCreateResult{
				Ontology: &Ontology{Namespace: "procurement.supplier", Version: "1.2", Status: OntologyDraft},
			})
		case r.Method == http.MethodGet && r.URL.Path == "/ontologies/procurement.supplier/latest":
			writeJSON(w, http.StatusOK, &Ontology{Namespace: "procurement.supplier", Version: "1.2", Status: OntologyDraft})
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/diff"):
			writeJSON(w, http.StatusOK, map[string]any{
				"changes": []Change{{Kind: "added", Type: "concept", ID: "LegalEntity"}},
			})
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
	defer ts.Close()

	ctx := context.Background()
	c := New(ts.URL)

	res, err := c.OntologyCreate(ctx, "ontology procurement.supplier @version 1.2\n")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if res.Ontology.Namespace != "procurement.supplier" {
		t.Errorf("namespace = %q", res.Ontology.Namespace)
	}

	latest, err := c.OntologyLatest(ctx, "procurement.supplier")
	if err != nil {
		t.Fatalf("latest: %v", err)
	}
	if latest.Version != "1.2" {
		t.Errorf("latest version = %q", latest.Version)
	}

	changes, err := c.OntologyDiff(ctx, "procurement.supplier", "1.1", "1.2")
	if err != nil {
		t.Fatalf("diff: %v", err)
	}
	if len(changes) != 1 || changes[0].ID != "LegalEntity" {
		t.Errorf("changes = %+v", changes)
	}
}

func TestAssertionAndResolve(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/assertions":
			var a KnowledgeAssertion
			_ = json.NewDecoder(r.Body).Decode(&a)
			a.AssertionID = "ka:test"
			a.Status = AssertionProposed
			writeJSON(w, http.StatusCreated, a)
		case r.Method == http.MethodPost && r.URL.Path == "/entities/resolve":
			writeJSON(w, http.StatusOK, []*EntityRef{{Namespace: "mdm", CanonicalID: "C-1042", Aliases: []string{"ACME"}}})
		case r.Method == http.MethodPost && r.URL.Path == "/entities/explain":
			writeJSON(w, http.StatusOK, &Resolution{
				Hint:      "ACME",
				Canonical: &EntityRef{Namespace: "mdm", CanonicalID: "C-1042"},
				Matches:   []ResolvedMatch{{Ref: EntityRef{Namespace: "mdm", CanonicalID: "C-1042"}, Score: 1, Status: "exact"}},
			})
		case r.Method == http.MethodPost && r.URL.Path == "/entities/merge-candidate":
			writeJSON(w, http.StatusCreated, map[string]any{
				"relation": &Relation{ID: "rel:1", Kind: RelationExactSameAs, Authority: "candidate"},
				"evidence": &MatchEvidence{Status: "auto", GeneratedBy: "rule", Score: 1},
			})
		case r.Method == http.MethodPost && r.URL.Path == "/grounding/assemble":
			writeJSON(w, http.StatusOK, &GroundingBundle{
				ResolvedIntent:   "entities=[mdm:C-1042]",
				PolicyDecision:   "allow tags=[procurement]",
				SemanticContract: &SemanticContract{OntologyVersion: "procurement.supplier@1.2"},
			})
		case r.Method == http.MethodPost && r.URL.Path == "/grounding/explain":
			writeJSON(w, http.StatusOK, &GroundingExplanation{
				Intent:               "entities=[mdm:C-1042]",
				CitationCompleteness: 1.0,
			})
		case r.Method == http.MethodPost && r.URL.Path == "/grounding/redact":
			writeJSON(w, http.StatusOK, &GroundingBundle{
				ResolvedIntent: "entities=[mdm:C-1042]",
				Knowledge:      []*KnowledgeAssertion{{AssertionID: "ka:g1"}},
			})
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
	defer ts.Close()

	ctx := context.Background()
	c := New(ts.URL)

	a, err := c.AssertionPropose(ctx, &KnowledgeAssertion{
		Subject:   EntityRef{Namespace: "mdm", CanonicalID: "C-1042"},
		Predicate: "risk:hasCreditRating",
	})
	if err != nil {
		t.Fatalf("propose: %v", err)
	}
	if a.AssertionID != "ka:test" || a.Status != AssertionProposed {
		t.Errorf("assertion = %+v", a)
	}

	matches, err := c.EntityResolve(ctx, "ACME")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if len(matches) != 1 || matches[0].CanonicalID != "C-1042" {
		t.Errorf("matches = %+v", matches)
	}

	expl, err := c.EntityExplain(ctx, "ACME")
	if err != nil {
		t.Fatalf("explain: %v", err)
	}
	if expl.Canonical == nil || expl.Canonical.CanonicalID != "C-1042" {
		t.Errorf("explain canonical = %+v", expl.Canonical)
	}

	rel, ev, err := c.EntityMergeCandidate(ctx,
		EntityRef{Namespace: "crm", CanonicalID: "ACME"},
		EntityRef{Namespace: "mdm", CanonicalID: "C-1042"},
		RelationExactSameAs, "rule")
	if err != nil {
		t.Fatalf("merge-candidate: %v", err)
	}
	if rel == nil || rel.Authority != "candidate" || ev.Status != "auto" {
		t.Errorf("merge candidate = %+v, evidence = %+v", rel, ev)
	}

	bundle, err := c.GroundingAssemble(ctx, "ACME risk", "alice", "onboarding", 0)
	if err != nil {
		t.Fatalf("grounding assemble: %v", err)
	}
	if bundle.ResolvedIntent != "entities=[mdm:C-1042]" || bundle.PolicyDecision != "allow tags=[procurement]" {
		t.Errorf("bundle = %+v", bundle)
	}

	ex, err := c.GroundingExplain(ctx, "ACME risk", "alice", "onboarding")
	if err != nil {
		t.Fatalf("grounding explain: %v", err)
	}
	if ex.CitationCompleteness != 1.0 {
		t.Errorf("explain = %+v", ex)
	}

	rb, err := c.GroundingRedact(ctx, "ACME risk", "alice", "onboarding", []string{"procurement"})
	if err != nil {
		t.Fatalf("grounding redact: %v", err)
	}
	if len(rb.Knowledge) != 1 || rb.Knowledge[0].AssertionID != "ka:g1" {
		t.Errorf("redacted bundle = %+v", rb)
	}
}

func TestHybridQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/query/hybrid" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		var body struct {
			Query   string        `json:"query"`
			Limit   int           `json:"limit"`
			Weights HybridWeights `json:"weights"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, HybridResult{
			Query: body.Query, Provider: "hashing/256", Fusion: "rrf",
			Weights: body.Weights,
			Hits:    []HybridHit{{ID: "ka:risk-1", Kind: "assertion", Score: 0.5, LexicalRank: 1}},
		})
	}))
	defer ts.Close()

	c := New(ts.URL)
	res, err := c.HybridQuery(context.Background(), "supplier risk", 5, &HybridWeights{Lexical: 1, Vector: 1, Graph: 0.5})
	if err != nil {
		t.Fatalf("hybrid query: %v", err)
	}
	if res.Fusion != "rrf" || len(res.Hits) != 1 || res.Hits[0].ID != "ka:risk-1" {
		t.Errorf("result = %+v", res)
	}
	if res.Weights.Graph != 0.5 {
		t.Errorf("weights = %+v", res.Weights)
	}
}

func TestStudioModel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/studio/model" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if r.URL.Query().Get("namespace") != "acme" {
			http.Error(w, "bad namespace", http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, StudioModel{
			Namespace: "acme", Version: "1.2", Status: "published",
			Concepts:  []StudioConcept{{ID: "Supplier", AssertionCount: 2}},
			Relations: []StudioRelation{{ID: "supplies"}},
			Sources:   []StudioSource{{Source: "erp", Count: 2}},
		})
	}))
	defer ts.Close()
	model, err := New(ts.URL).StudioModel(context.Background(), "acme")
	if err != nil {
		t.Fatalf("studio model: %v", err)
	}
	if model.Namespace != "acme" || len(model.Concepts) != 1 || model.Concepts[0].AssertionCount != 2 {
		t.Errorf("model = %+v", model)
	}
}

func TestActionMethods(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/actions":
			writeJSON(w, http.StatusOK, map[string]any{"actions": []ActionDefinition{{ID: "notify", Kind: "webhook", URL: "https://example.com"}}})
		case r.Method == http.MethodPost && r.URL.Path == "/actions":
			var def ActionDefinition
			if err := json.NewDecoder(r.Body).Decode(&def); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			writeJSON(w, http.StatusCreated, def)
		case r.Method == http.MethodPost && r.URL.Path == "/actions/notify/run":
			var req ActionRunRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			writeJSON(w, http.StatusOK, ActionRun{ID: "ACT-1", ActionID: "notify", Status: "succeeded", Principal: req.Principal})
		case r.Method == http.MethodGet && r.URL.Path == "/actions/runs":
			writeJSON(w, http.StatusOK, map[string]any{"runs": []ActionRun{{ID: "ACT-1", ActionID: "notify", Status: "succeeded"}}})
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
	defer ts.Close()
	ctx := context.Background()
	c := New(ts.URL)
	list, err := c.Actions(ctx)
	if err != nil || len(list) != 1 {
		t.Fatalf("actions = %+v, %v", list, err)
	}
	def, err := c.ActionRegister(ctx, ActionDefinition{ID: "notify", Kind: "webhook", URL: "https://example.com"})
	if err != nil || def.ID != "notify" {
		t.Fatalf("register = %+v, %v", def, err)
	}
	run, err := c.ActionRun(ctx, "notify", ActionRunRequest{Principal: "alice", Input: map[string]any{"x": 1}})
	if err != nil || run.Status != "succeeded" || run.Principal != "alice" {
		t.Fatalf("run = %+v, %v", run, err)
	}
	runs, err := c.ActionRuns(ctx)
	if err != nil || len(runs) != 1 {
		t.Fatalf("runs = %+v, %v", runs, err)
	}
}

func TestSparqlQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/sparql" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		var body struct {
			Query string `json:"query"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"head":    map[string]any{"vars": []string{"s"}},
			"results": map[string]any{"bindings": []any{map[string]any{"s": map[string]any{"type": "uri", "value": "https://example/C-1"}}}},
		})
	}))
	defer ts.Close()
	out, err := New(ts.URL).SparqlQuery(context.Background(), `SELECT ?s WHERE { ?s ?p ?o }`)
	if err != nil {
		t.Fatalf("sparql: %v", err)
	}
	head, _ := out["head"].(map[string]any)
	if head == nil {
		t.Fatalf("result = %+v", out)
	}
}

func TestShaclExport(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/interop/shacl/acme" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if r.URL.Query().Get("version") != "1.2" {
			http.Error(w, "version", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/turtle")
		_, _ = w.Write([]byte("@prefix sh: <http://www.w3.org/ns/shacl#> .\nex:SupplierShape a sh:NodeShape .\n"))
	}))
	defer ts.Close()
	out, err := New(ts.URL).ShaclExport(context.Background(), "acme", "1.2")
	if err != nil {
		t.Fatalf("shacl export: %v", err)
	}
	if !strings.Contains(out, "ex:SupplierShape") {
		t.Errorf("output = %q", out)
	}
}

func TestStudioImpact(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/studio/impact" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, ImpactReport{
			Namespace: "acme", FromVersion: "1.0", ToVersion: "1.1",
			Additive: []Change{{Kind: "added", Type: "concept", ID: "Contract"}},
		})
	}))
	defer ts.Close()
	report, err := New(ts.URL).StudioImpact(context.Background(), "acme")
	if err != nil {
		t.Fatalf("studio impact: %v", err)
	}
	if report.FromVersion != "1.0" || len(report.Additive) != 1 {
		t.Errorf("report = %+v", report)
	}
}

func TestConnectorMethods(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/connectors":
			writeJSON(w, http.StatusOK, map[string]any{"connectors": []ConnectorInfo{{Kind: "api"}, {Kind: "streaming"}}})
		case r.Method == http.MethodPost && r.URL.Path == "/connectors/introspect":
			var req ConnectorIntrospectRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			writeJSON(w, http.StatusOK, SourceSchema{
				Kind: req.Kind, Connector: req.Kind,
				Fields: []SourceField{{Name: "eventId", DataType: "string", Required: true}},
			})
		case r.Method == http.MethodPost && r.URL.Path == "/connectors/documents":
			var req DocumentBindingRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			writeJSON(w, http.StatusCreated, DocumentBinding{
				ID: "doc:" + req.DocID, DocID: req.DocID, BindingState: "bound",
				ArtifactDigest: "abc123",
				Chunks:         []DocumentChunk{{ID: "c1", DocID: req.DocID, Content: req.Content, Locator: "p#c"}},
			})
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
	defer ts.Close()

	ctx := context.Background()
	c := New(ts.URL)
	list, err := c.Connectors(ctx)
	if err != nil {
		t.Fatalf("connectors: %v", err)
	}
	if len(list) != 2 || list[0].Kind != "api" {
		t.Errorf("connectors = %+v", list)
	}
	schema, err := c.ConnectorIntrospect(ctx, "streaming", ConnectorSpec{Options: map[string]string{"inlineSchema": "{}"}})
	if err != nil {
		t.Fatalf("introspect: %v", err)
	}
	if schema.Kind != "streaming" || len(schema.Fields) != 1 {
		t.Errorf("schema = %+v", schema)
	}
	binding, err := c.ConnectorBindDocument(ctx, DocumentBindingRequest{DocID: "p1", Content: "must review"})
	if err != nil {
		t.Fatalf("bind: %v", err)
	}
	if binding.ID != "doc:p1" || len(binding.Chunks) != 1 || binding.ArtifactDigest != "abc123" {
		t.Errorf("binding = %+v", binding)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
