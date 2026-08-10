# GNOSIVELA Java SDK

Apache-2.0. Talks to the GNOSIVELA Semantic Control Plane over HTTP (never
links the AGPL core). Built on the JDK `HttpClient` only — no third-party
dependencies. Java 17+.

## Use

```java
GnosivelaClient c = new GnosivelaClient("http://localhost:8080");

String latest = c.ontologyLatest("procurement.supplier");      // raw JSON
String report = c.consistencyReport("procurement.supplier");
String view   = c.bridgeQuery("procurement.supplier", "ACME risk", "risk-officer", "onboarding");
```

Non-2xx responses raise `GnosivelaException` (carries `status()` and `body()`).

## Run the smoke test (no Maven/JUnit required)

```bash
cd sdk/java
mkdir -p target/classes
javac -d target/classes $(find src -name '*.java')
java -cp target/classes gnosivela.GnosivelaClientTest
```

## Build with Maven

```bash
mvn package
```

Optional dependency: Jackson (`jackson-databind`, marked optional) can be used
to map the JSON responses onto your own records.
