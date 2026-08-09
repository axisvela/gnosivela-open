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
				Hint: "ACME",
				Canonical: &EntityRef{Namespace: "mdm", CanonicalID: "C-1042"},
				Matches: []ResolvedMatch{{Ref: EntityRef{Namespace: "mdm", CanonicalID: "C-1042"}, Score: 1, Status: "exact"}},
			})
		case r.Method == http.MethodPost && r.URL.Path == "/entities/merge-candidate":
			writeJSON(w, http.StatusCreated, map[string]any{
				"relation": &Relation{ID: "rel:1", Kind: RelationExactSameAs, Authority: "candidate"},
				"evidence": &MatchEvidence{Status: "auto", GeneratedBy: "rule", Score: 1},
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
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
