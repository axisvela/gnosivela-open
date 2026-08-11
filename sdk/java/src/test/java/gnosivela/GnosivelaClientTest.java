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
