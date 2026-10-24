# Registry Publication

GitHub Releases are the current distribution channel. Public registry
publication (PyPI / npm / Maven Central) is prepared but requires credentials.

## Status

| Registry | Artifact | Status |
| --- | --- | --- |
| PyPI | `gnosivela` wheel/sdist | Script ready; token required |
| npm | `@axisvela/gnosivela` | Script ready; npm auth required |
| Maven Central | `com.axisvela:gnosivela` | Not configured; `distributionManagement` + settings required |

## Dry run

Builds all packages without publishing:

```bash
./scripts/publish-registry.sh --dry-run --target all
```

## PyPI

```bash
export TWINE_USERNAME=__token__
export TWINE_PASSWORD=...
./scripts/publish-registry.sh --target pypi
```

## npm

```bash
npm login          # or export NODE_AUTH_TOKEN with a publish token
./scripts/publish-registry.sh --target npm
```

## Maven Central

Maven publication is intentionally blocked until `sdk/java/pom.xml` declares
`distributionManagement` (and the usual `licenses`, `developers`, `scm`
metadata) and `~/.m2/settings.xml` holds the Sonatype credentials. GPG signing
is mandatory for Maven Central.

```bash
./scripts/publish-registry.sh --target maven
```

## Verification after publication

```bash
python -m pip download --no-deps gnosivela==1.6.0 -d /tmp/gnosivela-pypi
npm view @axisvela/gnosivela@1.6.0 dist.integrity
mvn dependency:get -Dartifact=com.axisvela:gnosivela:1.6.0
```

Compare hashes with the GitHub Release `sha256sums.txt`.

## Security

- Never commit tokens, `.pypirc`, `.npmrc` with auth tokens, or Maven settings.
- Use CI/secret-manager scoped credentials.
- Publish only from a clean tag; the script checks `VERSION` against all SDK
  metadata before building.
