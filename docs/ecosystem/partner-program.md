# Partner Program

The GNOSIVELA partner program lists **verified** connectors, industry packs,
benchmark reproductions and implementation services. Listing is evidence-based:
we link to source and test evidence rather than endorsing a vendor.

## Partner types

| Type | What is listed | Minimum evidence |
| --- | --- | --- |
| Connector Partner | A connector implementing the public `pkg/connector` SPI | Source repo, tests, live-source evidence, compatibility range |
| Industry Pack Partner | A neutral or commercial industry pack | Pack manifest, concepts/mappings, tests, licence |
| Benchmark Partner | An independent reproduction of GKF-Bench or scale-benchmark | Command, hardware, raw output, commit/version |
| Implementation Partner | Integration or delivery services | Public description, contact, case study (optional) |

## Requirements

- Public, versioned repository or downloadable artifact (commercial packs may be
  listed with a contact link instead of source).
- Compatible with the declared `CORE_COMPAT` range.
- No hidden copying of source data into GNOSIVELA.
- Security review: no embedded credentials, no SSRF bypass, explicit network
  boundaries.
- Clear licence and support/maintenance contact.
- Reproducible tests or a live verification report.

## Process

1. Open a **Partner / connector listing** issue with the template.
2. Maintainers review the evidence and may request a live verification run.
3. Accepted listings are added to this page with the supported version range.
4. Listings are re-reviewed on each core minor release.

## Current partners

**None yet.** This is expected for the P4-M4 scaffolding. Once two partner
connectors or packs are listed, the P4-M4 exit gate has its second condition.

## Removal

A listing is removed or marked stale if it fails a re-review, breaks the SPI,
drops below the supported core range, or is unmaintained for two minor releases.
