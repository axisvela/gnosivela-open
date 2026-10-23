"""Agent framework adapters (LangChain / LlamaIndex) for GNOSIVELA.

The adapters are dependency-free. If ``langchain_core`` or ``llama_index`` is
installed, the helper functions return native tools; otherwise they return
neutral :class:`ToolSpec` objects that any framework can wrap.

Every tool calls the governed control-plane API (grounding / semantic query /
entity resolution). The adapters never read the store directly and never copy
semantics into the agent framework.
"""

from __future__ import annotations

import json
from dataclasses import dataclass
from typing import Any, Callable, Dict, List, Optional

from .client import Client

JsonDict = Dict[str, Any]


@dataclass
class ToolSpec:
    """A framework-neutral tool descriptor."""

    name: str
    description: str
    parameters: JsonDict
    run: Callable[..., Any]

    def invoke(self, arguments: Optional[JsonDict] = None) -> Any:
        return self.run(**(arguments or {}))


def _string(description: str) -> JsonDict:
    return {"type": "string", "description": description}


def _integer(description: str) -> JsonDict:
    return {"type": "integer", "description": description}


def _object(properties: JsonDict, required: List[str]) -> JsonDict:
    return {
        "type": "object",
        "properties": properties,
        "required": required,
        "additionalProperties": False,
    }


def _json(value: Any) -> str:
    return json.dumps(value, ensure_ascii=False, sort_keys=True)


def tool_specs(
    client: Client,
    principal: str = "",
    purpose: str = "",
) -> List[ToolSpec]:
    """Return the GNOSIVELA tool set as neutral specs."""

    def ground(
        query: str,
        principal: str = principal,
        purpose: str = purpose,
        budget: int = 0,
    ) -> str:
        return _json(client.grounding_assemble(query, principal, purpose, budget))

    def cite(query: str, principal: str = principal, purpose: str = purpose) -> str:
        return _json(client.grounding_explain(query, principal, purpose))

    def search(query: str, principal: str = principal, purpose: str = purpose) -> str:
        return _json(client.semantic_query(query, principal, purpose))

    def hybrid_search(query: str, limit: int = 10) -> str:
        return _json(client.hybrid_query(query, limit=limit))

    def resolve_entity(hint: str) -> str:
        return _json(client.entity_resolve(hint))

    base = {
        "query": _string("Natural-language or structured question."),
        "principal": _string("Principal on whose behalf the tool runs."),
        "purpose": _string("Purpose scope for the knowledge view."),
    }
    return [
        ToolSpec(
            name="gnosivela_ground",
            description=(
                "Assemble a governed Grounding Bundle (knowledge, evidence, "
                "contract, conflicts, gaps) for a question."
            ),
            parameters=_object({**base, "budget": _integer("Max knowledge items (0 = unlimited).")}, ["query"]),
            run=ground,
        ),
        ToolSpec(
            name="gnosivela_cite",
            description="Explain the citations behind a grounded answer.",
            parameters=_object(dict(base), ["query"]),
            run=cite,
        ),
        ToolSpec(
            name="gnosivela_search",
            description="Run a semantic query and return a purpose-scoped KnowledgeView.",
            parameters=_object(dict(base), ["query"]),
            run=search,
        ),
        ToolSpec(
            name="gnosivela_hybrid_search",
            description="Hybrid BM25/vector/graph retrieval fused with RRF.",
            parameters=_object(
                {
                    "query": _string("Search query."),
                    "limit": _integer("Maximum number of hits (default 10)."),
                },
                ["query"],
            ),
            run=hybrid_search,
        ),
        ToolSpec(
            name="gnosivela_resolve_entity",
            description="Resolve an entity hint with precision-first match evidence.",
            parameters=_object({"hint": _string("Entity name, alias or canonical ID.")}, ["hint"]),
            run=resolve_entity,
        ),
    ]


def langchain_tools(
    client: Client,
    principal: str = "",
    purpose: str = "",
) -> List[Any]:
    """Return LangChain ``StructuredTool`` objects when LangChain is installed.

    Falls back to :func:`tool_specs` when ``langchain_core`` is unavailable.
    """
    specs = tool_specs(client, principal, purpose)
    try:
        from langchain_core.tools import StructuredTool
    except ImportError:
        return specs
    return [
        StructuredTool.from_function(
            func=spec.run,
            name=spec.name,
            description=spec.description,
        )
        for spec in specs
    ]


def llama_index_tools(
    client: Client,
    principal: str = "",
    purpose: str = "",
) -> List[Any]:
    """Return LlamaIndex ``FunctionTool`` objects when LlamaIndex is installed.

    Falls back to :func:`tool_specs` when ``llama_index`` is unavailable.
    """
    specs = tool_specs(client, principal, purpose)
    try:
        from llama_index.core.tools import FunctionTool
    except ImportError:
        return specs
    return [
        FunctionTool.from_defaults(fn=spec.run, name=spec.name, description=spec.description)
        for spec in specs
    ]
