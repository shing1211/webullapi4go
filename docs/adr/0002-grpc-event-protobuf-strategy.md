---
Status: Accepted
Date: 2026-09-18
Deciders: maintainer
---

# 0002. gRPC event protobuf strategy

## Context

Webull exposes event streams over gRPC for two groups of APIs:

- **Trading order-status events** (order lifecycle updates), relevant to v0.2/v0.3.
- **Broker account/instrument/corporate-actions/trade/funding/journal/master-data
  events**, relevant to v0.5.

These events are protobuf-encoded on the wire. Implementing them in a Go SDK
requires compiled Go types generated from the corresponding `.proto` schemas.

Facts established while planning v0.1:

- Webull **publishes** MQTT `Quote` / `Snapshot` / `Tick` protobuf definitions in
  its public documentation. Those are handled separately by the v0.1 MQTT work
  and are **not** the subject of this ADR.
- Webull **does not publish `.proto` files for the gRPC trade/broker event
  streams** in its public documentation, as of the date of this ADR.
- The official Python SDK, `github.com/webull-inc/webull-openapi-python-sdk`, is
  licensed **Apache-2.0** and bundles or generates the protobuf definitions used
  by the event streams.
- An official **Java SDK** also exists and ships generated descriptor/type
  artifacts.

**Inferred (not independently verified in this task):**

- The protobuf definitions inside the Python SDK are the same schemas Webull's
  servers use, and are complete enough to generate working Go stubs.
- The Java SDK's generated descriptors correspond to the same schemas.
- The Python SDK's bundled/generated `.proto` sources are themselves Apache-2.0
  (they are distributed inside an Apache-2.0 package), rather than carrying a
  separate license.

These inferences must be confirmed during the v0.3 spike before any vendoring is
merged (see Spike Plan).

## Problem

We want Go clients for Webull gRPC event streams, but the authoritative `.proto`
schemas are not published. We must choose how to obtain trustworthy message
definitions, or decide not to support these streams.

## Options Considered

### Option A — Vendor/derive `.proto` from the Apache-2.0 Python SDK

Extract the `.proto` sources (or derive them from generated Python modules) from
the official Python SDK, review them, and vendor them under `proto/` with a
`NOTICE`/attribution file. Generate Go types with `protoc` / `buf`.

**Pros**

- Uses the vendor's own definitions; highest chance of matching the wire format
  and field numbering.
- Apache-2.0 permits redistribution and derivative works with attribution.
- The schema set is versioned with the SDK, so upstream changes are trackable
  via dependency bumps/diffs.
- Works offline; no need to guess schemas from observed traffic.

**Cons**

- Requires a deliberate license/attribution process (NOTICE, retention of
  copyright notices, statement of modifications).
- The Python SDK may not ship clean `.proto` sources; definitions may be
  generated artifacts, requiring reverse-derivation.
- Bundled schemas may lag or lead the server; drift must be detected.
- Adds a build-time generation step and vendored third-party files to the repo.

### Option B — Observe/decode wire messages from the sandbox and hand-author `.proto`

Connect to the sandbox gRPC endpoints, capture encoded messages, decode them
with a generic protobuf reader, and hand-author `.proto` files.

**Pros**

- Schemas are validated against real sandbox traffic by construction.
- No third-party schema redistribution; authorship is ours.

**Cons**

- Reverse-engineering field numbers/types from wire bytes is error-prone.
- Unknown fields, optional fields, enums, and oneofs are hard to recover
  reliably; coverage depends on which events the sandbox happens to emit.
- Requires live sandbox access and credentials to (re)derive anything.
- Slow and fragile; likely to miss rarely-emitted message variants.

### Option C — Use the Java SDK's generated descriptors

Consume the Java SDK's generated descriptors (e.g. `FileDescriptorSet`) and use
them to generate Go types.

**Pros**

- Official artifacts, so schema fidelity should match the Java client.
- Generated descriptors are complete and machine-readable.

**Cons**

- Requires a Java toolchain (or prebuilt descriptor artifacts) in the build.
- Descriptor-to-Go generation is less direct than compiling `.proto` sources;
  may require custom extraction tooling.
- Same licensing/attribution review as Option A, plus the uncertainty of
  whether descriptor artifacts are redistributable and versioned conveniently.
- Heavier dependency footprint for a Go library.

### Option D — Do not support gRPC events until Webull publishes schemas; use the hosted Cloud MCP instead

Defer event-stream support and document the hosted Cloud MCP as the supported
path for event consumption.

**Pros**

- Zero licensing or reverse-engineering risk for our repo.
- No vendored third-party schemas to maintain.
- Focuses SDK scope on documented HTTP/MQTT surfaces.

**Cons**

- Leaves a documented Webull capability unsupported in the SDK.
- Hosted MCP is an external dependency and may not fit all users' deployment or
  latency requirements.
- "Until Webull publishes" has no committed date; this could mean "never".
- Contradicts the project roadmap, which lists gRPC events in v0.3 and v0.5.

## Decision

Adopt **Option A**: vendor/derive the gRPC event `.proto` files from the official
Apache-2.0 Python SDK, with full attribution, **contingent on the v0.3 spike
verifying schema fidelity against the sandbox**.

- Do **not** merge vendored protos until the spike confirms they round-trip with
  real sandbox events.
- If the spike shows the Python SDK definitions are unavailable, incomplete, or
  observably wrong against the sandbox, fall back to **Option B** as a
  stop-gap and continue monitoring for published schemas; **Option D** remains
  the fallback if neither is viable.

This decision is **Proposed** until the spike completes; it becomes **Accepted**
(or is superseded) at that point.

## Verification

Promoted to **Accepted** during the v0.3 spike (task T15.1). Findings:

