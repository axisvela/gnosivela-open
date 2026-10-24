# Connector Certification

Certification confirms that a connector is safe to use as a **metadata binder**
against a real source. It does not certify the source system itself.

## Levels

| Level | Meaning | Evidence |
| --- | --- | --- |
| Community Verified | A community member ran the real-source verification | Issue with redacted output and environment details |
| Certified | Maintainers reviewed code, tests, security and live evidence | Certification checklist + reproducible run |

## Checklist

- [ ] Implements the public `pkg/connector` SPI (`Kind`, `Introspect`).
- [ ] Returns `SourceSchema`; no vendor-specific types leak into the semantic API.
- [ ] Reads metadata only; source data is not copied or persisted.
- [ ] Handles authentication without embedding credentials.
- [ ] Network boundaries are explicit (private endpoints opt-in).
- [ ] Unit tests cover success, error and malformed responses.
- [ ] Live verification covers at least one real source.
- [ ] Documents required options, default ports, TLS and unsupported auth modes.
- [ ] Declares the supported core version range.
- [ ] Has a maintenance contact and licence.

## Evidence bundle

- connector repository and commit;
- core/open version used;
- source provider and version (no credentials);
- command and redacted output;
- known limitations (e.g. Kerberos/SASL unsupported in core);
- date and reviewer.

## Certification does not cover

- data correctness in the source;
- source-side access control quality;
- enterprise SSO/Kerberos/SASL packaging (EE scope);
- SLA or support guarantees beyond the partner's own statement.
