# Operations Runbook (Base L2, IPFS + Arweave)

## Key rotation
- Gateway TLS + signing keys: rotate weekly; publish new JWKs; overlap 48h; enforce kid.
- Org HSM keys: quarterly rotation; update DID docs; deprecate old assertion keys after overlap.
- Contract ownership keys: store in multisig; no hot wallets.

## Revocation SLA
- Standard revoke: propagate on-chain within seconds; gateway caches <= 30s with negative lookup refresh.
- Break-glass: auto-revoke immediately after first use; TTL <= 30m.
- Monitoring: alert if revocation check latency > 500ms or cache staleness > 60s.

## Backup/restore
- IPFS: pin in >=3 regions; monitor pinset divergence; weekly GC disabled on pin nodes.
- Arweave: mirror all ciphertext uploads; verify receipt IDs; monthly restore drill from Arweave-only.
- Chain data: rely on L2; archive node or third-party provider; checkpoint contract events to off-chain store.

## Monitoring & alerts
- Revocation registry lag (expected <2 blocks).
- Anchor confirmation lag (expected <2 blocks on Base).
- Signature verification failures rate.
- Break-glass activations and uses (page on activation, P1 on use).
- Mass-export detection (>N records/min per actor).
- IPFS fetch failures; Arweave fallback ratio.

## Incident response
- Suspected key compromise: freeze issuer DID keys; rotate; revoke affected capability handles; invalidate JWKs.
- Data poisoning: flag bad anchors; mark affected roots as invalid; require re-anchor with correct data.
- Storage DoS: throttle clients; switch to read-only mode; favor Arweave pulls.

## Change management
- All contract changes via multisig; dry-run on testnet; require audit before mainnet.
- Gateway deploys: blue/green; feature-flag new validation rules; canary 5% traffic.

## DR scenarios
- L2 outage: continue to serve reads with last confirmed anchors; queue writes; once chain resumes, reconcile queue.
- IPFS outage: serve from Arweave; re-pin when IPFS returns.
- Cache corruption: fail closed; require fresh chain reads; restart gateway nodes with empty cache.
