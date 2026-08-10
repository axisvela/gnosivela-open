"""Core DTOs for the GNOSIVELA control plane. Light dataclasses for the most
used objects; endpoints returning richer structures use plain dicts decoded
from JSON so the client stays small and dependency-free."""

from dataclasses import dataclass, field
from typing import Any, Dict, List, Optional


@dataclass
class EntityRef:
    namespace: str
    canonical_id: str
    type: str = ""
    aliases: List[str] = field(default_factory=list)
    authority: str = ""

    def to_dict(self) -> Dict[str, Any]:
        d: Dict[str, Any] = {"namespace": self.namespace, "canonicalId": self.canonical_id}
        if self.type:
            d["type"] = self.type
        if self.aliases:
            d["aliases"] = self.aliases
        if self.authority:
            d["authority"] = self.authority
        return d


@dataclass
class Value:
    type: str = "string"  # string | number | boolean | date | entity
    string: Optional[str] = None
    number: Optional[float] = None
    boolean: Optional[bool] = None
    entity: Optional[EntityRef] = None

    def to_dict(self) -> Dict[str, Any]:
        d: Dict[str, Any] = {"type": self.type}
        if self.string is not None:
            d["string"] = self.string
        if self.number is not None:
            d["number"] = self.number
        if self.boolean is not None:
            d["boolean"] = self.boolean
        if self.entity is not None:
            d["entityRef"] = self.entity.to_dict()
        return d


@dataclass
class AssertionContext:
    region: str = ""
    domain: str = ""
    contract: str = ""
    scenario: str = ""
    purpose: str = ""
    tags: List[str] = field(default_factory=list)

    def to_dict(self) -> Dict[str, Any]:
        d: Dict[str, Any] = {}
        if self.region:
            d["region"] = self.region
        if self.domain:
            d["domain"] = self.domain
        if self.contract:
            d["contract"] = self.contract
        if self.scenario:
            d["scenario"] = self.scenario
        if self.purpose:
            d["purpose"] = self.purpose
        if self.tags:
            d["tags"] = self.tags
        return d


@dataclass
class EvidenceRef:
    source: str = ""
    locator: str = ""
    artifact_digest: str = ""
    principal: str = ""

    def to_dict(self) -> Dict[str, Any]:
        d: Dict[str, Any] = {"source": self.source}
        if self.locator:
            d["locator"] = self.locator
        if self.artifact_digest:
            d["artifactDigest"] = self.artifact_digest
        if self.principal:
            d["principal"] = self.principal
        return d


@dataclass
class KnowledgeAssertion:
    assertion_id: str
    subject: EntityRef
    predicate: str
    object: Value
    source: str = ""
    status: str = "proposed"
    confidence: float = 0.0
    context: Optional[AssertionContext] = None
    evidence: List[EvidenceRef] = field(default_factory=list)
    valid_from: Optional[str] = None
    valid_to: Optional[str] = None
    recorded_at: str = ""
    authority: str = ""
    ontology_version: str = ""
    policy_tags: List[str] = field(default_factory=list)

    def to_dict(self) -> Dict[str, Any]:
        d: Dict[str, Any] = {
            "assertionId": self.assertion_id,
            "subject": self.subject.to_dict(),
            "predicate": self.predicate,
            "object": self.object.to_dict(),
            "status": self.status,
        }
        if self.source:
            d["source"] = self.source
        if self.confidence:
            d["confidence"] = self.confidence
        if self.context:
            d["context"] = self.context.to_dict()
        if self.evidence:
            d["evidence"] = [e.to_dict() for e in self.evidence]
        if self.valid_from:
            d["validFrom"] = self.valid_from
        if self.valid_to:
            d["validTo"] = self.valid_to
        if self.recorded_at:
            d["recordedAt"] = self.recorded_at
        if self.authority:
            d["authority"] = self.authority
        if self.ontology_version:
            d["ontologyVersion"] = self.ontology_version
        if self.policy_tags:
            d["policyTags"] = self.policy_tags
        return d
