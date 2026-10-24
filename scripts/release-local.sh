#!/usr/bin/env bash
# Local CI/release driver for GNOSIVELA-open SDK artifacts.
#
# Runs the four SDK test suites, builds Python/npm/Java artifacts, and optionally
# publishes them to the GitHub Release for the current tag. Registry publication
# (PyPI / npm / Maven Central) still requires registry credentials.
#
#   ./scripts/release-local.sh                  # test + package
#   ./scripts/release-local.sh --publish        # test + package + GitHub Release
#   ./scripts/release-local.sh --skip-tests --publish
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION="$(tr -d '[:space:]' < "$ROOT/VERSION")"
TAG="v$VERSION"
STAGE="$ROOT/.release/$VERSION"

PUBLISH=0
SKIP_TESTS=0
for arg in "$@"; do
  case "$arg" in
    --publish) PUBLISH=1 ;;
    --skip-tests) SKIP_TESTS=1 ;;
    -h|--help)
      sed -n '2,12p' "$0"
      exit 0
      ;;
    *)
      echo "unknown argument: $arg" >&2
      exit 2
      ;;
  esac
done

PY_VERSION="$(grep -o '^version = ".*"' "$ROOT/sdk/python/pyproject.toml" | cut -d'"' -f2)"
TS_VERSION="$(grep -o '"version": "[^"]*"' "$ROOT/sdk/typescript/package.json" | cut -d'"' -f4)"
JAVA_VERSION="$(grep -o '<version>[^<]*</version>' "$ROOT/sdk/java/pom.xml" | head -1 | sed 's/<[^>]*>//g')"
for meta in "$PY_VERSION" "$TS_VERSION" "$JAVA_VERSION"; do
  [[ "$meta" == "$VERSION" ]] || {
    echo "version mismatch: VERSION=$VERSION SDK=$meta" >&2
    exit 1
  }
done

if [[ "$SKIP_TESTS" -eq 0 ]]; then
  echo "==> Go SDK tests"
  (cd "$ROOT" && GOWORK=off go build ./... && GOWORK=off go vet ./... && GOWORK=off go test ./...)
  echo "==> Python SDK tests"
  (cd "$ROOT/sdk/python" && PYTHONPATH=. python -m unittest discover -s tests -q)
  echo "==> TypeScript SDK tests"
  (cd "$ROOT/sdk/typescript" && npm test)
  echo "==> Java SDK tests"
  (cd "$ROOT/sdk/java" && mvn -q test)
fi

echo "==> packaging SDK $VERSION"
rm -rf "$STAGE"
mkdir -p "$STAGE"
(cd "$ROOT/sdk/python" && python -m build --outdir "$STAGE")
(cd "$ROOT/sdk/typescript" && npm pack --pack-destination "$STAGE")
(cd "$ROOT/sdk/java" && mvn -q package -DskipTests)
cp "$ROOT/sdk/java/target/gnosivela-$VERSION.jar" "$STAGE/"
cp "$ROOT/contracts/openapi.yaml" "$STAGE/gnosivela-openapi-$VERSION.yaml"
(cd "$STAGE" && sha256sum * > sha256sums.txt)

if [[ "$PUBLISH" -eq 0 ]]; then
  echo "==> artifacts ready in $STAGE (not published)"
  echo "    rerun with --publish to create/update the GitHub Release $TAG"
  exit 0
fi

command -v gh >/dev/null 2>&1 || { echo "gh is required for --publish" >&2; exit 1; }
REPO="${GNOSIVELA_OPEN_REPO:-axisvela/GNOSIVELA-open}"

ASSETS=("$STAGE"/*)
if gh release view "$TAG" -R "$REPO" >/dev/null 2>&1; then
  echo "==> updating existing release $TAG"
  gh release upload "$TAG" -R "$REPO" "${ASSETS[@]}" --clobber
  gh release edit "$TAG" -R "$REPO" --draft=false --latest
else
  echo "==> creating release $TAG"
  gh release create "$TAG" -R "$REPO" --title "GNOSIVELA-open $VERSION" \
    --generate-notes --latest "${ASSETS[@]}"
fi

echo "==> published $TAG (GitHub Release)"
echo "    PyPI / npm / Maven Central publication requires registry credentials"
