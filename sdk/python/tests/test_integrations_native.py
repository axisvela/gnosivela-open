import importlib.util
import unittest

from gnosivela.client import Client
from gnosivela.integrations import langchain_tools, llama_index_tools, tool_specs

HAS_LANGCHAIN = importlib.util.find_spec("langchain_core") is not None
try:
    HAS_LLAMA_INDEX = importlib.util.find_spec("llama_index.core") is not None
except ModuleNotFoundError:
    HAS_LLAMA_INDEX = False

EXPECTED = {
    "gnosivela_ground",
    "gnosivela_cite",
    "gnosivela_search",
    "gnosivela_hybrid_search",
    "gnosivela_resolve_entity",
}


class NativeIntegrationTests(unittest.TestCase):
    def setUp(self):
        self.client = Client("http://127.0.0.1:1")

    def test_neutral_specs_always_available(self):
        specs = tool_specs(self.client)
        self.assertEqual({spec.name for spec in specs}, EXPECTED)
        self.assertEqual(len(specs), 5)

    @unittest.skipUnless(HAS_LANGCHAIN, "langchain_core is not installed")
    def test_langchain_structured_tools(self):
        from langchain_core.tools import StructuredTool

        tools = langchain_tools(self.client, principal="alice", purpose="review")
        self.assertEqual(len(tools), 5)
        self.assertTrue(all(isinstance(tool, StructuredTool) for tool in tools))
        self.assertEqual({tool.name for tool in tools}, EXPECTED)
        self.assertTrue(all(tool.description for tool in tools))

    @unittest.skipUnless(HAS_LLAMA_INDEX, "llama_index.core is not installed")
    def test_llama_index_function_tools(self):
        from llama_index.core.tools import FunctionTool

        tools = llama_index_tools(self.client, principal="alice", purpose="review")
        self.assertEqual(len(tools), 5)
        self.assertTrue(all(isinstance(tool, FunctionTool) for tool in tools))
        self.assertEqual({tool.metadata.name for tool in tools}, EXPECTED)
        self.assertTrue(all(tool.metadata.description for tool in tools))


if __name__ == "__main__":
    unittest.main()
