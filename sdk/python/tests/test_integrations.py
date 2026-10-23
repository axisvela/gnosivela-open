"""Tests for the dependency-free agent framework adapters."""

import json
import unittest

from gnosivela.integrations import ToolSpec, langchain_tools, llama_index_tools, tool_specs


class FakeClient:
    def __init__(self):
        self.calls = []

    def grounding_assemble(self, query, principal, purpose, budget):
        self.calls.append(("ground", query, principal, purpose, budget))
        return {"resolvedIntent": query, "knowledge": [{"assertionId": "a1"}]}

    def grounding_explain(self, query, principal, purpose):
        self.calls.append(("cite", query, principal, purpose))
        return {"intent": query, "citations": [{"assertionId": "a1"}], "citationCompleteness": 1.0}

    def semantic_query(self, query, principal, purpose):
        self.calls.append(("search", query, principal, purpose))
        return {"intent": {"query": query}, "view": {"assertions": []}}

    def hybrid_query(self, query, limit=10, weights=None):
        self.calls.append(("hybrid", query, limit))
        return {"query": query, "fusion": "rrf", "hits": []}

    def entity_resolve(self, hint):
        self.calls.append(("resolve", hint))
        return {"hint": hint, "canonical": {"canonicalId": "C-1"}}


class TestIntegrations(unittest.TestCase):
    def test_tool_set(self):
        client = FakeClient()
        tools = tool_specs(client, "alice", "onboarding")
        self.assertEqual(len(tools), 5)
        names = [t.name for t in tools]
        self.assertIn("gnosivela_ground", names)
        self.assertIn("gnosivela_hybrid_search", names)
        for tool in tools:
            self.assertIsInstance(tool, ToolSpec)
            self.assertEqual(tool.parameters["type"], "object")

    def test_ground_invocation(self):
        client = FakeClient()
        ground = next(t for t in tool_specs(client, "alice", "onboarding") if t.name == "gnosivela_ground")
        out = json.loads(ground.invoke({"query": "ACME risk", "budget": 3}))
        self.assertEqual(out["resolvedIntent"], "ACME risk")
        self.assertEqual(client.calls[0], ("ground", "ACME risk", "alice", "onboarding", 3))

    def test_other_tools(self):
        client = FakeClient()
        by_name = {t.name: t for t in tool_specs(client, "bob", "review")}
        self.assertEqual(json.loads(by_name["gnosivela_cite"].invoke({"query": "q"}))["citationCompleteness"], 1.0)
        self.assertEqual(json.loads(by_name["gnosivela_search"].invoke({"query": "q"}))["intent"]["query"], "q")
        hybrid = json.loads(by_name["gnosivela_hybrid_search"].invoke({"query": "q", "limit": 3}))
        self.assertEqual(hybrid["fusion"], "rrf")
        resolved = json.loads(by_name["gnosivela_resolve_entity"].invoke({"hint": "ACME"}))
        self.assertEqual(resolved["canonical"]["canonicalId"], "C-1")

    def test_langchain_fallback(self):
        tools = langchain_tools(FakeClient(), "alice", "onboarding")
        self.assertEqual(len(tools), 5)
        # Without langchain installed these remain neutral specs.
        try:
            import langchain_core  # noqa: F401
        except ImportError:
            self.assertTrue(all(isinstance(t, ToolSpec) for t in tools))

    def test_llama_index_fallback(self):
        tools = llama_index_tools(FakeClient(), "alice", "onboarding")
        self.assertEqual(len(tools), 5)
        try:
            import llama_index  # noqa: F401
        except ImportError:
            self.assertTrue(all(isinstance(t, ToolSpec) for t in tools))


if __name__ == "__main__":
    unittest.main()
