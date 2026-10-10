#!/usr/bin/env python3
"""
Webull API Live Probe — v0.5 Schema Discovery

Makes authenticated real API calls against the Webull HK sandbox to discover
correct endpoint paths and response schemas for Fund Data, Crypto Data,
Screener v2, and Broker FD endpoints.

Token is obtained via a Go helper binary (to use the SDK's own signing logic).
All subsequent API calls are made directly in Python.

Usage:
    python examples/probe/probe.py
"""

import os
import sys
import json
import time
import hmac
import hashlib
import base64
import urllib.request
import urllib.parse
import urllib.error
import subprocess
import secrets
from datetime import datetime, timezone

SCRIPT_DIR = os.path.dirname(__file__)
GETTOKEN_BIN = os.path.join(SCRIPT_DIR, "gettoken")

APP_KEY = os.environ.get("WEBULL_APP_KEY", "4b2b7acd2bf0d30d8aea173fceefa238")
APP_SECRET = os.environ.get("WEBULL_APP_SECRET", "840b4353a6a31ce3ab91e2f99a510272")
BASE_URL = "https://api.sandbox.webull.hk"
BROKER_API_URL = "https://broker-api.sandbox.webull.hk"
RESULTS_DIR = os.path.join(SCRIPT_DIR, "results")
REQUEST_COOLDOWN = 3
TOKEN_CACHE_TTL = 300  # 5 min, tokens valid ~2 weeks

os.makedirs(RESULTS_DIR, exist_ok=True)

_token_cache = {"token": "", "ts": 0}


def get_token() -> str:
    """Get a fresh access token, caching for TOKEN_CACHE_TTL seconds."""
    now = time.time()
    if _token_cache["token"] and (now - _token_cache["ts"]) < TOKEN_CACHE_TTL:
        return _token_cache["token"]

    # Prefer WEBULL_TOKEN env var if set (avoids Go binary + rate limiting)
    if os.environ.get("WEBULL_TOKEN"):
        _token_cache["token"] = os.environ["WEBULL_TOKEN"]
        _token_cache["ts"] = now
        print(f"[TOKEN] Using WEBULL_TOKEN env var")
        return _token_cache["token"]

    print(f"[TOKEN] Requesting new token via Go binary...")
    env = dict(os.environ)
    env["WEBULL_APP_KEY"] = APP_KEY
    env["WEBULL_APP_SECRET"] = APP_SECRET
    try:
        result = subprocess.run(
            [GETTOKEN_BIN + ".exe"],
            env=env,
            capture_output=True,
            text=True,
            timeout=30,
        )
        if result.returncode != 0:
            print(f"[TOKEN] Go binary error: {result.stderr.strip()}")
            sys.exit(1)
        _token_cache["token"] = result.stdout.strip()
        _token_cache["ts"] = now
        print(f"[TOKEN] Got: {_token_cache['token'][:20]}... (cached {TOKEN_CACHE_TTL}s)")
        return _token_cache["token"]
    except FileNotFoundError:
        fallback = os.environ.get("WEBULL_TOKEN")
        if fallback:
            print("[TOKEN] Go binary not found, using WEBULL_TOKEN env var")
            return fallback
        print("[TOKEN] Go binary not found and WEBULL_TOKEN not set")
        sys.exit(1)


def md5_hex(data: bytes) -> str:
    return hashlib.md5(data).hexdigest().upper()


def percent_encode(s: str) -> str:
    unreserved = frozenset(
        b"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_.~"
    )
    parts = []
    for c in s.encode("utf-8"):
        if c in unreserved:
            parts.append(chr(c))
        else:
            parts.append(f"%{c:02X}")
    return "".join(parts)


def get_token() -> str:
    """Get a fresh access token using the Go binary (SDK's signing)."""
    env = dict(os.environ)
    env["WEBULL_APP_KEY"] = APP_KEY
    env["WEBULL_APP_SECRET"] = APP_SECRET
    try:
        result = subprocess.run(
            [GETTOKEN_BIN + ".exe"],
            env=env,
            capture_output=True,
            text=True,
            timeout=30,
        )
        if result.returncode != 0:
            print(f"[TOKEN] Go binary error: {result.stderr}")
            sys.exit(1)
        return result.stdout.strip()
    except FileNotFoundError:
        fallback = os.environ.get("WEBULL_TOKEN")
        if fallback:
            print("[TOKEN] Go binary not found, using WEBULL_TOKEN env var")
            return fallback
        print("[TOKEN] Go binary not found and WEBULL_TOKEN not set")
        sys.exit(1)


