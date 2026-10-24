#!/usr/bin/env python3
"""Validate the public ecosystem registry and compatibility matrix."""

from __future__ import annotations

import argparse
import hashlib
import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
CORE = ROOT.parent / "GNOSIVELA"
EE = ROOT.parent / "GNOSIVELA-ee"

VERSION_RE = re.compile(r"\d+\.\d+\.\d+")
RANGE_RE = re.compile(r">=(\d+\.\d+\.\d+),<(\d+\.\d+\.\d+)")
PARTNER_TYPES = {"connector", "industry_pack", "benchmark", "implementation"}
PARTNER_FIELDS = ("name", "type", "repository", "licence", "coreRange", "evidence")
DOCS = (
    "docs/ecosystem/README.md",
    "docs/ecosystem/partner-program.md",
    "docs/ecosystem/connector-certification.md",
    "docs/ecosystem/compatibility-matrix.md",
    "docs/ecosystem/benchmark-reproduction.md",
    "docs/ecosystem/registry-publication.md",
    "docs/ecosystem/compatibility.json",
    "docs/ecosystem/partners.json",
    "docs/ecosystem/partners.schema.json",
    ".github/ISSUE_TEMPLATE/partner-listing.yml",
    ".github/ISSUE_TEMPLATE/connector-certification.yml",
    ".github/ISSUE_TEMPLATE/benchmark-reproduction.yml",
    "scripts/publish-registry.sh",
)

errors: list[str] = []
warnings: list[str] = []


def read_text(path: Path) -> str:
    return path.read_text(encoding="utf-8")


def read_version(path: Path, label: str) -> str | None:
    if not path.exists():
        warnings.append(f"{label} not found at {path}; skipped")
        return None
    value = read_text(path).strip()
    if not VERSION_RE.fullmatch(value):
        errors.append(f"{label}: invalid version {value!r}")
        return None
    return value


def parse_version(value: str) -> tuple[int, int, int]:
    return tuple(int(part) for part in value.split("."))  # type: ignore[return-value]


def parse_range(value: str) -> tuple[tuple[int, int, int], tuple[int, int, int]] | None:
    match = RANGE_RE.fullmatch(value)
    if not match:
        return None
    return parse_version(match.group(1)), parse_version(match.group(2))


def in_range(version: str, value: str) -> bool:
    bounds = parse_range(value)
    if bounds is None:
        return False
    return bounds[0] <= parse_version(version) < bounds[1]


def sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def display_path(path: Path) -> str:
    try:
        return str(path.relative_to(ROOT))
    except ValueError:
        return str(path)


def load_json(path: Path) -> dict:
    try:
        return json.loads(read_text(path))
    except (OSError, json.JSONDecodeError) as exc:
        errors.append(f"{display_path(path)}: {exc}")
        return {}


