package gnosivela;

import java.io.IOException;
import java.net.URI;
import java.net.URLEncoder;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import java.time.Duration;

/**
 * Thin client for the GNOSIVELA Semantic Control Plane API. Built on the JDK
 * HttpClient only — no third-party dependencies. Methods return the raw JSON
 * response body; non-2xx responses raise {@link GnosivelaException}.
 */
public final class GnosivelaClient {
    private final String baseURL;
    private final HttpClient http;

    public GnosivelaClient(String baseURL) {
        this.baseURL = baseURL.replaceAll("/$", "");
        this.http = HttpClient.newBuilder().connectTimeout(Duration.ofSeconds(10)).build();
    }

    // ---- low level ----

    public String get(String path) {
        return request("GET", path, null);
    }

    public String post(String path, String jsonBody) {
        return request("POST", path, jsonBody);
    }

    private String request(String method, String path, String jsonBody) {
        HttpRequest.Builder b = HttpRequest.newBuilder()
                .uri(URI.create(baseURL + path))
                .timeout(Duration.ofSeconds(30));
        if (jsonBody != null) {
            b.header("Content-Type", "application/json")
                    .method(method, HttpRequest.BodyPublishers.ofString(jsonBody, StandardCharsets.UTF_8));
        } else {
            b.method(method, HttpRequest.BodyPublishers.noBody());
        }
        try {
            HttpResponse<String> resp = http.send(b.build(), HttpResponse.BodyHandlers.ofString());
            if (resp.statusCode() >= 300) {
                throw new GnosivelaException(resp.statusCode(), resp.body(), path);
            }
            return resp.body();
        } catch (IOException e) {
            throw new GnosivelaException(0, e.getMessage(), path);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new GnosivelaException(0, e.getMessage(), path);
        }
    }

    private String requestTextBody(String method, String path, String body, String contentType) {
        HttpRequest.Builder b = HttpRequest.newBuilder()
                .uri(URI.create(baseURL + path))
                .timeout(Duration.ofSeconds(30))
                .header("Content-Type", contentType)
                .method(method, HttpRequest.BodyPublishers.ofString(body, StandardCharsets.UTF_8));
        try {
            HttpResponse<String> resp = http.send(b.build(), HttpResponse.BodyHandlers.ofString());
            if (resp.statusCode() >= 300) {
                throw new GnosivelaException(resp.statusCode(), resp.body(), path);
            }
            return resp.body();
        } catch (IOException e) {
            throw new GnosivelaException(0, e.getMessage(), path);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new GnosivelaException(0, e.getMessage(), path);
        }
    }

    private static String enc(String s) {
        return URLEncoder.encode(s, StandardCharsets.UTF_8);
    }

    // ---- ontology ----

    public String ontologyCreate(String dsl) {
        return request("POST", "/ontologies", dsl);
    }

    public String ontologyCreateJSON(String json) {
        return request("POST", "/ontologies", json);
    }

    public String ontologyLatest(String namespace) {
        return get("/ontologies/" + enc(namespace) + "/latest");
    }

    public String ontologyGet(String namespace, String version) {
        return get("/ontologies/" + enc(namespace) + "/versions/" + enc(version));
    }

    public String ontologyPublish(String namespace, String version, String approval) {
        String q = approval == null || approval.isEmpty() ? "" : "?approval=" + enc(approval);
        return post("/ontologies/" + enc(namespace) + "/versions/" + enc(version) + "/publish" + q, null);
    }

    public String ontologyImpact(String namespace, String version) {
        return get("/ontologies/" + enc(namespace) + "/versions/" + enc(version) + "/impact");
    }

    public String ontologyRollback(String namespace, String version) {
        return post("/ontologies/" + enc(namespace) + "/versions/" + enc(version) + "/rollback", null);
    }

    public String ontologyDiff(String namespace, String version, String other) {
        return get("/ontologies/" + enc(namespace) + "/versions/" + enc(version) + "/diff?other=" + enc(other));
    }

