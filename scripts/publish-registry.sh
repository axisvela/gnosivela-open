#!/usr/bin/env bash
# Publish GNOSIVELA SDK artifacts to public registries.
#
#   ./scripts/publish-registry.sh --dry-run --target all
#   TWINE_USERNAME=__token__ TWINE_PASSWORD=... ./scripts/publish-registry.sh --target pypi
#   npm login && ./scripts/publish-registry.sh --target npm
#
# Maven Central requires a configured distributionManagement and ~/.m2/settings.xml.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION="$(tr -d '[:space:]' < "$ROOT/VERSION")"
TARGET="all"
DRY_RUN=0

while [[ "$#" -gt 0 ]]; do
  case "$1" in
    --dry-run) DRY_RUN=1; shift ;;
    --target) TARGET="${2:-}"; shift 2 ;;
    --target=*) TARGET="${1#--target=}"; shift ;;
    -h|--help)
      sed -n '2,12p' "$0"
      exit 0
      ;;
    *)
      echo "unknown argument: $1" >&2
      exit 2
      ;;
  esac
done

case "$TARGET" in
  all|pypi|npm|maven) ;;
  *) echo "invalid target: $TARGET (use all|pypi|npm|maven)" >&2; exit 2 ;;
esac

PY_VERSION="$(grep -o '^version = ".*"' "$ROOT/sdk/python/pyproject.toml" | cut -d'"' -f2)"
TS_VERSION="$(grep -o '"version": "[^"]*"' "$ROOT/sdk/typescript/package.json" | cut -d'"' -f4)"
JAVA_VERSION="$(grep -o '<version>[^<]*</version>' "$ROOT/sdk/java/pom.xml" | head -1 | sed 's/<[^>]*>//g')"
for meta in "$PY_VERSION" "$TS_VERSION" "$JAVA_VERSION"; do
  [[ "$meta" == "$VERSION" ]] || {
    echo "version mismatch: VERSION=$VERSION SDK=$meta" >&2
    exit 1
  }
done

echo "==> building SDK artifacts for $VERSION"
(cd "$ROOT/sdk/python" && python -m build --outdir "$ROOT/sdk/python/dist")
(cd "$ROOT/sdk/typescript" && npm pack)
(cd "$ROOT/sdk/java" && mvn -q package -DskipTests)

publish_pypi() {
  if [[ "$DRY_RUN" == "1" ]]; then
    echo "==> [dry-run] PyPI: twine upload sdk/python/dist/*"
    return
  fi
  command -v python >/dev/null 2>&1 || { echo "python is required" >&2; exit 1; }
  python -m twine --version >/dev/null 2>&1 || {
    echo "twine is required: python -m pip install twine" >&2
    exit 1
  }
  if [[ -z "${TWINE_USERNAME:-}" || -z "${TWINE_PASSWORD:-}" ]]; then
    echo "PyPI credentials missing: set TWINE_USERNAME=__token__ and TWINE_PASSWORD" >&2
    exit 1
  fi
  python -m twine upload "$ROOT/sdk/python/dist/"*
}

publish_npm() {
  if [[ "$DRY_RUN" == "1" ]]; then
    echo "==> [dry-run] npm: npm publish --access public"
    return
  fi
  npm whoami >/dev/null 2>&1 || {
    echo "npm is not authenticated: run npm login or set NODE_AUTH_TOKEN" >&2
    exit 1
  }
  (cd "$ROOT/sdk/typescript" && npm publish --access public)
}

publish_maven() {
  if ! grep -q '<distributionManagement>' "$ROOT/sdk/java/pom.xml"; then
    echo "Maven Central is not configured: add distributionManagement to sdk/java/pom.xml" >&2
    echo "and credentials to ~/.m2/settings.xml; see docs/ecosystem/registry-publication.md" >&2
    [[ "$DRY_RUN" == "1" ]] && return
    exit 1
  fi
  if [[ "$DRY_RUN" == "1" ]]; then
    echo "==> [dry-run] Maven: mvn deploy -DskipTests"
    return
  fi
  (cd "$ROOT/sdk/java" && mvn -q deploy -DskipTests)
}

case "$TARGET" in
  all)
    publish_pypi
    publish_npm
    publish_maven
    ;;
  pypi) publish_pypi ;;
  npm) publish_npm ;;
  maven) publish_maven ;;
esac

if [[ "$DRY_RUN" == "1" ]]; then
  echo "==> dry-run complete for $VERSION (nothing published)"
else
  echo "==> registry publication complete for $VERSION"
fi
