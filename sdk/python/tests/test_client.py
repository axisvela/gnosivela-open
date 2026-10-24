"""Smoke tests against a local mock control plane (standard library only)."""

import http.server
import json
import threading
import unittest

from gnosivela import Client, GnosivelaError
from gnosivela.types import EntityRef, KnowledgeAssertion, Value


class MockHandler(http.server.BaseHTTPRequestHandler):
    def log_message(self, *args):  # silence
        pass

    def _reply(self, payload, status=200):
        body = json.dumps(payload).encode("utf-8")
        self.send_response(status)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def _reply_text(self, text, status=200, content_type="text/turtle"):
        body = text.encode("utf-8")
        self.send_response(status)
        self.send_header("Content-Type", content_type)
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def do_GET(self):
        if self.path == "/ontologies/acme/latest":
            self._reply({"namespace": "acme", "version": "1.2", "status": "published"})
        elif self.path == "/consistency/conflicts":
            self._reply({"conflicts": [{"subject": {"namespace": "m", "canonicalId": "E-1"}, "predicate": "p", "assertions": []}]})
        elif self.path == "/audit":
            self._reply({"entries": [{"id": "AUD-000001", "action": "publish", "actor": "alice"}]})
        elif self.path == "/events/contracts":
            self._reply({"contracts": [{"id": "c-1", "type": "price.updated"}]})
        elif self.path == "/quality":
            self._reply({"citationCompleteness": 0.98, "unresolvedRate": 0.0, "conflicts": 0})
        elif self.path == "/metrics":
            self._reply({"counts": {"query.semantic": 3}, "asOf": "2024-01-01T00:00:00Z"})
        elif self.path == "/bridge/procurement.supplier/contract":
            self._reply({"id": "procurement.supplier@1.2", "signature": "a" * 64, "concepts": []})
        elif self.path == "/metrics/definitions":
            self._reply({"definitions": [{"id": "m:1", "formula": "sum(x)"}]})
        elif self.path.startswith("/interop/shacl/"):
            self._reply_text('@prefix sh: <http://www.w3.org/ns/shacl#> .\nex:SupplierShape a sh:NodeShape .\n')
        elif self.path == "/actions":
            self._reply({"actions": [{"id": "notify", "kind": "webhook"}]})
        elif self.path == "/actions/runs":
            self._reply({"runs": [{"id": "ACT-1", "status": "succeeded"}]})
        elif self.path == "/connectors":
            self._reply({"connectors": [{"kind": "api"}, {"kind": "streaming"}]})
        elif self.path.startswith("/studio/impact"):
            self._reply({"namespace": "acme", "fromVersion": "1.0", "toVersion": "1.1",
                         "additive": [{"kind": "added", "type": "concept", "id": "Contract"}], "blocking": False})
        elif self.path.startswith("/studio/model"):
            self._reply({"namespace": "acme", "version": "1.2", "status": "published",
                         "concepts": [{"id": "Supplier", "assertionCount": 2}],
                         "relations": [], "sources": [{"source": "erp", "count": 2}]})
        else:
            self._reply({"error": "not found"}, 404)

    def do_POST(self):
        length = int(self.headers.get("Content-Length", 0))
        raw = self.rfile.read(length) if length else b"{}"
        try:
            body = json.loads(raw)
        except Exception:
            body = {}
        if self.path == "/entities":
            self._reply({"namespace": body.get("namespace"), "canonicalId": body.get("canonicalId")}, 201)
        elif self.path == "/ontologies":
            self._reply({"namespace": body.get("namespace"), "version": body.get("version"), "status": "draft"}, 201)
        elif self.path == "/assertions":
            self._reply({"assertionId": body.get("assertionId")}, 201)
        elif self.path == "/events/ingest":
            self._reply({"assertions": [{"assertionId": "ev:e-1:p"}], "resolved": ["ev:e-1:p"], "gaps": []}, 201)
        elif self.path == "/policy/evaluate":
            self._reply({"allowed": True, "reason": "open default"})
        elif self.path == "/query/hybrid":
            self._reply({
                "query": body.get("query"),
                "provider": "hashing/256",
                "fusion": "rrf",
                "hits": [{"id": "ka:risk-1", "kind": "assertion", "score": 0.5, "lexicalRank": 1}],
                "lexicalCandidates": 1,
            })
        elif self.path == "/federation/query":
            self._reply({"query": body.get("query"), "domainsQueried": 1, "domainHits": {"local": 1}})
        elif self.path == "/quality":
            self._reply({"citationCompleteness": 0.98, "unresolvedRate": 0.0, "conflicts": 0})
        elif self.path == "/metrics":
            self._reply({"counts": {"query.semantic": 3}, "asOf": "2024-01-01T00:00:00Z"})
        elif self.path == "/incidents/check":
            self._reply({"opened": [{"id": "INC-001", "status": "open"}], "quality": {}})
        elif self.path == "/bridge/procurement.supplier/contract":
            self._reply({"id": "procurement.supplier@1.2", "signature": "a" * 64, "concepts": []})
        elif self.path == "/metrics/definitions":
            if "POST" == self.command:
                self._reply({"id": "m:1", "formula": "sum(x)"}, 201)
            else:
                self._reply({"definitions": [{"id": "m:1", "formula": "sum(x)"}]})
        elif self.path == "/sparql":
            self._reply({"head": {"vars": ["s"]}, "results": {"bindings": [{"s": {"type": "uri", "value": "https://example/C-1"}}]}})
        elif self.path == "/actions":
            self._reply({"id": body.get("id"), "kind": body.get("kind")}, 201)
        elif self.path.startswith("/actions/") and self.path.endswith("/run"):
            self._reply({"id": "ACT-1", "actionId": self.path.split("/")[2], "status": "succeeded",
                         "principal": body.get("principal")})
        elif self.path == "/connectors/introspect":
            self._reply({"kind": body.get("kind"), "connector": body.get("kind"), "fields": [{"name": "eventId", "dataType": "string", "required": True}]})
        elif self.path == "/connectors/documents":
            self._reply({"id": "doc:" + body.get("docId", ""), "docId": body.get("docId"), "bindingState": "bound",
                         "chunks": [{"id": "c1", "docId": body.get("docId"), "content": body.get("content"), "locator": "p#c"}]}, 201)
        elif self.path == "/lineage/columns":
            schema = body.get("schema", {})
            mapping = body.get("mapping", {})
            self._reply({"sourceSchema": schema.get("id"), "mappingId": mapping.get("id"),
                         "mappingVersion": mapping.get("version"),
                         "columns": [{"sourceField": "supplier_id", "targetProperty": "supplierId", "status": "mapped"}],
                         "unmapped": ["active"]})
        else:
            self._reply({"error": "not found"}, 404)