    // ---- assertion / entity ----

    public String assertionPropose(String json) {
        return post("/assertions", json);
    }

    public String assertionList(String subjectNs, String subjectId) {
        return get("/assertions?subjectNs=" + enc(subjectNs) + "&subjectId=" + enc(subjectId));
    }

    public String entitySave(String json) {
        return post("/entities", json);
    }

    public String entityResolve(String hint) {
        return post("/entities/resolve", "{\"hint\":\"" + escape(hint) + "\"}");
    }

    public String entityExplain(String hint) {
        return post("/entities/explain", "{\"hint\":\"" + escape(hint) + "\"}");
    }

    // ---- query / grounding ----

    public String semanticQuery(String query, String principal, String purpose) {
        return post("/query/semantic", "{\"query\":\"" + escape(query) + "\",\"principal\":\"" + escape(principal) + "\",\"purpose\":\"" + escape(purpose) + "\"}");
    }

    public String hybridQuery(String query, int limit) {
        return post("/query/hybrid", "{\"query\":\"" + escape(query) + "\",\"limit\":" + limit + "}");
    }

    public String studioModel(String namespace) {
        return get("/studio/model?namespace=" + enc(namespace));
    }

    public String studioImpact(String namespace) {
        return get("/studio/impact?namespace=" + enc(namespace));
    }

    public String actions() {
        return get("/actions");
    }

    public String actionRegister(String definitionJson) {
        return post("/actions", definitionJson);
    }

    public String actionRun(String actionId, String principal, String approvalId, String inputJson) {
        String input = (inputJson == null || inputJson.isEmpty()) ? "{}" : inputJson;
        return post("/actions/" + enc(actionId) + "/run", "{\"principal\":\"" + escape(principal) +
                "\",\"approvalId\":\"" + escape(approvalId) + "\",\"input\":" + input + "}");
    }

    public String actionRuns() {
        return get("/actions/runs");
    }

    public String sparqlQuery(String query) {
        return post("/sparql", "{\"query\":\"" + escape(query) + "\"}");
    }

    public String sparqlQueryText(String query) {
        return requestTextBody("POST", "/sparql", query, "application/sparql-query");
    }

    public String shaclExport(String namespace) {
        return get("/interop/shacl/" + enc(namespace));
    }

    public String shaclExport(String namespace, String version) {
        return get("/interop/shacl/" + enc(namespace) + "?version=" + enc(version));
    }

    public String shaclImport(String namespace, String turtle) {
        return requestTextBody("POST", "/interop/shacl/" + enc(namespace), turtle, "text/turtle");
    }

    public String shaclImport(String namespace, String version, String turtle) {
        return requestTextBody("POST", "/interop/shacl/" + enc(namespace) + "?version=" + enc(version), turtle, "text/turtle");
    }

    public String shaclValidateStored(String namespace) {
        return get("/interop/shacl/" + enc(namespace) + "/validate");
    }

    public String shaclValidateStored(String namespace, String version) {
        return get("/interop/shacl/" + enc(namespace) + "/validate?version=" + enc(version));
    }

    public String shaclValidate(String namespace, String assertionsJson) {
        return post("/interop/shacl/" + enc(namespace) + "/validate", "{\"assertions\":" + assertionsJson + "}");
    }

    public String shaclValidate(String namespace, String version, String assertionsJson) {
        return post("/interop/shacl/" + enc(namespace) + "/validate?version=" + enc(version),
                "{\"assertions\":" + assertionsJson + "}");
    }

    public String owlExport(String namespace) {
        return get("/interop/owl/" + enc(namespace));
    }

    public String owlExport(String namespace, String version) {
        return get("/interop/owl/" + enc(namespace) + "?version=" + enc(version));
    }

    public String owlInfer(String namespace) {
        return get("/interop/owl/" + enc(namespace) + "/infer");
    }

