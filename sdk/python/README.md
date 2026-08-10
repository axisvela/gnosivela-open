# GNOSIVELA Python SDK

Apache-2.0. Talks to the GNOSIVELA Semantic Control Plane over HTTP (never
links the AGPL core). Standard-library only — no third-party dependencies.

## Install

```bash
pip install -e ./sdk/python
```

## Usage

```python
from gnosivela import Client
from gnosivela.types import EntityRef, KnowledgeAssertion, Value

c = Client("http://localhost:8080")

c.ontology_create("ontology procurement.supplier @version 1.0\nentity Supplier identifiedBy canonicalSupplierId")
c.entity_save(EntityRef(namespace="mdm", canonical_id="S-1042", type="Supplier", aliases=["ACME"]))
c.assertion_propose(KnowledgeAssertion(
    assertion_id="ka:risk", subject=EntityRef(namespace="mdm", canonical_id="S-1042", type="Supplier"),
    predicate="risk:score", object=Value(type="number", number=0.82),
    source="RiskOffice", status="validated",
))
report = c.consistency_report("procurement.supplier")
```

See `examples/quickstart.py` for the full walkthrough.

## Tests

```bash
python -m unittest discover -s sdk/python/tests -v
```