class TestClient(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.server = http.server.ThreadingHTTPServer(("127.0.0.1", 0), MockHandler)
        cls.thread = threading.Thread(target=cls.server.serve_forever, daemon=True)
        cls.thread.start()
        cls.base = "http://127.0.0.1:%d" % cls.server.server_address[1]

    @classmethod
    def tearDownClass(cls):
        cls.server.shutdown()

    def test_ontology_latest(self):
        c = Client(self.base)
        o = c.ontology_latest("acme")
        self.assertEqual(o["version"], "1.2")

    def test_ontology_create_json(self):
        c = Client(self.base)
        o = c.ontology_create_json({"namespace": "acme", "version": "2.0"})
        self.assertEqual(o["namespace"], "acme")
        self.assertEqual(o["version"], "2.0")

    def test_entity_and_assertion_roundtrip(self):
        c = Client(self.base)
        ent = c.entity_save(EntityRef(namespace="mdm", canonical_id="C-1", type="Company"))
        self.assertEqual(ent["canonicalId"], "C-1")
        a = c.assertion_propose(
            KnowledgeAssertion(
                assertion_id="ka-1",
                subject=EntityRef(namespace="mdm", canonical_id="C-1", type="Company"),
                predicate="company:status",
                object=Value(type="string", string="active"),
                source="erp",
            )
        )
        self.assertEqual(a["assertionId"], "ka-1")

    def test_governance_endpoints(self):
        c = Client(self.base)
        conflicts = c.consistency_conflicts()
        self.assertEqual(len(conflicts["conflicts"]), 1)
        audit = c.audit_list()
        self.assertEqual(audit["entries"][0]["action"], "publish")
        d = c.policy_evaluate({"action": "ontology.publish", "resource": "x"})
        self.assertTrue(d["allowed"])

    def test_events(self):
        c = Client(self.base)
        contracts = c.event_contract_list()
        self.assertEqual(len(contracts["contracts"]), 1)
        ing = c.event_ingest("c-1", {"id": "e-1", "type": "price.updated", "source": "s", "payload": {}})
        self.assertEqual(len(ing["resolved"]), 1)

    def test_federation(self):
        c = Client(self.base)
        view = c.federation_query("ACME", "auditor", "review")
        self.assertEqual(view["domainsQueried"], 1)

    def test_hybrid_query(self):
        c = Client(self.base)
        result = c.hybrid_query("supplier risk", limit=5, weights={"lexical": 1, "vector": 1, "graph": 0.5})
        self.assertEqual(result["fusion"], "rrf")
        self.assertEqual(result["hits"][0]["id"], "ka:risk-1")

    def test_studio_model(self):
        c = Client(self.base)
        model = c.studio_model("acme")
        self.assertEqual(model["namespace"], "acme")
        self.assertEqual(model["concepts"][0]["assertionCount"], 2)

    def test_studio_impact(self):
        c = Client(self.base)
        report = c.studio_impact("acme")
        self.assertEqual(report["fromVersion"], "1.0")
        self.assertFalse(report["blocking"])
        self.assertEqual(report["additive"][0]["id"], "Contract")

    def test_actions(self):
        c = Client(self.base)
        self.assertEqual(c.actions()["actions"][0]["id"], "notify")
        created = c.action_register({"id": "notify", "kind": "webhook", "url": "https://example.com"})
        self.assertEqual(created["id"], "notify")
        run = c.action_run("notify", principal="alice", input_data={"x": 1})
        self.assertEqual(run["status"], "succeeded")
        self.assertEqual(run["principal"], "alice")
        self.assertEqual(c.action_runs()["runs"][0]["id"], "ACT-1")

    def test_sparql_query(self):
        c = Client(self.base)
        result = c.sparql_query("SELECT ?s WHERE { ?s ?p ?o }")
        self.assertEqual(result["head"]["vars"], ["s"])

    def test_shacl_export(self):
        c = Client(self.base)
        text = c.shacl_export("acme", "1.2")
        self.assertIn("ex:SupplierShape", text)

    def test_connectors(self):
        c = Client(self.base)
        kinds = [x["kind"] for x in c.connectors()["connectors"]]
        self.assertIn("streaming", kinds)
        schema = c.connector_introspect("streaming", {"options": {"inlineSchema": "{}"}})
        lineage = c.column_lineage(
            {"id": "glue://catalog/suppliers", "kind": "glue", "name": "suppliers", "connector": "glue",
             "fields": [{"name": "supplier_id", "dataType": "integer"}]},
            {"id": "MAP-1", "version": "1.0", "rules": []},
        )
        self.assertEqual(lineage["sourceSchema"], "glue://catalog/suppliers")
        self.assertEqual(lineage["unmapped"], ["active"])
        self.assertEqual(schema["kind"], "streaming")
        self.assertEqual(schema["fields"][0]["name"], "eventId")
        binding = c.connector_bind_document("p1", "Suppliers must be reviewed.", extract_claims=True)
        self.assertEqual(binding["docId"], "p1")
        self.assertEqual(len(binding["chunks"]), 1)

    def test_operations_endpoints(self):
        c = Client(self.base)
        q = c.quality("")
        self.assertGreater(q["citationCompleteness"], 0.9)
        counts = c.metrics()["counts"]
        self.assertEqual(counts["query.semantic"], 3)
        opened = c.incident_check()["opened"]
        self.assertEqual(opened[0]["status"], "open")
        contract = c.bridge_contract_export("procurement.supplier")
        self.assertEqual(len(contract["signature"]), 64)
        c.metric_definition_register({"id": "m:1", "formula": "sum(x)"})
        defs = c.metric_definitions()["definitions"]
        self.assertEqual(defs[0]["id"], "m:1")

    def test_error(self):
        c = Client(self.base)
        with self.assertRaises(GnosivelaError) as cm:
            c.ontology_latest("missing")
        self.assertEqual(cm.exception.status, 404)


if __name__ == "__main__":
    unittest.main()
