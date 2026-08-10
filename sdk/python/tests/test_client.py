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

    def do_GET(self):
        if self.path == "/ontologies/acme/latest":
            self._reply({"namespace": "acme", "version": "1.2", "status": "published"})
        elif self.path == "/consistency/conflicts":
            self._reply({"conflicts": [{"subject": {"namespace": "m", "canonicalId": "E-1"}, "predicate": "p", "assertions": []}]})
        elif self.path == "/audit":
            self._reply({"entries": [{"id": "AUD-000001", "action": "publish", "actor": "alice"}]})
        elif self.path == "/events/contracts":
            self._reply({"contracts": [{"id": "c-1", "type": "price.updated"}]})
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
        elif self.path == "/assertions":
            self._reply({"assertionId": body.get("assertionId")}, 201)
        elif self.path == "/events/ingest":
            self._reply({"assertions": [{"assertionId": "ev:e-1:p"}], "resolved": ["ev:e-1:p"], "gaps": []}, 201)
        elif self.path == "/policy/evaluate":
            self._reply({"allowed": True, "reason": "open default"})
        elif self.path == "/federation/query":
            self._reply({"query": body.get("query"), "domainsQueried": 1, "domainHits": {"local": 1}})
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

    def test_error(self):
        c = Client(self.base)
        with self.assertRaises(GnosivelaError) as cm:
            c.ontology_latest("missing")
        self.assertEqual(cm.exception.status, 404)


if __name__ == "__main__":
    unittest.main()
