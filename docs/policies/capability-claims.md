# Capability Claims (Gateway Validation Rules)

## Required claims
- iss: issuer DID; must match trusted issuer set.
- sub: subject DID (actor using token).
- pat: patient DID.
- aud: must equal gateway audience.
- scope: object with fhirTypes[], optional fields mask, actions[], breakGlass? boolean.
- purpose: e.g., treatment, operations, emergency.
- nbf / exp: short TTL (minutes-hours); clock skew tolerance <= 2m.
- rev: revocation handle (bytes32 base58/hex) resolvable on-chain.
- policy: policy identifier with version (e.g., policy-ER-v3).
- delegationDepth: 0 unless caregiver delegation is permitted.
- nonce: random per token issuance.
- cnf: DPoP key thumbprint (binds to TLS key).
- Optional: caseId for break-glass.

## Validation logic (gateway)
1) Verify signature against iss DID keys (assertion method); ensure key not revoked/rotated out.
2) Check nbf <= now + skew; exp > now; enforce max TTL.
3) Match aud; verify purpose is allowed for requested endpoint.
4) Check revocation handle on-chain; if revoked -> deny.
5) Enforce policy version; if outdated -> deny/refresh.
6) Enforce delegationDepth (no re-delegation unless allowed).
7) Verify cnf/DPoP: proof key thumbprint matches cnf; nonce fresh.
8) For breakGlass=true: confirm chain event and policyRef, enforce single-use TTL.
9) Field mask (if present): restrict returned fields to mask.

## Error codes (HTTP/gRPC)
- invalid_signature
- token_expired
- token_not_yet_valid
- audience_mismatch
- purpose_denied
- revoked
- policy_version_mismatch
- delegation_not_allowed
- dpop_invalid
- breakglass_policy_violation
- scope_denied

## Security controls
- Short TTL for all capabilities; break-glass <= 30m and single-use.
- All responses bound to DPoP key; no bearer semantics.
- Audit every decision with revStateAtUse.
- Enforce rate limits per DID and IP.
