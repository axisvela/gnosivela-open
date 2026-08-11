// Command quickstart demonstrates the GNOSIVELA Go SDK against a running
// core server. Start the server first:
//
//	cd GNOSIVELA && go run ./cmd/gnosivela
//
// Then run this example from the GNOSIVELA-open repo:
//
//	go run ./examples/quickstart
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	gnosivela "github.com/axisrobo/GNOSIVELA-open/sdk/go/gnosivela"
)

func main() {
	base := os.Getenv("GNOSIVELA_URL")
	if base == "" {
		base = "http://localhost:8080"
	}
	c := gnosivela.New(base)
	ctx := context.Background()

	// 1. Create an ontology from a Semantic Contract DSL document.
	doc, err := os.ReadFile("contracts/examples/supplier.dsl")
	if err != nil {
		log.Fatalf("read dsl: %v", err)
	}
	res, err := c.OntologyCreate(ctx, string(doc))
	if err != nil {
		log.Fatalf("create ontology: %v", err)
	}
	fmt.Printf("ontology %s@%s status=%s concepts=%d\n",
		res.Ontology.Namespace, res.Ontology.Version, res.Ontology.Status, len(res.Ontology.Concepts))

	// 1b. Publish the ontology so contract export / bridge queries work.
	if err := c.OntologyPublish(ctx, res.Ontology.Namespace, res.Ontology.Version); err != nil {
		log.Fatalf("publish ontology: %v", err)
	}
	fmt.Printf("ontology published\n")

	// 2. Register a canonical entity with aliases.
	_, err = c.EntitySave(ctx, &gnosivela.EntityRef{
		Namespace: "mdm", CanonicalID: "C-1042", Type: "Supplier",
		Aliases: []string{"ACME", "Acme Pte. Ltd."}, Authority: "mdm",
	})
	if err != nil {
		log.Fatalf("save entity: %v", err)
	}

	// 3. Resolve a hint to the canonical entity.
	matches, err := c.EntityResolve(ctx, "ACME")
	if err != nil {
		log.Fatalf("resolve: %v", err)
	}
	fmt.Printf("resolve 'ACME' -> %d match(es):\n", len(matches))
	for _, m := range matches {
		fmt.Printf("  %s:%s (type=%s, authority=%s)\n", m.Namespace, m.CanonicalID, m.Type, m.Authority)
	}

	// 4. Propose a knowledge assertion bound to that entity.
	a, err := c.AssertionPropose(ctx, &gnosivela.KnowledgeAssertion{
		Subject:         gnosivela.EntityRef{Namespace: "mdm", CanonicalID: "C-1042", Type: "Supplier"},
		Predicate:       "risk:hasCreditRating",
		Object:          gnosivela.Value{Type: "string", String: str("A")},
		Context:         gnosivela.AssertionContext{Region: "SG", Purpose: "supplier onboarding"},
		Source:          "ERP",
		Authority:       "RiskOffice",
		Confidence:      0.98,
		OntologyVersion: "procurement.supplier@1.2",
	})
	if err != nil {
		log.Fatalf("propose assertion: %v", err)
	}
	fmt.Printf("assertion %s status=%s (predicate=%s)\n", a.AssertionID, a.Status, a.Predicate)

	// 5. Governance: consistency + conflicts + audit.
	report, err := c.ConsistencyReport(ctx, "procurement.supplier")
	if err != nil {
		log.Fatalf("consistency report: %v", err)
	}
	fmt.Printf("consistency failures=%d warnings=%d\n", report.Failures, report.Warnings)

	conflicts, err := c.ConsistencyConflicts(ctx)
	if err != nil {
		log.Fatalf("consistency conflicts: %v", err)
	}
	fmt.Printf("conflicts=%d\n", len(conflicts))

	audit, err := c.AuditList(ctx)
	if err != nil {
		log.Fatalf("audit list: %v", err)
	}
	fmt.Printf("audit entries=%d\n", len(audit))

	// 6. Semantic bridge: export the signed contract and query through it.
	contract, err := c.BridgeContractExport(ctx, "procurement.supplier")
	if err != nil {
		log.Fatalf("bridge export: %v", err)
	}
	fmt.Printf("contract %s signature=%s… jsonSchema=%d chars\n",
		contract.ID, contract.Signature[:8], len(contract.JSONSchema))

	view, err := c.BridgeQuery(ctx, "procurement.supplier", "ACME credit", "risk-officer", "onboarding")
	if err != nil {
		log.Fatalf("bridge query: %v", err)
	}
	fmt.Printf("contract-driven view assertions=%d\n", len(view.Assertions))

	// 7. Federation: query across the registered domain (local by default).
	fedView, err := c.FederationQuery(ctx, "ACME credit", "risk-officer", "onboarding")
	if err != nil {
		log.Fatalf("federation query: %v", err)
	}
	fmt.Printf("federation domains=%d hits=%d latencyMs=%d\n",
		fedView.DomainsQueried, len(fedView.DomainHits), fedView.LatencyMillis)

	// 8. Real-time events: register a contract and ingest an event.
	err = c.EventContractRegister(ctx, gnosivela.EventContract{
		ID: "procurement.price.updated", Type: "price.updated", Ontology: "procurement.supplier@1.2",
		Templates: []gnosivela.EventTemplate{{
			Predicate: "Supplier:price", SubjectField: "company", SubjectType: "Supplier",
			ObjectField: "amount", ObjectType: "number", Region: "SG", ValidFor: "90d",
		}},
	})
	if err != nil {
		log.Fatalf("event contract register: %v", err)
	}
	ing, err := c.EventIngest(ctx, "procurement.price.updated", &gnosivela.Event{
		ID: "e-1", Type: "price.updated", Source: "market-feed",
		Payload: map[string]any{"company": "ACME", "amount": 12.5},
	})
	if err != nil {
		log.Fatalf("event ingest: %v", err)
	}
	fmt.Printf("event ingest assertions=%d resolved=%d gaps=%d\n",
		len(ing.Assertions), len(ing.Resolved), len(ing.Gaps))

	// 9. Operations: quality snapshot + SLO incidents + metrics.
	q, err := c.Quality(ctx, "")
	if err != nil {
		log.Fatalf("quality: %v", err)
	}
	fmt.Printf("quality citation=%.3f unresolved=%.3f conflicts=%d\n",
		q.CitationCompleteness, q.UnresolvedRate, q.Conflicts)

	err = c.IncidentRuleAdd(ctx, gnosivela.IncidentRule{
		ID: "r-conflicts", Metric: "conflicts", Operator: ">=", Threshold: 0, Severity: "warning",
	})
	if err != nil {
		log.Fatalf("incident rule add: %v", err)
	}
	opened, _, err := c.IncidentCheck(ctx)
	if err != nil {
		log.Fatalf("incident check: %v", err)
	}
	fmt.Printf("incidents opened=%d\n", len(opened))

	metrics, err := c.Metrics(ctx)
	if err != nil {
		log.Fatalf("metrics: %v", err)
	}
	fmt.Printf("metrics=%v\n", metrics)
}

func str(s string) *string { return &s }