def check_versions() -> tuple[str | None, str | None]:
    open_version = read_version(ROOT / "VERSION", "open VERSION")
    core_version = read_version(CORE / "VERSION", "core VERSION") if CORE.exists() else None
    ee_version = read_version(EE / "VERSION", "EE VERSION") if EE.exists() else None

    py_version = None
    pyproject = ROOT / "sdk/python/pyproject.toml"
    if pyproject.exists():
        match = re.search(r'^version = "([^"]+)"', read_text(pyproject), re.MULTILINE)
        py_version = match.group(1) if match else None
    ts_version = None
    package = ROOT / "sdk/typescript/package.json"
    if package.exists():
        ts_version = load_json(package).get("version")
    java_version = None
    pom = ROOT / "sdk/java/pom.xml"
    if pom.exists():
        match = re.search(r"<version>([^<]+)</version>", read_text(pom))
        java_version = match.group(1) if match else None

    compat = load_json(ROOT / "docs/ecosystem/compatibility.json")
    expected = {
        "open": open_version,
        "core": core_version,
        "ee": ee_version,
        "python": py_version,
        "typescript": ts_version,
        "java": java_version,
    }
    if compat:
        if compat.get("open") != open_version:
            errors.append(f"compatibility.json open={compat.get('open')} != VERSION={open_version}")
        if core_version and compat.get("core") != core_version:
            errors.append(f"compatibility.json core={compat.get('core')} != core VERSION={core_version}")
        if ee_version and compat.get("ee") != ee_version:
            errors.append(f"compatibility.json ee={compat.get('ee')} != EE VERSION={ee_version}")
        for key, actual in (("python", py_version), ("typescript", ts_version), ("java", java_version)):
            if actual and compat.get("sdks", {}).get(key) != actual:
                errors.append(f"compatibility.json sdk {key}={compat.get('sdks', {}).get(key)} != {actual}")

    if core_version and open_version and core_version != open_version:
        errors.append(f"open VERSION={open_version} != core VERSION={core_version}")

    if EE.exists():
        compat_file = EE / "CORE_COMPAT"
        if compat_file.exists() and compat:
            ee_compat = read_text(compat_file).strip()
            if compat.get("coreCompat") != ee_compat:
                errors.append(f"compatibility.json coreCompat={compat.get('coreCompat')} != EE CORE_COMPAT={ee_compat}")
            if core_version and not in_range(core_version, ee_compat):
                errors.append(f"core {core_version} does not satisfy EE CORE_COMPAT {ee_compat}")

    core_contract = CORE / "contracts/openapi.yaml"
    open_contract = ROOT / "contracts/openapi.yaml"
    if core_contract.exists() and open_contract.exists():
        if sha256(core_contract) != sha256(open_contract):
            errors.append("core and open OpenAPI contracts differ")

    for name, value in expected.items():
        if value is None:
            warnings.append(f"version {name} unavailable")
    return core_version, open_version


def check_partners(core_version: str | None, registry_path: Path) -> int:
    registry = load_json(registry_path)
    if not registry:
        return 0
    label_prefix = display_path(registry_path)
    if registry.get("schemaVersion") != 1:
        errors.append(f"{label_prefix} schemaVersion must be 1")
    partners = registry.get("partners")
    if not isinstance(partners, list):
        errors.append(f"{label_prefix} partners must be an array")
        return 0
    seen: set[str] = set()
    for index, partner in enumerate(partners):
        label = f"partners[{index}]"
        if not isinstance(partner, dict):
            errors.append(f"{label}: must be an object")
            continue
        for field in PARTNER_FIELDS:
            if not str(partner.get(field, "")).strip():
                errors.append(f"{label}: missing {field}")
        name = str(partner.get("name", "")).strip()
        if name in seen:
            errors.append(f"{label}: duplicate partner name {name}")
        seen.add(name)
        if partner.get("type") not in PARTNER_TYPES:
            errors.append(f"{label}: invalid type {partner.get('type')!r}")
        core_range = str(partner.get("coreRange", ""))
        if parse_range(core_range) is None:
            errors.append(f"{label}: invalid coreRange {core_range!r}")
        elif core_version and not in_range(core_version, core_range):
            errors.append(f"{label}: core {core_version} outside declared range {core_range}")
    return len(partners)


def check_docs() -> None:
    for rel in DOCS:
        if not (ROOT / rel).exists():
            errors.append(f"missing ecosystem file: {rel}")


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--partners-file",
        default="docs/ecosystem/partners.json",
        help="partner registry file to validate (default: docs/ecosystem/partners.json)",
    )
    args = parser.parse_args()
    registry_path = Path(args.partners_file)
    if not registry_path.is_absolute():
        registry_path = ROOT / registry_path

    check_docs()
    core_version, open_version = check_versions()
    partner_count = check_partners(core_version, registry_path)

    for warning in warnings:
        print(f"warning: {warning}")
    if errors:
        for error in errors:
            print(f"error: {error}")
        return 1
    print(
        "ecosystem gate passed: "
        f"core={core_version or 'n/a'} open={open_version or 'n/a'} "
        f"partners={partner_count} registry={display_path(registry_path)}"
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