- A real `.proto` source exists upstream at
  `webull/trade/events/events.proto`; it was pinned at commit
  `6d1418795449098404a9dbf9ac66e9c1dd9e9c47` (branch `main`).
- The upstream `LICENSE` is Apache-2.0 and the upstream `NOTICE`
  (`Webull OpenAPI Python SDK`, `Copyright 2022 Webull`) was located. The
  vendored `.proto` carries no per-file license header, so provenance rests on
  the repository-level `NOTICE`.
- The schema was vendored byte-identical (verified via matching git blob SHA-1)
  into `proto/webull/trade/events/v1/events.proto`, with no `go_package` edit;
  the Go package/import path is supplied by an `M` mapping in `buf.gen.yaml`.
- Go generation succeeded (buf v1.73.0, `protoc-gen-go` v1.36.1,
  `protoc-gen-go-grpc` v1.6.2) into `gen/webull/trade/events/v1/`, including the
  `EventService` server-streaming client. `go build`, `go vet`, `gofmt -l`, and
  `go test ./gen/...` pass offline.
- Attribution is recorded in `THIRD_PARTY_NOTICES.md`. The live-sandbox
  round-trip that this ADR gates on is verified separately by T15.3.

## Rationale

- Option A best balances schema fidelity, effort, and legal cleanliness:
  Apache-2.0 is explicitly designed to permit this kind of redistribution, and
  the definitions come from Webull itself.
- The v0.3 gate means we get the vendor's schemas while still requiring direct
  sandbox verification, which mitigates the main risk of trusting bundled
  artifacts.
- Options B and C are higher-effort or higher-uncertainty; D gives up a planned
  capability. They are retained as fallbacks rather than the primary path.

### Licensing analysis

- The Python SDK is Apache-2.0. Apache-2.0 is compatible with this project's
  Apache-2.0 license for redistribution of unmodified or modified files.
- Attribution requirements (Apache-2.0 section 4): retain copyright, patent,
  trademark, and attribution notices from the source; include a copy of the
  license; include a `NOTICE` file if the upstream distribution provides one;
  and state significant modifications to the vendored files.
- Trademarks: Apache-2.0 grants no trademark rights. "Webull" and related marks
  must not be relicensed or implied to be ours. Vendored files must keep
  upstream notices, and any derived files must clearly state they are derived
  from the Webull Python SDK, not authored by us.
- Generated Go code from vendored/derived protos is a derivative work; it must
  carry the same attribution headers or a repo-level NOTICE pointing to them.
- **To verify:** whether the bundled `.proto` sources carry a separate license
  header and whether the SDK ships an upstream `NOTICE` file.

## Consequences

- v0.3 adds a `proto/` vendor directory, a `NOTICE`/attribution file, and a
  protobuf generation step to the build.
- The repository will contain third-party-derived schema files; contributors
  must not edit them in place without recording modifications as required by
  Apache-2.0 section 4.
- Upstream Python SDK releases become a signal for schema drift; we should
  diff the vendored protos on SDK version bumps.
- Message and enum names, package paths, and field numbers are **deliberately
  unspecified here** and are to be determined during the v0.3 spike.

## Risks and Unknowns

- **Unknown:** whether the Python SDK ships `.proto` sources or only generated
  Python modules.
- **Unknown:** exact license headers on bundled protobuf artifacts, and presence
  of an upstream `NOTICE`.
- **Unknown:** whether sandbox gRPC endpoints and credentials are available to
  verify against (the v0.3 spike assumes they are; to be confirmed).
- **Risk:** schema drift between vendored definitions and the live server.
  Mitigation: verify against sandbox each time we bump the upstream SDK, and
  keep a generation/pin record.
- **Risk:** original protos may not be byte-for-byte complete; some fields may
  be server-populated extensions not present in the Python package.
- **Risk:** Apache-2.0 attribution being done incorrectly. Mitigation: review
  the `NOTICE` and generated-file headers before merge.
- **Risk:** if protos are hand-derived from generated Python, derivation errors
  are possible; mitigation: prefer shipping actual `.proto` sources when
  available.

## Spike Plan (v0.3)

Goal: prove that schemas vendored from the Apache-2.0 Python SDK generate Go
types that correctly decode real sandbox gRPC events, before merging any
vendored files.

1. **Inventory upstream.** Pin a specific Python SDK version/tag. Enumerate
   where protobuf definitions live (`.proto` sources vs. generated modules vs.
   packaged descriptors). Record findings in this ADR or a follow-up.
2. **License audit.** Locate the upstream `LICENSE` and `NOTICE`. Inspect
   per-file headers on the protobuf artifacts. Confirm Apache-2.0 applies and
   document attribution obligations.
3. **Extract.** Produce the candidate `.proto` set. If only generated Python is
   available, derive `.proto` with an appropriate tool and mark the derivation.
4. **Generate Go.** Run `protoc` (or `buf`) to generate Go types into an
   internal, non-public package.
5. **Verify against sandbox.** Connect to the sandbox gRPC event stream for at
   least one trading order-status flow. Decode real messages and assert that
   required fields populate sensibly.
6. **Cross-check (optional).** Compare against the Java SDK descriptors to
   confirm field numbers/types agree; treat disagreements as blockers until
   resolved.
7. **Decide.** If verification passes, promote this ADR from Proposed to
   Accepted and merge vendored protos with `NOTICE`. If it fails, open a new ADR
   selecting Option B or D and mark this one superseded.

Artifacts from the spike: the pinned upstream version, the extracted `.proto`
set, a `NOTICE` draft, and a short verification note.

## Notes

- This ADR intentionally contains no proto message names, package names, or
  field numbers. Those are **to be determined during the v0.3 spike** and must
  not be guessed.
- Anything labeled **Inferred** or **Unknown** is not a verified fact.
