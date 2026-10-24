#!/usr/bin/env bash
# Publish GNOSIVELA SDK artifacts to public registries.
#
#   ./scripts/publish-registry.sh --check --target all
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
CHECK=0

while [[ "$#" -gt 0 ]]; do
  case "$1" in
    --check) CHECK=1; shift ;;
    --dry-run) DRY_RUN=1; shift ;;
    --target) TARGET="${2:-}"; shift 2 ;;
    --target=*) TARGET="${1#--target=}"; shift ;;
    -h|--help)
      sed -n '2,14p' "$0"
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

check_pypi() {
  local failed=0
  if python -m twine --version >/dev/null 2>&1; then
    echo "pypi: twine ok"
  else
    echo "pypi: twine missing (python -m pip install twine)"
    failed=1
  fi
  if [[ -n "${TWINE_USERNAME:-}" && -n "${TWINE_PASSWORD:-}" ]]; then
    echo "pypi: credentials ok"
  else
    echo "pypi: credentials missing (TWINE_USERNAME / TWINE_PASSWORD)"
    failed=1
  fi
  return "$failed"
}

check_npm() {
  if npm whoami >/dev/null 2>&1; then
    echo "npm: authenticated as $(npm whoami)"
  else
    echo "npm: not authenticated (npm login or NODE_AUTH_TOKEN)"
    return 1
  fi
}

check_maven() {
  local failed=0
  if grep -q '<distributionManagement>' "$ROOT/sdk/java/pom.xml"; then
    echo "maven: distributionManagement present"
  else
    echo "maven: distributionManagement missing in sdk/java/pom.xml"
    failed=1
  fi
  if [[ -f "$HOME/.m2/settings.xml" ]]; then
    echo "maven: ~/.m2/settings.xml present"
  else
    echo "maven: ~/.m2/settings.xml missing"
    failed=1
  fi
  return "$failed"
}

if [[ "$CHECK" == "1" ]]; then
  failed=0
  case "$TARGET" in
    all)
      check_pypi || failed=1
      check_npm || failed=1
      check_maven || failed=1
      ;;
    pypi) check_pypi || failed=1 ;;
    npm) check_npm || failed=1 ;;
    maven) check_maven || failed=1 ;;
  esac
  if [[ "$failed" != "0" ]]; then
    echo "==> registry check failed: credentials or configuration missing"
    exit 1
  fi
  echo "==> registry check passed for $VERSION"
  exit 0
fi

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
