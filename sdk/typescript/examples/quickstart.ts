// GNOSIVELA TypeScript quickstart.
// Run with: node examples/quickstart.ts  (Node >= 22.6, no build step)
import { Client } from "../src/client.ts";

const c = new Client("http://localhost:8080");

async function main() {
  const dsl = `ontology procurement.supplier @version 1.0
entity Supplier identifiedBy canonicalSupplierId
property riskScore: Decimal [0..1]
`;

  const created = await c.ontologyCreate(dsl);
  console.log("ontology:", created.ontology?.namespace, created.ontology?.version);

  const ent = await c.entitySave({ namespace: "mdm", canonicalId: "S-1042", type: "Supplier", aliases: ["ACME"] });
  console.log("entity:", ent.canonicalId);

  await c.assertionPropose({
    assertionId: "ka:risk",
    subject: { namespace: "mdm", canonicalId: "S-1042", type: "Supplier" },
    predicate: "risk:score",
    object: { type: "number", number: 0.82 },
    source: "RiskOffice",
    status: "validated",
    confidence: 0.9,
  });

  const report = await c.consistencyReport("procurement.supplier");
  console.log("consistency failures:", report.failures);

  const conflicts = await c.consistencyConflicts();
  console.log("conflicts:", conflicts.conflicts.length);

  const view = await c.bridgeQuery("procurement.supplier", "ACME risk", "risk-officer", "onboarding");
  console.log("bridge view assertions:", (view.assertions as unknown[] | undefined)?.length ?? 0);

  const fed = await c.federationQuery("ACME risk", "risk-officer", "onboarding");
  console.log("federation domains:", fed.domainsQueried, "hits:", fed.domainHits);

  await c.eventContractRegister({
    id: "procurement.price.updated",
    type: "price.updated",
    ontology: "procurement.supplier@1.2",
    templates: [{
      predicate: "Supplier:price", subjectField: "company", subjectType: "Supplier",
      objectField: "amount", objectType: "number", region: "SG", validFor: "90d",
    }],
  });
  const ing = await c.eventIngest("procurement.price.updated", {
    id: "e-1", type: "price.updated", source: "market-feed",
    payload: { company: "ACME", amount: 12.5 },
  });
  console.log("event ingest assertions:", ing.assertions.length, "resolved:", ing.resolved.length);

  const q = await c.quality("");
  console.log("quality citation:", q.citationCompleteness, "conflicts:", q.conflicts);
  await c.incidentRuleAdd({ id: "r-conflicts", metric: "conflicts", operator: ">=", threshold: 0, severity: "warning" });
  const opened = await c.incidentCheck();
  console.log("incidents opened:", opened.opened.length);
  const counts = await c.metrics();
  console.log("metrics:", counts.counts);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
