// Package gnosivela is the official Go client SDK for the GNOSIVELA
// Semantic Control Plane API. It is generated from the semantic contract
// (contracts/openapi.yaml) and is Apache-2.0 licensed — it never links the AGPL
// core; it talks to it over HTTP.
package gnosivela

import "time"

// OntologyStatus is the lifecycle state of an ontology.
type OntologyStatus string

const (
	OntologyDraft      OntologyStatus = "draft"
	OntologyValidated  OntologyStatus = "validated"
	OntologyPublished  OntologyStatus = "published"
	OntologyDeprecated OntologyStatus = "deprecated"
	OntologyRevoked    OntologyStatus = "revoked"
)

// AssertionStatus is the lifecycle state of a knowledge assertion.
type AssertionStatus string

const (
	AssertionProposed   AssertionStatus = "proposed"
	AssertionValidated  AssertionStatus = "validated"
	AssertionContested  AssertionStatus = "contested"
	AssertionSuperseded AssertionStatus = "superseded"
	AssertionRevoked    AssertionStatus = "revoked"
)

// Ontology is the machine-executable semantic contract.
type Ontology struct {
	ID           string          `json:"id"`
	Namespace    string          `json:"namespace"`
	Version      string          `json:"version"`
	Dependencies []string        `json:"dependencies,omitempty"`
	Concepts     []*Concept      `json:"concepts,omitempty"`
	Relations    []*RelationSpec `json:"relations,omitempty"`
	Constraints  []*Constraint   `json:"constraints,omitempty"`
	Status       OntologyStatus  `json:"status"`
	Owner        string          `json:"owner,omitempty"`
	Signature    string          `json:"signature,omitempty"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
}

// Concept defines a business object type.
type Concept struct {
	ID           string      `json:"id"`
	Type         string      `json:"type,omitempty"`
	Label        string      `json:"label"`
	Definition   string      `json:"definition,omitempty"`
	Scope        string      `json:"scope,omitempty"`
	Owner        string      `json:"owner,omitempty"`
	Properties   []*Property `json:"properties,omitempty"`
	DeprecatedBy string      `json:"deprecatedBy,omitempty"`
}

// Property defines an attribute of a concept.
type Property struct {
	Name        string   `json:"name"`
	DataType    string   `json:"dataType"`
	Unit        string   `json:"unit,omitempty"`
	Required    bool     `json:"required,omitempty"`
	Cardinality string   `json:"cardinality,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

// RelationSpec defines a typed relationship between two concepts.
type RelationSpec struct {
	ID          string   `json:"id"`
	Domain      string   `json:"domain"`
	Range       string   `json:"range"`
	Roles       []string `json:"roles,omitempty"`
	Cardinality string   `json:"cardinality,omitempty"`
	Inverse     string   `json:"inverse,omitempty"`
	Properties  []string `json:"properties,omitempty"`
}

// Constraint declares a validation or inference rule.
type Constraint struct {
	ID         string `json:"id"`
	Target     string `json:"target"`
	Expression string `json:"expression"`
	Severity   string `json:"severity"`
	Message    string `json:"message,omitempty"`
	Version    string `json:"version,omitempty"`
}

// EntityRef is a canonical reference to an entity.
type EntityRef struct {
	Namespace   string   `json:"namespace"`
	CanonicalID string   `json:"canonicalId"`
	Type        string   `json:"type,omitempty"`
	Aliases     []string `json:"aliases,omitempty"`
	Authority   string   `json:"authority,omitempty"`
}

// Value is a typed object value or target entity of an assertion.
type Value struct {
	Type      string      `json:"type,omitempty"`
	EntityRef *EntityRef  `json:"entityRef,omitempty"`
	String    *string     `json:"string,omitempty"`
	Number    *float64    `json:"number,omitempty"`
	Boolean   *bool       `json:"boolean,omitempty"`
}

// AssertionContext describes the business scope of an assertion.
type AssertionContext struct {
	Domain   string   `json:"domain,omitempty"`
	Region   string   `json:"region,omitempty"`
	Contract string   `json:"contract,omitempty"`
	Scenario string   `json:"scenario,omitempty"`
	Purpose  string   `json:"purpose,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

// KnowledgeAssertion is the ontology-governed, evidenced knowledge claim.
type KnowledgeAssertion struct {
	AssertionID     string           `json:"assertionId"`
	Subject         EntityRef        `json:"subject"`
	Predicate       string           `json:"predicate"`
	Object          Value            `json:"object"`
	Context         AssertionContext `json:"context"`
	ValidFrom       time.Time        `json:"validFrom,omitempty"`
	ValidTo         time.Time        `json:"validTo,omitempty"`
	RecordedAt      time.Time        `json:"recordedAt"`
	Source          string           `json:"source"`
	Authority       string           `json:"authority,omitempty"`
	Confidence      float64          `json:"confidence,omitempty"`
	Status          AssertionStatus  `json:"status"`
	PolicyTags      []string         `json:"policyTags,omitempty"`
	OntologyVersion string           `json:"ontologyVersion,omitempty"`
}

// EvidenceRef points at the source artifact supporting an assertion.
type EvidenceRef struct {
	Source         string    `json:"source"`
	ArtifactDigest string    `json:"artifactDigest,omitempty"`
	Locator        string    `json:"locator,omitempty"`
	CapturedAt     time.Time `json:"capturedAt,omitempty"`
	Principal      string    `json:"principal,omitempty"`
}

// Issue is a validation result item.
type Issue struct {
	Severity string `json:"severity"`
	Path     string `json:"path"`
	Message  string `json:"message"`
}

// Change describes a single semantic diff entry.
type Change struct {
	Kind     string `json:"kind"`
	Type     string `json:"type"`
	ID       string `json:"id"`
	Before   string `json:"before,omitempty"`
	After    string `json:"after,omitempty"`
	Breaking bool   `json:"breaking,omitempty"`
}

// OntologyCreateResult is returned by OntologyCreate.
type OntologyCreateResult struct {
	Ontology *Ontology `json:"ontology"`
	Issues   []Issue   `json:"issues"`
}

// RelationKind is the type of an identity relationship.
type RelationKind string

const (
	RelationExactSameAs       RelationKind = "exactSameAs"
	RelationLegallyRepresents RelationKind = "legallyRepresents"
	RelationOperationalAlias  RelationKind = "operationalAlias"
	RelationMemberOf          RelationKind = "memberOf"
	RelationPartOf            RelationKind = "partOf"
	RelationDerivedFrom       RelationKind = "derivedFrom"
	RelationSimilarTo         RelationKind = "similarTo"
)

// Relation links two entity references with a typed identity relation.
type Relation struct {
	ID        string       `json:"id"`
	From      EntityRef    `json:"from"`
	To        EntityRef    `json:"to"`
	Kind      RelationKind `json:"kind"`
	Evidence  []EvidenceRef `json:"evidence,omitempty"`
	Authority string       `json:"authority,omitempty"`
}

// MatchEvidence records why two references were matched.
type MatchEvidence struct {
	Left        EntityRef  `json:"left"`
	Right       EntityRef  `json:"right"`
	Score       float64    `json:"score"`
	MatchedOn   []string   `json:"matchedOn"`
	GeneratedBy string     `json:"generatedBy"`
	Status      string     `json:"status"`
	CapturedAt  time.Time  `json:"capturedAt"`
}

// ResolvedMatch is a candidate match with its evidence.
type ResolvedMatch struct {
	Ref       EntityRef `json:"ref"`
	Score     float64   `json:"score"`
	MatchedOn []string  `json:"matchedOn"`
	Status    string    `json:"status"`
}

// Resolution is the purpose-scoped outcome of an entity resolution request.
type Resolution struct {
	Hint      string       `json:"hint"`
	Canonical *EntityRef   `json:"canonical,omitempty"`
	Matches   []ResolvedMatch `json:"matches"`
	Related   []*Relation  `json:"related,omitempty"`
	Gaps      []string     `json:"gaps,omitempty"`
	AsOf      time.Time    `json:"asOf"`
}

// Intent is the resolved business intent of a semantic query.
type Intent struct {
	Query     string      `json:"query"`
	Concepts  []string    `json:"concepts,omitempty"`
	Entities  []EntityRef `json:"entities,omitempty"`
	Relations []string    `json:"relations,omitempty"`
	Region    string      `json:"region,omitempty"`
	Purpose   string      `json:"purpose,omitempty"`
	Tokens    []string    `json:"tokens,omitempty"`
}

// Conflict captures competing assertions that are not silently merged.
type Conflict struct {
	Subject    EntityRef             `json:"subject"`
	Predicate  string                `json:"predicate"`
	Assertions []*KnowledgeAssertion `json:"assertions"`
}

// KnowledgeView is the purpose-scoped resolution of assertions.
type KnowledgeView struct {
	Purpose    string                `json:"purpose"`
	Principal  string                `json:"principal"`
	AsOf       time.Time             `json:"asOf"`
	Assertions []*KnowledgeAssertion `json:"assertions,omitempty"`
	Conflicts  []*Conflict           `json:"conflicts,omitempty"`
	Gaps       []string              `json:"gaps,omitempty"`
}

// EvidenceChain ties an assertion to its sources and derivations.
type EvidenceChain struct {
	AssertionID string   `json:"assertionId"`
	Sources     []string `json:"sources,omitempty"`
	Derivations []string `json:"derivations,omitempty"`
}

// PathStep is one hop in a multi-hop path.
type PathStep struct {
	From     string `json:"from"`
	Relation string `json:"relation"`
	To       string `json:"to"`
}

// Path is a multi-hop relationship path with explanations.
type Path struct {
	Nodes        []string   `json:"nodes"`
	Steps        []PathStep `json:"steps"`
	Explanations []string   `json:"explanations"`
}

// MetricDefinition defines a governed metric's semantics.
type MetricDefinition struct {
	ID         string   `json:"id"`
	Formula    string   `json:"formula"`
	Dimensions []string `json:"dimensions,omitempty"`
	Grain      string   `json:"grain,omitempty"`
	Unit       string   `json:"unit,omitempty"`
	Source     string   `json:"source,omitempty"`
	Owner      string   `json:"owner,omitempty"`
	Version    string   `json:"version,omitempty"`
}

// SemanticResult is the outcome of a semantic query.
type SemanticResult struct {
	Intent         Intent             `json:"intent"`
	View           *KnowledgeView     `json:"view"`
	Relations      []*Relation        `json:"relations,omitempty"`
	Paths          []*Path            `json:"paths,omitempty"`
	EvidenceChains []EvidenceChain    `json:"evidenceChains,omitempty"`
	Metrics        []*MetricDefinition `json:"metrics,omitempty"`
}

// SemanticContract carries the ontology/mapping versions behind a bundle.
type SemanticContract struct {
	OntologyVersion string   `json:"ontologyVersion"`
	MappingVersion  string   `json:"mappingVersion,omitempty"`
	Units           []string `json:"units,omitempty"`
}

// GroundingBundle is the assembled, evidence-backed context for one task.
type GroundingBundle struct {
	ResolvedIntent   string                `json:"resolvedIntent"`
	Knowledge        []*KnowledgeAssertion `json:"knowledge,omitempty"`
	Evidence         []EvidenceRef         `json:"evidence,omitempty"`
	SemanticContract *SemanticContract     `json:"semanticContract,omitempty"`
	PolicyDecision   string                `json:"policyDecision,omitempty"`
	Conflict         []*Conflict           `json:"conflict,omitempty"`
	Gap              []string              `json:"gap,omitempty"`
	ContextBudget    int                   `json:"contextBudget,omitempty"`
}

// GroundingCitation is one assertion's evidence citation in an explanation.
type GroundingCitation struct {
	AssertionID     string        `json:"assertionId"`
	Subject         string        `json:"subject"`
	Predicate       string        `json:"predicate"`
	Status          string        `json:"status"`
	Confidence      float64       `json:"confidence,omitempty"`
	OntologyVersion string        `json:"ontologyVersion,omitempty"`
	Sources         []string      `json:"sources,omitempty"`
	Derivations     []string      `json:"derivations,omitempty"`
	Evidence        []EvidenceRef `json:"evidence,omitempty"`
}

// GroundingExplanation is the citation-level breakdown of an assembly.
type GroundingExplanation struct {
	Intent               string              `json:"intent"`
	Citations            []GroundingCitation `json:"citations"`
	Conflicts            []string            `json:"conflicts,omitempty"`
	Gaps                 []string            `json:"gaps,omitempty"`
	Policy               string              `json:"policy"`
	CitationCompleteness float64             `json:"citationCompleteness"`
}