def sign(method: str, path: str, query: dict, body: bytes) -> tuple:
    """Compute HMAC-SHA1 signature. Returns (signature, nonce, timestamp)."""
    now = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
    nonce = secrets.token_hex(16)
    host = urllib.parse.urlparse(BASE_URL).netloc

    sig_headers = {
        "x-app-key": APP_KEY,
        "x-signature-algorithm": "HMAC-SHA1",
        "x-signature-version": "1.0",
        "x-signature-nonce": nonce,
        "x-timestamp": now,
        "host": host,
    }

    all_vals = {}
    for k, v in query.items():
        all_vals[k] = v if isinstance(v, list) else [v]
    for k, v in sig_headers.items():
        all_vals[k] = v if isinstance(v, list) else [v]

    entries = []
    for k in sorted(all_vals):
        vals = sorted(all_vals[k])
        entries.append(f"{k}={','.join(vals)}")

    if body:
        canonical = f"{path}&{'&'.join(entries)}&{md5_hex(body)}"
    else:
        canonical = f"{path}&{'&'.join(entries)}"

    encoded = percent_encode(canonical)
    mac = hmac.new(
        (APP_SECRET + "&").encode("utf-8"),
        encoded.encode("utf-8"),
        hashlib.sha1,
    )
    return base64.b64encode(mac.digest()).decode("utf-8"), nonce, now


def api_get(path: str, query: dict = None, category: str = "US") -> dict:
    """Make an authenticated GET request."""
    time.sleep(REQUEST_COOLDOWN)
    q = dict(query or {})
    if category:
        q["category"] = category

    body = b""
    sig, nonce, ts = sign("GET", path, q, body)
    token = get_token()

    headers = {
        "x-app-key": APP_KEY,
        "x-signature-algorithm": "HMAC-SHA1",
        "x-signature-version": "1.0",
        "x-signature-nonce": nonce,
        "x-timestamp": ts,
        "host": urllib.parse.urlparse(BASE_URL).netloc,
        "Authorization": f"Bearer {token}",
        "x-signature": sig,
    }

    url = BASE_URL + path
    if q:
        url += "?" + urllib.parse.urlencode(q)

    print(f"[GET ] {path}  q={q}")
    req = urllib.request.Request(url, headers=headers, method="GET")

    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            data = json.loads(resp.read().decode("utf-8"))
            print(f"[GET ] {path} -> HTTP {resp.status}")
            return {"status": resp.status, "data": data}
    except urllib.error.HTTPError as e:
        body_text = e.read().decode("utf-8") if e.fp else ""
        print(f"[GET ] {path} -> HTTP {e.code}")
        return {"status": e.code, "error": body_text[:300]}


def api_post(path: str, query: dict = None, body: dict = None) -> dict:
    """Make an authenticated POST request."""
    time.sleep(REQUEST_COOLDOWN)
    q = dict(query or {})
    body_bytes = json.dumps(body or {}).encode("utf-8")

    sig, nonce, ts = sign("POST", path, q, body_bytes)
    token = get_token()

    headers = {
        "x-app-key": APP_KEY,
        "x-signature-algorithm": "HMAC-SHA1",
        "x-signature-version": "1.0",
        "x-signature-nonce": nonce,
        "x-timestamp": ts,
        "host": urllib.parse.urlparse(BASE_URL).netloc,
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json",
        "x-signature": sig,
    }

    url = BASE_URL + path
    if q:
        url += "?" + urllib.parse.urlencode(q)

    print(f"[POST] {path}  body={body or {}}")
    req = urllib.request.Request(url, data=body_bytes, headers=headers, method="POST")

    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            data = json.loads(resp.read().decode("utf-8"))
            print(f"[POST] {path} -> HTTP {resp.status}")
            return {"status": resp.status, "data": data}
    except urllib.error.HTTPError as e:
        body_text = e.read().decode("utf-8") if e.fp else ""
        print(f"[POST] {path} -> HTTP {e.code}")
        return {"status": e.code, "error": body_text[:300]}


def broker_get(path: str, query: dict = None) -> dict:
    """Make an authenticated GET request to the Broker API host."""
    time.sleep(REQUEST_COOLDOWN)
    q = dict(query or {})

    body = b""
    sig, nonce, ts = sign("GET", path, q, body)
    token = get_token()

    headers = {
        "x-app-key": APP_KEY,
        "x-signature-algorithm": "HMAC-SHA1",
        "x-signature-version": "1.0",
        "x-signature-nonce": nonce,
        "x-timestamp": ts,
        "host": urllib.parse.urlparse(BROKER_API_URL).netloc,
        "Authorization": f"Bearer {token}",
        "x-signature": sig,
    }

    url = BROKER_API_URL + path
    if q:
        url += "?" + urllib.parse.urlencode(q)

    print(f"[BRK ] {path}  q={q}")
    req = urllib.request.Request(url, headers=headers, method="GET")

    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            data = json.loads(resp.read().decode("utf-8"))
            print(f"[BRK ] {path} -> HTTP {resp.status}")
            return {"status": resp.status, "data": data}
    except urllib.error.HTTPError as e:
        body_text = e.read().decode("utf-8") if e.fp else ""
        print(f"[BRK ] {path} -> HTTP {e.code}")
        return {"status": e.code, "error": body_text[:300]}


