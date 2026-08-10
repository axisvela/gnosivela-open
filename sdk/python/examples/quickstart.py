"""GNOSIVELA Python quickstart.

Assumes a control plane running on http://localhost:8080 (in-memory mode is
fine): `cd backend && go run ./cmd/gnosivela`.
"""

from gnosivela import Client
from gnosivela.types import EntityRef, KnowledgeAssertion, Value

c = Client("http://localhost:8080")

# 1. create an ontology from the semantic contract DSL
dsl = """ontology procurement.supplier @version 1.0
entity Supplier identifiedBy canonicalSupplierId
property riskScore: Decimal [0..1]
"""
result = c.ontology_create(dsl)
print("ontology:", result["ontology"]["namespace"], result["ontology"]["version"])

# 2. register a canonical entity and resolve a hint
ent = c.entity_save(EntityRef(namespace="mdm", canonical_id="S-1042", type="Supplier", aliases=["ACME"]))
print("entity:", ent["canonicalId"])

# 3. propose a sourced assertion
a = KnowledgeAssertion(
    assertion_id="ka:risk",
    subject=EntityRef(namespace="mdm", canonical_id="S-1042", type="Supplier"),
    predicate="risk:score",
    object=Value(type="number", number=0.82),
    source="RiskOffice",
    status="validated",
    confidence=0.9,
)
print("assertion:", c.assertion_propose(a)["assertionId"])

# 4. run the consistency report and surface conflicts
report = c.consistency_report("procurement.supplier")
print("consistency failures:", report["failures"])

# 5. contract-driven query through the semantic bridge
view = c.bridge_query("procurement.supplier", "ACME risk", "risk-officer", "onboarding")
print("bridge view assertions:", len(view.get("assertions", [])))
