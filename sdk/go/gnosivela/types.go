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
	Type      string     `json:"type,omitempty"`
	EntityRef *EntityRef `json:"entityRef,omitempty"`
	String    *string    `json:"string,omitempty"`
	Number    *float64   `json:"number,omitempty"`
	Boolean   *bool      `json:"boolean,omitempty"`
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

// AffectedAssertion references a knowledge assertion produced under an old
// ontology version that now references removed or narrowed elements.
type AffectedAssertion struct {
	AssertionID     string    `json:"assertionId"`
	Subject         EntityRef `json:"subject"`
	Predicate       string    `json:"predicate"`
	OntologyVersion string    `json:"ontologyVersion"`
	Reason          string    `json:"reason"`
}

// ImpactSummary reduces an ImpactReport to gate-relevant counters.
type ImpactSummary struct {
	Changes            int `json:"changes"`
	Breaking           int `json:"breaking"`
	Additive           int `json:"additive"`
	Deprecations       int `json:"deprecations"`
	AffectedAssertions int `json:"affectedAssertions"`
	AffectedEntities   int `json:"affectedEntities"`
}

// ImpactReport is the outcome of analysing a candidate ontology version
// against the currently published version of the same namespace.
type ImpactReport struct {
	Namespace          string              `json:"namespace"`
	FromVersion        string              `json:"fromVersion"`
	ToVersion          string              `json:"toVersion"`
	Breaking           []Change            `json:"breaking,omitempty"`
	Additive           []Change            `json:"additive,omitempty"`
	Deprecations       []Change            `json:"deprecations,omitempty"`
	AffectedAssertions []AffectedAssertion `json:"affectedAssertions,omitempty"`
	AffectedEntities   []EntityRef         `json:"affectedEntities,omitempty"`
	Blocking           bool                `json:"blocking"`
}

// PublishResult is returned by publish, carrying the impact summary.
type PublishResult struct {
	Status string         `json:"status"`
	Impact *ImpactSummary `json:"impact"`
}

// RuleResult is the outcome of evaluating one ontology constraint against one
// assertion.
type RuleResult struct {
	ConstraintID string `json:"constraintId"`
	Expression   string `json:"expression"`
	AssertionID  string `json:"assertionId"`
	Passed       bool   `json:"passed"`
	Severity     string `json:"severity"`
	Message      string `json:"message,omitempty"`
}