def save(name: str, info: dict, result: dict):
    path = os.path.join(RESULTS_DIR, name)
    with open(path, "w", encoding="utf-8") as f:
        json.dump({**info, "result": result}, f, indent=2, ensure_ascii=False)
    print(f"[SAVE] {path}")


def tag_from_path(path: str) -> str:
    return path.lstrip("/").replace("/", "_").replace("-", "_").replace("?", "_").replace("&", "_")


def probe_corporate_actions():
    print("\n=== Corporate Actions ===")
    save("corp_actions_list.json",
         {"path": "/market-data/instruments/stocks/corporate-actions/list"},
         api_get("/market-data/instruments/stocks/corporate-actions/list",
                 {"symbols": "AAPL", "page_size": "5"}))
    save("corp_actions_market.json",
         {"path": "/market-data/instruments/stocks/corporate-actions/market"},
         api_get("/market-data/instruments/stocks/corporate-actions/market",
                 {"market": "US", "page_size": "5"}))


def probe_fund_data():
    print("\n=== Fund Data ===")
    tests = [
        ("/market-data/fund/AAPL/nav", {"page_size": "5"}),
        ("/market-data/funds/AAPL/nav", {"page_size": "5"}),
        ("/market-data/etf/list", {"page_size": "5"}),
        ("/market-data/fund/list", {"page_size": "5"}),
    ]
    for path, q in tests:
        r = api_get(path, q)
        save(f"fund_{tag_from_path(path)}.json", {"path": path, "query": q}, r)
        if r.get("status") == 200:
            print(f"[FUND] {path} -> 200 OK")
            break


def probe_crypto_data():
    print("\n=== Crypto Data ===")
    tests = [
        ("/market-data/crypto/bars", {"ticker": "BTCUSD", "interval": "1m", "count": "5"}),
        ("/market-data/crypto/bars", {"ticker": "BTC", "type": "ETF", "interval": "1m", "count": "5"}),
        ("/market-data/crypto/BTCUSD/bars", {"interval": "1m", "count": "5"}),
    ]
    for path, q in tests:
        r = api_get(path, q, category=None)
        save(f"crypto_{tag_from_path(path)}.json", {"path": path, "query": q}, r)
        if r.get("status") == 200:
            print(f"[CRYPTO] {path} -> 200 OK")
            break


def probe_screener_v2():
    print("\n=== Screener v2 ===")
    body = {
        "filter": {"field": "price", "operator": "gt", "value": "100"},
        "sort": {"field": "price", "direction": "DESC"},
        "page_size": 5,
    }
    for path in ["/wlas/screener/ng/query", "/market-data/screener/ng/query"]:
        r = api_post(path, query={}, body=body)
        save(f"screener_v2_{tag_from_path(path)}.json", {"path": path, "body": body}, r)
        if r.get("status") == 200:
            print(f"[SCREENER] {path} -> 200 OK")
            break


def probe_broker_fd():
    print("\n=== Broker FD (broker-api host) ===")
    for path in [
        "/broker-fd/accounts",
        "/broker-fd/positions",
        "/broker-fd/assets",
        "/account/summary",
        "/account/positions",
    ]:
        r = broker_get(path, {})
        save(f"broker_api_{tag_from_path(path)}.json", {"path": path, "host": BROKER_API_URL}, r)
        if r.get("status") == 200:
            print(f"[BRK ] {path} -> 200 OK")


def main():
    print("=" * 60)
    print("Webull v0.5 Live Probe")
    print(f"Base: {BASE_URL}")
    print(f"Results: {RESULTS_DIR}")
    print("=" * 60)

    print("\n[INIT] Getting token via Go binary...")
    token = get_token()
    print(f"[INIT] Token: {token[:20]}...\n")

    probe_corporate_actions()
    probe_fund_data()
    probe_crypto_data()
    probe_screener_v2()
    probe_broker_fd()

    print("\n" + "=" * 60)
    print("Results:")
    for fn in sorted(os.listdir(RESULTS_DIR)):
        sz = os.path.getsize(os.path.join(RESULTS_DIR, fn))
        print(f"  {fn}  ({sz} bytes)")
    print("=" * 60)


if __name__ == "__main__":
    main()
