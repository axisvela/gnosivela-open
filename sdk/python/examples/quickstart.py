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

# 6. federated query across the registered domains (local by default)
fed = c.federation_query("ACME risk", "risk-officer", "onboarding")
print("federation domains:", fed["domainsQueried"], "hits:", fed.get("domainHits"))

# 7. real-time events: register a contract and ingest an event
c.event_contract_register({
    "id": "procurement.price.updated", "type": "price.updated", "ontology": "procurement.supplier@1.2",
    "templates": [{
        "predicate": "Supplier:price", "subjectField": "company", "subjectType": "Supplier",
        "objectField": "amount", "objectType": "number", "region": "SG", "validFor": "90d",
    }],
})
ing = c.event_ingest("procurement.price.updated", {
    "id": "e-1", "type": "price.updated", "source": "market-feed",
    "payload": {"company": "ACME", "amount": 12.5},
})
print("event ingest assertions:", len(ing["assertions"]), "resolved:", len(ing["resolved"]))

# 8. operations: quality snapshot + SLO incidents + metrics
q = c.quality("")
print("quality citation:", round(q["citationCompleteness"], 3), "conflicts:", q["conflicts"])
c.incident_rule_add({"id": "r-conflicts", "metric": "conflicts", "operator": ">=", "threshold": 0, "severity": "warning"})
opened = c.incident_check()["opened"]
print("incidents opened:", len(opened))
counts = c.metrics()["counts"]
print("metrics:", counts)