    public String owlInfer(String namespace, String version) {
        return get("/interop/owl/" + enc(namespace) + "/infer?version=" + enc(version));
    }

    public String graphqlQuery(String query) {
        return post("/graphql", "{\"query\":\"" + escape(query) + "\"}");
    }

    public String graphqlQuery(String query, String variablesJson) {
        return post("/graphql", "{\"query\":\"" + escape(query) + "\",\"variables\":" + variablesJson + "}");
    }

    public String graphqlSchema() {
        return get("/graphql/schema");
    }

    public String connectors() {
        return get("/connectors");
    }

    public String connectorIntrospect(String kind, String specJson) {
        return post("/connectors/introspect", "{\"kind\":\"" + escape(kind) + "\",\"spec\":" + specJson + "}");
    }

    public String connectorBindDocument(String docId, String content) {
        return connectorBindDocument(docId, content, "markdown");
    }

    public String connectorBindDocument(String docId, String content, String contentType) {
        return post("/connectors/documents", "{\"docId\":\"" + escape(docId) + "\",\"contentType\":\"" + escape(contentType) +
                "\",\"content\":\"" + escape(content) + "\"}");
    }

    public String columnLineage(String schemaJson, String mappingJson) {
        return post("/lineage/columns", "{\"schema\":" + schemaJson + ",\"mapping\":" + mappingJson + "}");
    }

    public String pathQuery(String from, String to) {
        return post("/query/path", "{\"from\":\"" + escape(from) + "\",\"to\":\"" + escape(to) + "\"}");
    }

    public String groundingAssemble(String query, String principal, String purpose) {
        return post("/grounding/assemble", "{\"query\":\"" + escape(query) + "\",\"principal\":\"" + escape(principal) + "\",\"purpose\":\"" + escape(purpose) + "\"}");
    }

    public String groundingExplain(String query, String principal, String purpose) {
        return post("/grounding/explain", "{\"query\":\"" + escape(query) + "\",\"principal\":\"" + escape(principal) + "\",\"purpose\":\"" + escape(purpose) + "\"}");
    }

    public String groundingRedact(String query, String principal, String purpose) {
        return post("/grounding/redact", "{\"query\":\"" + escape(query) + "\",\"principal\":\"" + escape(principal) + "\",\"purpose\":\"" + escape(purpose) + "\"}");
    }

    public String subgraphQuery(String node) {
        return post("/query/subgraph", "{\"node\":\"" + escape(node) + "\"}");
    }

    // ---- governance ----

    public String consistencyReport(String ontologyNamespace) {
        String q = ontologyNamespace == null || ontologyNamespace.isEmpty() ? "" : "?ontology=" + enc(ontologyNamespace);
        return get("/consistency/report" + q);
    }

    public String consistencyConflicts() {
        return get("/consistency/conflicts");
    }

    public String consistencyResolve() {
        return post("/consistency/resolve", null);
    }

    public String consistencyAudit() {
        return get("/consistency/audit");
    }

    public String policyList() {
        return get("/policy/policies");
    }

    public String policyEvaluate(String json) {
        return post("/policy/evaluate", json);
    }

    public String approvalCreate(String action, String resource, String requester) {
        return post("/approval/requests", "{\"action\":\"" + escape(action) + "\",\"resource\":\"" + escape(resource) + "\",\"requester\":\"" + escape(requester) + "\"}");
    }

    public String approvalList() {
        return get("/approval/requests");
    }

    public String approvalApprove(String id, String approver, String role) {
        return post("/approval/requests/" + enc(id) + "/approve", "{\"approver\":\"" + escape(approver) + "\",\"role\":\"" + escape(role) + "\"}");
    }

    public String approvalGet(String id) {
        return get("/approval/requests/" + enc(id));
    }

    public String approvalReject(String id, String approver, String role) {
        return post("/approval/requests/" + enc(id) + "/reject", "{\"approver\":\"" + escape(approver) + "\",\"role\":\"" + escape(role) + "\"}");
    }

