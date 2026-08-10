/** Core DTOs for the GNOSIVELA Semantic Control Plane. */

export interface EntityRef {
  namespace: string;
  canonicalId: string;
  type?: string;
  aliases?: string[];
  authority?: string;
}

export interface AssertionContext {
  region?: string;
  domain?: string;
  contract?: string;
  scenario?: string;
  purpose?: string;
  tags?: string[];
}

export interface Value {
  type: string; // string | number | boolean | date | entity
  string?: string;
  number?: number;
  boolean?: boolean;
  entityRef?: EntityRef;
}

export interface EvidenceRef {
  source: string;
  locator?: string;
  artifactDigest?: string;
  principal?: string;
}

export interface KnowledgeAssertion {
  assertionId: string;
  subject: EntityRef;
  predicate: string;
  object: Value;
  source?: string;
  status?: string; // proposed | validated | contested | superseded | revoked
  confidence?: number;
  context?: AssertionContext;
  evidence?: EvidenceRef[];
  validFrom?: string;
  validTo?: string;
  recordedAt?: string;
  authority?: string;
  ontologyVersion?: string;
  policyTags?: string[];
}

export interface Conflict {
  reason?: string;
  subject: EntityRef;
  predicate: string;
  assertions: KnowledgeAssertion[];
}

export interface PolicyRequest {
  principal?: string;
  action: string;
  resource?: string;
  tenant?: string;
  region?: string;
  attributes?: Record<string, string>;
}

export interface PolicyDecision {
  allowed: boolean;
  reason: string;
  policies?: string[];
}

export interface PipelineRequest {
  namespace: string;
  version: string;
  principal?: string;
  approvalId?: string;
  requireApproval?: boolean;
  maxFailures?: number;
  autoRollback?: boolean;
}

export interface EventContract {
  id: string;
  type: string;
  version?: string;
  ontology?: string;
  templates: EventTemplate[];
}

export interface EventTemplate {
  predicate: string;
  subjectField: string;
  subjectType?: string;
  objectField: string;
  objectType: string; // string | number | boolean | date
  region?: string;
  confidence?: number;
  validFor?: string;
  status?: string;
}

export interface Event {
  id: string;
  type: string;
  at?: string;
  source?: string;
  payload: Record<string, unknown>;
}
