"""HTTP client for the GNOSIVELA Semantic Control Plane API.

Standard-library only (urllib), so the SDK has no third-party dependencies.
"""

import json
import urllib.parse
import urllib.request
from typing import Any, Dict, Optional

from .types import EntityRef, KnowledgeAssertion


class GnosivelaError(Exception):
    """Raised when the control plane returns a non-2xx response."""

    def __init__(self, status: int, body: str, path: str):
        self.status = status
        self.body = body
        self.path = path
        super().__init__("gnosivela: %s -> %d: %s" % (path, status, body))


class Client:
    """Thin client for the GNOSIVELA control plane."""

    def __init__(self, base_url: str = "http://localhost:8080"):
        self.base_url = base_url.rstrip("/")

    # ---- low level ----

    def _request(
        self,
        method: str,
        path: str,
        body: Optional[Any] = None,
        content_type: str = "application/json",
        headers: Optional[Dict[str, str]] = None,
    ) -> Any:
        url = self.base_url + path
        data = None
        req_headers = dict(headers or {})
        if body is not None:
            if isinstance(body, (dict, list)):
                data = json.dumps(body).encode("utf-8")
            elif isinstance(body, str):
                data = body.encode("utf-8")
                content_type = "text/plain"
            else:
                data = body
            req_headers["Content-Type"] = content_type
        req = urllib.request.Request(url, data=data, method=method, headers=req_headers)
        try:
            with urllib.request.urlopen(req) as resp:
                raw = resp.read().decode("utf-8")
                if not raw:
                    return None
                return json.loads(raw)
        except urllib.error.HTTPError as e:
            raise GnosivelaError(e.code, e.read().decode("utf-8"), path)

    # ---- ontology ----

    def ontology_create(self, dsl: str) -> Dict[str, Any]:
        return self._request("POST", "/ontologies", dsl)

    def ontology_latest(self, namespace: str) -> Dict[str, Any]:
        return self._request("GET", "/ontologies/%s/latest" % urllib.parse.quote(namespace))

    def ontology_get(self, namespace: str, version: str) -> Dict[str, Any]:
        return self._request(
            "GET", "/ontologies/%s/versions/%s" % (urllib.parse.quote(namespace), urllib.parse.quote(version))
        )

    def ontology_publish(self, namespace: str, version: str, approval: str = "") -> Dict[str, Any]:
        path = "/ontologies/%s/versions/%s/publish" % (urllib.parse.quote(namespace), urllib.parse.quote(version))
        if approval:
            path += "?approval=" + urllib.parse.quote(approval)
        return self._request("POST", path)

    def ontology_impact(self, namespace: str, version: str) -> Dict[str, Any]:
        return self._request("GET", "/ontologies/%s/versions/%s/impact" % (urllib.parse.quote(namespace), urllib.parse.quote(version)))

    def ontology_rollback(self, namespace: str, version: str) -> Dict[str, Any]:
        return self._request("POST", "/ontologies/%s/versions/%s/rollback" % (urllib.parse.quote(namespace), urllib.parse.quote(version)))

    def ontology_diff(self, namespace: str, version: str, other: str) -> Dict[str, Any]:
        return self._request(
            "GET",
            "/ontologies/%s/versions/%s/diff?other=%s"
            % (urllib.parse.quote(namespace), urllib.parse.quote(version), urllib.parse.quote(other)),
        )

    # ---- assertion ----

    def assertion_propose(self, a: KnowledgeAssertion) -> Dict[str, Any]:
        return self._request("POST", "/assertions", a.to_dict())

    def assertion_list(self, subject: EntityRef) -> Dict[str, Any]:
        return self._request(
            "GET",
            "/assertions?subjectNs=%s&subjectId=%s"
            % (urllib.parse.quote(subject.namespace), urllib.parse.quote(subject.canonical_id)),
        )

    # ---- entity ----

    def entity_save(self, e: EntityRef) -> Dict[str, Any]:
        return self._request("POST", "/entities", e.to_dict())

    def entity_resolve(self, hint: str) -> Dict[str, Any]:
        return self._request("POST", "/entities/resolve", {"hint": hint})

    def entity_explain(self, hint: str) -> Dict[str, Any]:
        return self._request("POST", "/entities/explain", {"hint": hint})

    # ---- query / grounding ----

    def semantic_query(self, query: str, principal: str = "", purpose: str = "") -> Dict[str, Any]:
        return self._request("POST", "/query/semantic", {"query": query, "principal": principal, "purpose": purpose})

    def path_query(self, frm: str, to: str) -> Dict[str, Any]:
        return self._request("POST", "/query/path", {"from": frm, "to": to})

    def subgraph_query(self, node: str) -> Dict[str, Any]:
        return self._request("POST", "/query/subgraph", {"node": node})

    def grounding_assemble(self, query: str, principal: str, purpose: str, budget: int = 0) -> Dict[str, Any]:
        return self._request("POST", "/grounding/assemble", {"query": query, "principal": principal, "purpose": purpose, "budget": budget})

    def grounding_explain(self, query: str, principal: str, purpose: str) -> Dict[str, Any]:
        return self._request("POST", "/grounding/explain", {"query": query, "principal": principal, "purpose": purpose})

    def grounding_redact(self, query: str, principal: str, purpose: str, hide_tags=None) -> Dict[str, Any]:
        return self._request("POST", "/grounding/redact", {"query": query, "principal": principal, "purpose": purpose, "hideTags": hide_tags or []})

    # ---- consistency / governance ----

    def consistency_report(self, ontology_namespace: str = "") -> Dict[str, Any]:
        path = "/consistency/report"
        if ontology_namespace:
            path += "?ontology=" + urllib.parse.quote(ontology_namespace)
        return self._request("GET", path)

    def consistency_conflicts(self) -> Dict[str, Any]:
        return self._request("GET", "/consistency/conflicts")

    def consistency_resolve(self) -> Dict[str, Any]:
        return self._request("POST", "/consistency/resolve")

    def consistency_audit(self) -> Dict[str, Any]:
        return self._request("GET", "/consistency/audit")

    def policy_list(self) -> Dict[str, Any]:
        return self._request("GET", "/policy/policies")

    def policy_evaluate(self, req: Dict[str, Any]) -> Dict[str, Any]:
        return self._request("POST", "/policy/evaluate", req)

    def approval_create(self, action: str, resource: str, requester: str = "", reason: str = "") -> Dict[str, Any]:
        return self._request("POST", "/approval/requests", {"action": action, "resource": resource, "requester": requester, "reason": reason})

    def approval_list(self) -> Dict[str, Any]:
        return self._request("GET", "/approval/requests")

    def approval_get(self, req_id: str) -> Dict[str, Any]:
        return self._request("GET", "/approval/requests/" + urllib.parse.quote(req_id))

    def approval_approve(self, req_id: str, approver: str, role: str, comment: str = "") -> Dict[str, Any]:
        return self._request("POST", "/approval/requests/%s/approve" % urllib.parse.quote(req_id), {"approver": approver, "role": role, "comment": comment})

    def approval_reject(self, req_id: str, approver: str, role: str, comment: str = "") -> Dict[str, Any]:
        return self._request("POST", "/approval/requests/%s/reject" % urllib.parse.quote(req_id), {"approver": approver, "role": role, "comment": comment})

    def audit_list(self) -> Dict[str, Any]:
        return self._request("GET", "/audit")

    def audit_attest(self, entry_id: str, ref: str, content: str, by: str) -> Dict[str, Any]:
        return self._request("POST", "/audit/attest", {"entryId": entry_id, "ref": ref, "content": content, "by": by})

    # ---- pipeline / federation / bridge / events ----

    def pipeline_run(self, req: Dict[str, Any]) -> Dict[str, Any]:
        return self._request("POST", "/pipeline/run", req)

    def federation_query(self, query: str, principal: str = "", purpose: str = "") -> Dict[str, Any]:
        return self._request("POST", "/federation/query", {"query": query, "principal": principal, "purpose": purpose})

    def bridge_contract_export(self, namespace: str) -> Dict[str, Any]:
        return self._request("GET", "/bridge/%s/contract" % urllib.parse.quote(namespace))

    def bridge_query(self, namespace: str, query: str, principal: str = "", purpose: str = "") -> Dict[str, Any]:
        return self._request("POST", "/bridge/query", {"namespace": namespace, "query": query, "principal": principal, "purpose": purpose})

    def event_contract_register(self, contract: Dict[str, Any]) -> Dict[str, Any]:
        return self._request("POST", "/events/contracts", contract)

    def event_contract_list(self) -> Dict[str, Any]:
        return self._request("GET", "/events/contracts")

    def event_ingest(self, contract_id: str, event: Dict[str, Any]) -> Dict[str, Any]:
        return self._request("POST", "/events/ingest", {"contractId": contract_id, "event": event})

    def metrics(self) -> Dict[str, Any]:
        return self._request("GET", "/metrics")