// CheckResult is a single consistency finding with severity and a reference.
type CheckResult struct {
	Type     string `json:"type"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Ref      string `json:"ref,omitempty"`
}

// ConsistencyReport aggregates rule validation and the seven consistency
// checks (lexical/identity/structural/temporal/contextual/inferential/authority).
type ConsistencyReport struct {
	Rules    []RuleResult  `json:"rules,omitempty"`
	Checks   []CheckResult `json:"checks,omitempty"`
	Failures int           `json:"failures"`
	Warnings int           `json:"warnings"`
}

// ConsistencyResolution records how a group of competing assertions was
// settled. Conflicts are never resolved silently: every decision carries a
// basis, losers are marked contested, and ties are escalated without mutation.
type ConsistencyResolution struct {
	At        time.Time `json:"at"`
	Subject   EntityRef `json:"subject"`
	Predicate string    `json:"predicate"`
	Winner    string    `json:"winner,omitempty"`
	Loser     string    `json:"loser,omitempty"`
	Basis     string    `json:"basis"`
	Escalated bool      `json:"escalated"`
	Ties      []string  `json:"ties,omitempty"`
}

// PolicyEffect is allow or deny.
type PolicyEffect string

const (
	PolicyEffectAllow PolicyEffect = "allow"
	PolicyEffectDeny  PolicyEffect = "deny"
)

// Condition constrains a policy to requests whose field matches.
type Condition struct {
	Field  string   `json:"field"`
	Op     string   `json:"op"` // eq | neq | in
	Values []string `json:"values,omitempty"`
}

// Policy is a declarative attribute-based access rule.
type Policy struct {
	ID         string       `json:"id"`
	Effect     PolicyEffect `json:"effect"`
	Actions    []string     `json:"actions,omitempty"`
	Resources  []string     `json:"resources,omitempty"`
	Conditions []Condition  `json:"conditions,omitempty"`
	Callout    string       `json:"callout,omitempty"`
}

// PolicyRequest is the attribute set evaluated by the policy engine.
type PolicyRequest struct {
	Principal  string            `json:"principal,omitempty"`
	Action     string            `json:"action"`
	Resource   string            `json:"resource,omitempty"`
	Tenant     string            `json:"tenant,omitempty"`
	Region     string            `json:"region,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// PolicyDecision is the evaluation outcome with matched policies and a reason.
type PolicyDecision struct {
	Allowed  bool     `json:"allowed"`
	Reason   string   `json:"reason"`
	Policies []string `json:"policies,omitempty"`
}

// ApprovalRole is a role in an approval chain.
type ApprovalRole string

const (
	ApprovalRoleDomainOwner ApprovalRole = "domain-owner"
	ApprovalRoleSteward     ApprovalRole = "steward"
	ApprovalRoleSecurity    ApprovalRole = "security"
)

// ApprovalStep is one role's decision in an approval chain.
type ApprovalStep struct {
	Role     ApprovalRole `json:"role"`
	Status   string       `json:"status"` // pending | approved | rejected
	Approver string       `json:"approver,omitempty"`
	Comment  string       `json:"comment,omitempty"`
	At       time.Time    `json:"at,omitempty"`
}

// ApprovalRequest is a pending or settled approval over an action/resource.
type ApprovalRequest struct {
	ID        string         `json:"id"`
	Action    string         `json:"action"`
	Resource  string         `json:"resource"`
	Requester string         `json:"requester"`
	Reason    string         `json:"reason,omitempty"`
	Steps     []ApprovalStep `json:"steps"`
	Status    string         `json:"status"` // pending | approved | rejected
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

// ApprovalCreate is the payload to open an approval request.
type ApprovalCreate struct {
	Action    string `json:"action"`
	Resource  string `json:"resource"`
	Requester string `json:"requester,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

// AuditEvidence is a digest-attested source artifact.
type AuditEvidence struct {
	Ref        string    `json:"ref"`
	Digest     string    `json:"digest,omitempty"`
	AttestedAt time.Time `json:"attestedAt,omitempty"`
	By         string    `json:"by,omitempty"`
}

// AuditEntry is one audited governance decision.
type AuditEntry struct {
	ID       string          `json:"id"`
	At       time.Time       `json:"at"`
	Actor    string          `json:"actor"`
	Action   string          `json:"action"`
	Resource string          `json:"resource"`
	Outcome  string          `json:"outcome"`
	Detail   string          `json:"detail,omitempty"`
	Evidence []AuditEvidence `json:"evidence,omitempty"`
}

// PipelineRequest is a single mapping CI/CD run over one ontology version.
type PipelineRequest struct {
	Namespace       string `json:"namespace"`
	Version         string `json:"version"`
	Principal       string `json:"principal,omitempty"`
	ApprovalID      string `json:"approvalId,omitempty"`
	RequireApproval bool   `json:"requireApproval,omitempty"`
	MaxFailures     int    `json:"maxFailures,omitempty"`
	AutoRollback    bool   `json:"autoRollback,omitempty"`
}

// PipelineStage is the outcome of one stage of a pipeline run.
type PipelineStage struct {
	Stage  string `json:"stage"`
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

// PipelineReport is the full record of a pipeline run.
type PipelineReport struct {
	Namespace   string          `json:"namespace"`
	Version     string          `json:"version"`
	FinalStatus string          `json:"finalStatus"` // deployed | rolled-back | blocked | failed | degraded
	Stages      []PipelineStage `json:"stages"`
}

// FederatedView is the merged, conflict-preserving result of a federated
// semantic query across autonomous domains.
type FederatedView struct {
	Query          string                `json:"query"`
	Purpose        string                `json:"purpose"`
	AsOf           time.Time             `json:"asOf"`
	Assertions     []*KnowledgeAssertion `json:"assertions,omitempty"`
	Conflicts      []*Conflict           `json:"conflicts,omitempty"`
	Gaps           []string              `json:"gaps,omitempty"`
	DomainHits     map[string]int        `json:"domainHits"`
	DomainOf       map[string]string     `json:"domainOf,omitempty"`
	LatencyMillis  int64                 `json:"latencyMillis"`
	DomainsQueried int                   `json:"domainsQueried"`
}

// ContractProperty is a governed attribute exposed by a semantic contract.
type ContractProperty struct {
	Name     string   `json:"name"`
	DataType string   `json:"dataType"`
	Unit     string   `json:"unit,omitempty"`
	Required bool     `json:"required,omitempty"`
	Enum     []string `json:"enum,omitempty"`
}

// ContractConcept is a governed business object type in a contract.
type ContractConcept struct {
	ID         string             `json:"id"`
	Label      string             `json:"label"`
	Properties []ContractProperty `json:"properties"`
}

// BridgeContract is the versioned, signed semantic contract exported from a
// published ontology for downstream consumption.
type BridgeContract struct {
	ID          string            `json:"id"`
	Namespace   string            `json:"namespace"`
	Version     string            `json:"version"`
	Owner       string            `json:"owner,omitempty"`
	Concepts    []ContractConcept `json:"concepts"`
	Relations   []string          `json:"relations,omitempty"`
	Constraints []string          `json:"constraints,omitempty"`
	Units       []string          `json:"units,omitempty"`
	JSONSchema  string            `json:"jsonSchema,omitempty"`
	DDL         string            `json:"ddl,omitempty"`
	Signature   string            `json:"signature"`
	ExportedAt  time.Time         `json:"exportedAt"`
}

// ContractView is the result of a contract-driven query.
type ContractView struct {
	Contract   BridgeContract        `json:"contract"`
	Assertions []*KnowledgeAssertion `json:"assertions,omitempty"`
	AsOf       time.Time             `json:"asOf"`
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
	ID        string        `json:"id"`
	From      EntityRef     `json:"from"`
	To        EntityRef     `json:"to"`
	Kind      RelationKind  `json:"kind"`
	Evidence  []EvidenceRef `json:"evidence,omitempty"`
	Authority string        `json:"authority,omitempty"`
}

// MatchEvidence records why two references were matched.
type MatchEvidence struct {
	Left        EntityRef `json:"left"`
	Right       EntityRef `json:"right"`
	Score       float64   `json:"score"`
	MatchedOn   []string  `json:"matchedOn"`
	GeneratedBy string    `json:"generatedBy"`
	Status      string    `json:"status"`
	CapturedAt  time.Time `json:"capturedAt"`
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
	Hint      string          `json:"hint"`
	Canonical *EntityRef      `json:"canonical,omitempty"`
	Matches   []ResolvedMatch `json:"matches"`
	Related   []*Relation     `json:"related,omitempty"`
	Gaps      []string        `json:"gaps,omitempty"`
	AsOf      time.Time       `json:"asOf"`
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
	Reason     string                `json:"reason,omitempty"`
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
	Intent         Intent              `json:"intent"`
	View           *KnowledgeView      `json:"view"`
	Relations      []*Relation         `json:"relations,omitempty"`
	Paths          []*Path             `json:"paths,omitempty"`
	EvidenceChains []EvidenceChain     `json:"evidenceChains,omitempty"`
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
