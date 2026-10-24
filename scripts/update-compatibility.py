#!/usr/bin/env python3
"""Update or check docs/ecosystem/compatibility.json from repository versions."""

from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
CORE = ROOT.parent / "GNOSIVELA"
EE = ROOT.parent / "GNOSIVELA-ee"
TARGET = ROOT / "docs/ecosystem/compatibility.json"


def read_version(path: Path, label: str) -> str:
    try:
        value = path.read_text(encoding="utf-8").strip()
    except OSError as exc:
        raise SystemExit(f"{label}: {exc}")
    if not re.fullmatch(r"\d+\.\d+\.\d+", value):
        raise SystemExit(f"{label}: invalid version {value!r}")
    return value


def sdk_versions() -> dict[str, str]:
    pyproject = ROOT / "sdk/python/pyproject.toml"
    package = ROOT / "sdk/typescript/package.json"
    pom = ROOT / "sdk/java/pom.xml"
    py = re.search(r'^version = "([^"]+)"', pyproject.read_text(encoding="utf-8"), re.MULTILINE)
    ts = json.loads(package.read_text(encoding="utf-8")).get("version")
    java = re.search(r"<version>([^<]+)</version>", pom.read_text(encoding="utf-8"))
    if not py or not ts or not java:
        raise SystemExit("cannot read SDK versions")
    return {"go": py.group(1), "python": py.group(1), "typescript": ts, "java": java.group(1)}


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--check", action="store_true", help="exit 1 if the file is stale")
    args = parser.parse_args()

    if not TARGET.exists():
        raise SystemExit(f"missing {TARGET}")

    current = json.loads(TARGET.read_text(encoding="utf-8"))
    open_version = read_version(ROOT / "VERSION", "open VERSION")
    core_version = read_version(CORE / "VERSION", "core VERSION")
    ee_version = read_version(EE / "VERSION", "EE VERSION")
    ee_compat = (EE / "CORE_COMPAT").read_text(encoding="utf-8").strip()
    sdks = sdk_versions()

    updated = {
        "schemaVersion": current.get("schemaVersion", 1),
        "core": core_version,
        "open": open_version,
        "ee": ee_version,
        "coreCompat": ee_compat,
        "contract": open_version,
        "sdks": sdks,
        "runtimes": current.get("runtimes", {}),
        "frameworks": current.get("frameworks", {}),
    }

    if current == updated:
        print(f"compatibility.json is current (core/open {open_version}, EE {ee_version})")
        return 0

    if args.check:
        for key in updated:
            if current.get(key) != updated.get(key):
                print(f"stale {key}: {current.get(key)!r} -> {updated.get(key)!r}", file=sys.stderr)
        return 1

    TARGET.write_text(json.dumps(updated, indent=2) + "\n", encoding="utf-8")
    print(f"updated compatibility.json (core/open {open_version}, EE {ee_version})")
    return 0


if __name__ == "__main__":
    sys.exit(main())
