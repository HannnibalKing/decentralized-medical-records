# Decentralized Medical Records Platform (DMRP)

Version: 0.1 (Jan 7, 2026)

## Abstract
DMRP is a patient-centric health data platform for life-critical scenarios, combining off-chain encrypted FHIR records, on-chain attestations, and granular capability tokens. Data remains patient-controlled with zero-trust storage (IPFS/Arweave), revocable granular access, emergency break-glass with auditability, multi-party prescription attestation, time-locked delegation, and accessibility-first UX.

## Requirements Summary
- Life-critical data: high availability, integrity, and verifiable provenance.
- Zero-knowledge privacy posture: ciphertext-only storage, minimal metadata leakage.
- Granular permissions: capability tokens with scopes (resource types, fields), purpose-of-use, TTL, revocation handles.
- Emergency access: break-glass by multi-party approval, auto-expire, full audit trail, patient notification.
- Multi-party prescription attestation: prescriber, pharmacist, insurer.
- Time-locked revocation and caregiver delegation with expiry.
- Cross-hospital reconciliation UI.
- Accessibility for elderly/disabled.

## Architecture Overview
- Storage: Encrypted FHIR bundles on IPFS; Arweave immutable backup. Only ciphertext is stored.
- Chain: Base L2 (EVM). Contracts: `RevocationRegistry`, `AttestationRegistry`, `MerkleAnchor`.
- Gateway: Validation, revocation checks, storage fetch, audit emission, HPKE unwrap (planned).
- Identity: DIDs (did:ion/did:key). Biometric-backed passkeys for patients/caregivers; HSM keys for organizations.
- Data model: Standard FHIR resources; field-level encryption possible for sensitive data.
- UX: Web client (React) for consent dashboard, break-glass, reconciliation, attestations.

Key repo references:
- Contracts: `contracts/RevocationRegistry.sol`, `contracts/AttestationRegistry.sol`, `contracts/MerkleAnchor.sol`.
- Gateway API: `docs/api/gateway.proto`, `docs/api/http-openapi.yaml`.
- Diagrams: `docs/diagrams/*.puml`.

## Cryptography & Keys
- Per-resource symmetric encryption: AES-256-GCM with per-object nonce; AAD binds `{CID, resourceType, patientDID, version, timestamp}`.
- Key distribution: HPKE (X25519/HKDF-SHA256/AES-256-GCM) to wrap per-resource keys to recipient DIDs.
- Identity keys: Assertion, authentication, agreement keys in DID Docs; rotation supported.
- Audit integrity: Hash-chained records; periodic Merkle root anchored on-chain.

### Capability Token Model
- Claims: `iss, sub, pat, aud, scope{fhirTypes, fields?, actions}, purpose, nbf, exp, rev, policy, delegationDepth, nonce, cnf(jkt), breakGlass?, caseId?`.
- Format: JWT/CWT with signature (issuer assertion key). DPoP-style confirmation binds to TLS key.
- Revocation: On-chain registry keyed by `rev` handle (bytes32). Short TTL enforced to reduce window.

## Core Protocol Flows
See detailed sequence diagrams in `docs/diagrams`.

1. Standard Read
- Client presents capability + DPoP.
- Gateway validates signature, TTL, audience, purpose; checks revocation; fetches CID from IPFS (Arweave fallback); verifies anchor; unwraps key (planned) and returns ciphertext.

2. Write (Clinician)
- Create FHIR resource, encrypt, store to IPFS; anchor hash+metadata via `MerkleAnchor`; issue capabilities; audit event.

3. Revocation
- Mark handle in `RevocationRegistry`; clients recheck per call; short TTL shrinks stale windows.

4. Caregiver Delegation
- Patient step-up (biometric) to issue scope-limited, expiring token; revocable; audit on use.

5. Break-Glass
- Patient policy predefines institutions/scope/TTL.
- Activation via N-of-M emergency officer approvals; gateway issues single-use, short TTL capability; auto-revokes after use; real-time patient notifications and audit.

