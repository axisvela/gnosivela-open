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
	doc, err := os.ReadFile("examples/supplier.dsl")
	if err != nil {
		log.Fatalf("read dsl: %v", err)
	}
	res, err := c.OntologyCreate(ctx, string(doc))
	if err != nil {
		log.Fatalf("create ontology: %v", err)
	}
	fmt.Printf("ontology %s@%s status=%s concepts=%d\n",
		res.Ontology.Namespace, res.Ontology.Version, res.Ontology.Status, len(res.Ontology.Concepts))

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
}

func str(s string) *string { return &s }
