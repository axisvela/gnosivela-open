package gnosivela;

/**
 * Quickstart for the GNOSIVELA Java SDK against a running control plane.
 *
 * <pre>
 *   cd backend && go run ./cmd/gnosivela     # start the server
 *   cd sdk/java && mvn -q compile            # or: javac -d target $(find src -name '*.java')
 * </pre>
 */
public final class Quickstart {
    public static void main(String[] args) {
        String base = System.getenv().getOrDefault("GNOSIVELA_URL", "http://localhost:8080");
        GnosivelaClient c = new GnosivelaClient(base);

        // 1. ontology from DSL
        String dsl = "ontology procurement.supplier @version 1.0\n"
                + "entity Supplier identifiedBy canonicalSupplierId\n"
                + "property riskScore: Decimal [0..1]\n";
        System.out.println("ontology: " + c.ontologyCreate(dsl));

        // 2. entity + assertion
        System.out.println("entity: " + c.entitySave(
                "{\"namespace\":\"mdm\",\"canonicalId\":\"S-1042\",\"type\":\"Supplier\",\"aliases\":[\"ACME\"]}"));
        System.out.println("assertion: " + c.assertionPropose(
                "{\"assertionId\":\"ka:risk\",\"subject\":{\"namespace\":\"mdm\",\"canonicalId\":\"S-1042\",\"type\":\"Supplier\"},"
                + "\"predicate\":\"risk:score\",\"object\":{\"type\":\"number\",\"number\":0.82},"
                + "\"source\":\"RiskOffice\",\"status\":\"validated\",\"confidence\":0.9}"));

        // 3. governance + operations
        System.out.println("consistency: " + c.consistencyReport("procurement.supplier"));
        System.out.println("quality: " + c.quality("procurement.supplier"));
        System.out.println("metrics: " + c.metrics());
        System.out.println("industry packs: " + c.industryPacks());

        // 4. semantic bridge contract
        System.out.println("contract: " + c.bridgeContractExport("procurement.supplier"));
    }
}