6. Prescription Attestation
- Prescriber signs `MedicationRequest`, pharmacist signs `MedicationDispense`, insurer optionally signs coverage; `AttestationRegistry` stores status with aggregate signature hash; clients verify before fulfill.

## Smart Contracts (EVM - Base L2)
- `RevocationRegistry`: `revoke(bytes32 handle, string reason)`, `isRevoked(bytes32)`, `get(bytes32)`.
- `AttestationRegistry`: `record(bytes32 cidMr, bytes32 cidMd, bytes32 aggSig, Status status)`, `updateStatus(bytes32 cidMr, Status)`.
- `MerkleAnchor`: `anchor(bytes32 root, bytes32 cid, string resourceType, bytes32 patient)`; emits `Anchored` events.

Security considerations:
- Idempotency, event emissions for indexers, minimized storage writes, no arbitrary reentrancy.

## Gateway (Validation & Storage)
- Revocation checks with cache; pluggable revocation client (HTTP mock in `gateway/cmd/mockrev`).
- DPoP header enforced (presence; full crypto binding planned).
- IPFS via Kubo HTTP API; Arweave via public gateway fallback.
- gRPC service stub (build tag `grpc`) aligns with `gateway.proto`.

Endpoints (REST):
- `POST /v1/validate-capability` → revocation-aware validation.
- `POST /v1/fetch-record` → returns ciphertext + AAD + wrappedKey (planned).
- `POST /v1/revoke` → posts revocation.
- `POST /v1/breakglass/activate` → issues single-use capability (planned).
- `GET /v1/attestations/{cid_mr}` → retrieve attestation (planned).

## Data Integrity
- Each bundle hashed; Merkle root anchored via `MerkleAnchor`. Clients verify CID + Merkle proof before decrypting.
- All audit logs hash-chained with periodic root anchoring.

## Threat Model & Mitigations
- Replay: Short TTL, nonces, DPoP token binding, server challenges.
- MITM: mTLS and certificate pinning; DPoP required.
- Key theft: HSM/TEE for custodial keys; device-bound passkeys; biometric step-up.
- Data poisoning: Signature and policy checks on writes; provenance displayed; reconciliation workflow requires explicit accept/merge.
- Storage DoS/ransom: Dual-backend (IPFS + Arweave), regional pinning, backpressure.
- Privacy: On-device biometrics; minimal metadata; field-level encryption for sensitive fields.

## Accessibility & UX
- WCAG 2.1 AA baseline: high contrast, large touch targets, text scaling, ARIA labels, voice commands for essential actions.
- Consent Dashboard: list delegates, scopes, expiries, audit highlights, one-tap revoke.
- Break-Glass: clear activation UI, officer signatures collection, banners, notifications.
- Reconciliation: timeline + diff view, provenance chips, accept/merge with re-anchor.

## Compliance
- SMART-on-FHIR compatibility for EHR integration; OAuth2 edge translation to capabilities.
- HIPAA/GDPR: purpose-of-use tagging, revocation registries, immutable audits, data minimization.
- Biometric privacy: on-device matching, liveness checks, non-extractable templates.

## Operations
- Key rotation: weekly gateway keys, quarterly org keys; overlapping validity windows.
- Revocation SLA: seconds-level propagation; break-glass immediate auto-revoke.
- Backups: IPFS pins ≥3 regions; Arweave mirror; monthly restore drills.
- Monitoring: revocation/anchor lags, signature failures, break-glass events, mass export detection.

## Performance & Availability
- Gateway stateless design with revocation cache; horizontal scaling.
- Anchoring batched by epoch; IPFS pinning managed asynchronously.
- Read path designed to tolerate L2 delays (serve last confirmed anchors, queue writes).

## Roadmap
- HPKE wrap/unwrap end-to-end in gateway.
- Full DPoP crypto validation and nonce service.
- ZK proofs for selective field-disclosure (optional future work).
- Attestation registry indexer and verification UI.

## License
This project is released under the MIT License. See `LICENSE`.
