package gnosivela;

import com.sun.net.httpserver.HttpServer;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.nio.charset.StandardCharsets;

/**
 * Smoke test for the Java SDK using the JDK HTTP server as a mock control
 * plane. No JUnit dependency: run with
 * {@code javac ... && java -cp ... gnosivela.GnosivelaClientTest}.
 */
public final class GnosivelaClientTest {
    private static int failures = 0;

    private static void check(boolean cond, String msg) {
        if (!cond) {
            failures++;
            System.out.println("FAIL: " + msg);
        } else {
            System.out.println("ok: " + msg);
        }
    }

    public static void main(String[] args) throws Exception {
        HttpServer server = HttpServer.create(new InetSocketAddress("127.0.0.1", 0), 0);
        server.createContext("/", exchange -> {
            String path = exchange.getRequestURI().getPath();
            String body = "{}";
            int status = 200;
            if ("/ontologies/acme/latest".equals(path)) {
                body = "{\"namespace\":\"acme\",\"version\":\"1.2\",\"status\":\"published\"}";
            } else if ("/consistency/conflicts".equals(path)) {
                body = "{\"conflicts\":[]}";
            } else if ("/events/contracts".equals(path) && "GET".equals(exchange.getRequestMethod())) {
                body = "{\"contracts\":[{\"id\":\"c-1\",\"type\":\"price.updated\"}]}";
            } else if ("/events/ingest".equals(path)) {
                body = "{\"assertions\":[{\"assertionId\":\"ev:e-1:p\"}],\"resolved\":[\"ev:e-1:p\"],\"gaps\":[]}";
                status = 201;
            } else if ("/entities".equals(path)) {
                body = "{\"namespace\":\"mdm\",\"canonicalId\":\"C-1\"}";
                status = 201;
            } else if ("/ontologies".equals(path) && "POST".equals(exchange.getRequestMethod())) {
                body = "{\"namespace\":\"acme\",\"version\":\"2.0\",\"status\":\"draft\"}";
                status = 201;
            } else if ("/audit/verify".equals(path)) {
                body = "{\"intact\":true,\"brokenAt\":-1}";
            } else if ("/audit".equals(path)) {
                body = "{\"entries\":[{\"id\":\"AUD-1\",\"action\":\"publish\",\"actor\":\"alice\"}]}";
            } else if (path.startsWith("/approval/requests/") && "GET".equals(exchange.getRequestMethod())) {
                body = "{\"id\":\"AR-001\",\"status\":\"pending\",\"steps\":[]}";
            } else if (path.endsWith("/reject")) {
                body = "{\"id\":\"AR-001\",\"status\":\"rejected\"}";
            } else if ("/entities/merge-candidate".equals(path)) {
                body = "{\"relation\":{\"id\":\"rel:1\",\"authority\":\"candidate\"},\"evidence\":{\"status\":\"auto\"}}";
                status = 201;
            } else if ("/grounding/explain".equals(path) || "/grounding/redact".equals(path)) {
                body = "{\"intent\":{\"query\":\"x\"},\"citations\":[]}";
            } else if ("/query/hybrid".equals(path)) {
                body = "{\"query\":\"supplier risk\",\"fusion\":\"rrf\",\"hits\":[{\"id\":\"ka:risk-1\",\"score\":0.5}]}";
            } else if ("/studio/model".equals(path)) {
                body = "{\"namespace\":\"acme\",\"version\":\"1.2\",\"concepts\":[{\"id\":\"Supplier\",\"assertionCount\":2}],\"relations\":[]}";
            } else if ("/studio/impact".equals(path)) {
                body = "{\"namespace\":\"acme\",\"fromVersion\":\"1.0\",\"toVersion\":\"1.1\",\"additive\":[{\"id\":\"Contract\"}],\"blocking\":false}";
            } else if (path.startsWith("/interop/shacl/")) {
                body = "@prefix sh: <http://www.w3.org/ns/shacl#> .\nex:SupplierShape a sh:NodeShape .\n";
            } else if ("/sparql".equals(path)) {
                body = "{\"head\":{\"vars\":[\"s\"]},\"results\":{\"bindings\":[{\"s\":{\"type\":\"uri\",\"value\":\"https://example/C-1\"}}]}}";
            } else if ("/actions".equals(path)) {
                body = "{\"actions\":[{\"id\":\"notify\",\"kind\":\"webhook\"}]}";
            } else if ("/actions/runs".equals(path)) {
                body = "{\"runs\":[{\"id\":\"ACT-1\",\"status\":\"succeeded\"}]}";
            } else if (path.startsWith("/actions/") && path.endsWith("/run")) {
                body = "{\"id\":\"ACT-1\",\"actionId\":\"notify\",\"status\":\"succeeded\"}";
            } else if ("/connectors".equals(path)) {
                body = "{\"connectors\":[{\"kind\":\"api\"},{\"kind\":\"streaming\"}]}";
            } else if ("/connectors/introspect".equals(path)) {
                body = "{\"kind\":\"streaming\",\"connector\":\"streaming\",\"fields\":[{\"name\":\"eventId\",\"dataType\":\"string\"}]}";
            } else if ("/connectors/documents".equals(path)) {
                body = "{\"id\":\"doc:p1\",\"docId\":\"p1\",\"bindingState\":\"bound\",\"chunks\":[{\"id\":\"c1\"}]}";
                status = 201;
            } else if ("/query/subgraph".equals(path)) {
                body = "[{\"id\":\"r1\",\"kind\":\"exactSameAs\"}]";
            } else {
                status = 404;
                body = "{\"error\":\"not found\"}";
            }
            byte[] out = body.getBytes(StandardCharsets.UTF_8);
            exchange.sendResponseHeaders(status, out.length);
            try (OutputStream os = exchange.getResponseBody()) {
                os.write(out);
            }
        });
        server.start();
        try {
            int port = server.getAddress().getPort();
            String base = "http://127.0.0.1:" + port;
            GnosivelaClient c = new GnosivelaClient(base);

            check(c.ontologyLatest("acme").contains("\"version\":\"1.2\""), "ontologyLatest returns 1.2");
            check(c.consistencyConflicts().contains("\"conflicts\":[]"), "consistencyConflicts empty");
            check(c.entitySave("{\"namespace\":\"mdm\",\"canonicalId\":\"C-1\"}").contains("\"canonicalId\":\"C-1\""), "entitySave roundtrip");
            check(c.eventContractList().contains("\"c-1\""), "eventContractList finds c-1");
            check(c.eventIngest("c-1", "{\"id\":\"e-1\",\"type\":\"price.updated\",\"payload\":{}}").contains("\"resolved\":[\"ev:e-1:p\"]"), "eventIngest resolves");

            // api-surface parity methods
            check(c.ontologyCreateJSON("{\"namespace\":\"acme\",\"version\":\"2.0\"}").contains("\"version\":\"2.0\""), "ontologyCreateJSON");
            check(c.approvalGet("AR-001").contains("\"status\":\"pending\""), "approvalGet");
            check(c.approvalReject("AR-001", "alice", "steward").contains("\"rejected\""), "approvalReject");
            check(c.auditListFiltered("alice", "publish", "").contains("\"publish\""), "auditListFiltered");
            check(c.auditVerify().contains("\"intact\":true"), "auditVerify");
            check(c.entityMergeCandidate(
                    "{\"namespace\":\"crm\",\"canonicalId\":\"ACME\"}",
                    "{\"namespace\":\"mdm\",\"canonicalId\":\"C-1042\"}",
                    "exactSameAs", "rule").contains("\"candidate\""), "entityMergeCandidate");
            check(c.groundingExplain("ACME", "alice", "onboarding").contains("\"intent\""), "groundingExplain");
            check(c.groundingRedact("ACME", "alice", "onboarding").contains("\"intent\""), "groundingRedact");
            check(c.subgraphQuery("capability:C-2").contains("\"r1\""), "subgraphQuery");
            check(c.hybridQuery("supplier risk", 5).contains("\"ka:risk-1\""), "hybridQuery");
            check(c.studioModel("acme").contains("\"Supplier\""), "studioModel");
            check(c.studioImpact("acme").contains("\"Contract\""), "studioImpact");
            check(c.shaclExport("acme").contains("ex:SupplierShape"), "shaclExport");
            check(c.sparqlQuery("SELECT ?s WHERE { ?s ?p ?o }").contains("\"vars\":[\"s\"]"), "sparqlQuery");
            check(c.actions().contains("\"notify\""), "actions");
            check(c.actionRegister("{\"id\":\"notify\",\"kind\":\"webhook\",\"url\":\"https://example.com\"}").contains("notify"), "actionRegister");
            check(c.actionRun("notify", "alice", "", "{\"x\":1}").contains("\"succeeded\""), "actionRun");
            check(c.actionRuns().contains("\"ACT-1\""), "actionRuns");
            check(c.connectors().contains("\"streaming\""), "connectors");
            check(c.connectorIntrospect("streaming", "{\"options\":{\"inlineSchema\":\"{}\"}}").contains("\"eventId\""), "connectorIntrospect");
            check(c.connectorBindDocument("p1", "must review").contains("\"doc:p1\""), "connectorBindDocument");

            try {
                c.ontologyLatest("missing");
                check(false, "404 must raise GnosivelaException");
            } catch (GnosivelaException e) {
                check(e.status() == 404, "404 raises GnosivelaException with status 404");
            }
        } finally {
            server.stop(0);
        }

        if (failures > 0) {
            System.out.println(failures + " failure(s)");
            System.exit(1);
        }
        System.out.println("all Java SDK tests passed");
    }
}
