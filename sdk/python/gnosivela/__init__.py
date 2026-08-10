"""GNOSIVELA — Semantic Control Plane Python client (Apache-2.0).

The client talks to the control plane over HTTP only; it never links the AGPL
core. It mirrors the Go SDK: ontologies, assertions, entities, semantic query,
grounding, consistency, policy, approval, audit, pipeline, federation, bridge
and real-time events.
"""

from .client import Client, GnosivelaError

__all__ = ["Client", "GnosivelaError"]
__version__ = "1.0.0-beta.1"