    public String auditList() {
        return get("/audit");
    }

    public String auditListFiltered(String actor, String action, String resource) {
        StringBuilder q = new StringBuilder();
        if (actor != null && !actor.isEmpty()) {
            q.append("actor=").append(enc(actor));
        }
        if (action != null && !action.isEmpty()) {
            if (q.length() > 0) q.append("&");
            q.append("action=").append(enc(action));
        }
        if (resource != null && !resource.isEmpty()) {
            if (q.length() > 0) q.append("&");
            q.append("resource=").append(enc(resource));
        }
        return get("/audit" + (q.length() > 0 ? "?" + q : ""));
    }

    public String auditVerify() {
        return get("/audit/verify");
    }

    public String entityMergeCandidate(String leftJson, String rightJson, String kind, String generatedBy) {
        return post("/entities/merge-candidate", "{\"left\":" + leftJson + ",\"right\":" + rightJson + ",\"kind\":\"" + escape(kind) + "\",\"generatedBy\":\"" + escape(generatedBy) + "\"}");
    }

    public String auditAttest(String entryId, String ref, String content, String by) {
        return post("/audit/attest", "{\"entryId\":\"" + escape(entryId) + "\",\"ref\":\"" + escape(ref) + "\",\"content\":\"" + escape(content) + "\",\"by\":\"" + escape(by) + "\"}");
    }

    // ---- pipeline / federation / bridge / events ----

    public String pipelineRun(String json) {
        return post("/pipeline/run", json);
    }

    public String federationQuery(String query, String principal, String purpose) {
        return post("/federation/query", "{\"query\":\"" + escape(query) + "\",\"principal\":\"" + escape(principal) + "\",\"purpose\":\"" + escape(purpose) + "\"}");
    }

    public String federationDomains() {
        return get("/federation/domains");
    }

    public String federationDomainAdd(String name, String baseUrl) {
        return post("/federation/domains", "{\"name\":\"" + escape(name) + "\",\"baseUrl\":\"" + escape(baseUrl) + "\"}");
    }

    public String bridgeContractExport(String namespace) {
        return get("/bridge/" + enc(namespace) + "/contract");
    }

    public String bridgeQuery(String namespace, String query, String principal, String purpose) {
        return post("/bridge/query", "{\"namespace\":\"" + escape(namespace) + "\",\"query\":\"" + escape(query) + "\",\"principal\":\"" + escape(principal) + "\",\"purpose\":\"" + escape(purpose) + "\"}");
    }

    public String eventContractRegister(String json) {
        return post("/events/contracts", json);
    }

    public String eventContractList() {
        return get("/events/contracts");
    }

    public String eventIngest(String contractId, String eventJson) {
        return post("/events/ingest", "{\"contractId\":\"" + escape(contractId) + "\",\"event\":" + eventJson + "}");
    }

    public String metrics() {
        return get("/metrics");
    }

    public String quality(String ontologyNamespace) {
        String q = ontologyNamespace == null || ontologyNamespace.isEmpty() ? "" : "?ontology=" + enc(ontologyNamespace);
        return get("/quality" + q);
    }

    public String incidentRuleAdd(String json) {
        return post("/incidents/rules", json);
    }

    public String incidentCheck() {
        return post("/incidents/check", null);
    }

    public String incidentList() {
        return get("/incidents");
    }

    public String incidentResolve(String id) {
        return post("/incidents/" + enc(id) + "/resolve", null);
    }

    public String metricDefinitions() {
        return get("/metrics/definitions");
    }

    public String metricDefinitionRegister(String json) {
        return post("/metrics/definitions", json);
    }

    public String industryPacks() {
        return get("/industry/packs");
    }

    public String industryPack(String id) {
        return get("/industry/packs/" + enc(id));
    }

    private static String escape(String s) {
        if (s == null) {
            return "";
        }
        return s.replace("\\", "\\\\").replace("\"", "\\\"");
    }
}
