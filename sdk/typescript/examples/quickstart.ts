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
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
