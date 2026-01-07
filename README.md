# Decentralized Medical Records Platform

Production-focused scaffold with:
- Smart contracts: `RevocationRegistry`, `AttestationRegistry`, `MerkleAnchor` (Foundry).
- Gateway (Go): revocation-aware handlers, IPFS/Arweave clients, REST + gRPC proto.
- Web UX (React/Vite): consent dashboard, break-glass, reconciliation, attestations.
- Docs: whitepaper, API specs, diagrams, policies, runbook. MIT licensed.

## Quick Start

### Contracts
```powershell
cd "c:\Users\HellKnight\Documents\GitHub\UX\Decentralized Medical Records Platform"
forge install foundry-rs/forge-std
forge test
```

### Mock revocation RPC (dev only)
```powershell
cd "c:\Users\HellKnight\Documents\GitHub\UX\Decentralized Medical Records Platform\gateway"
go run .\cmd\mockrev
```

### Gateway (HTTP :8080)
```powershell
cd "c:\Users\HellKnight\Documents\GitHub\UX\Decentralized Medical Records Platform\gateway"
$env:GATEWAY_ADDR=":8080"
$env:REVOCATION_RPC="http://localhost:8181"
$env:IPFS_URL="http://localhost:5001"
$env:ARWEAVE_URL="https://arweave.net"
go run .\cmd\gateway
```

### Web UX (http://localhost:5173)
```powershell
cd "c:\Users\HellKnight\Documents\GitHub\UX\Decentralized Medical Records Platform\web"
if (Test-Path node_modules) { Remove-Item node_modules -Recurse -Force }
if (Test-Path package-lock.json) { Remove-Item package-lock.json -Force }
npm install
set VITE_GATEWAY_URL=http://localhost:8080
npm run dev
```

### Diagrams
Open `docs/diagrams/*.puml` in any PlantUML renderer.

## Docs
- Whitepaper: `docs/whitepaper.md`
- API: `docs/api/gateway.proto`, `docs/api/http-openapi.yaml`
- Policies: `docs/policies/capability-claims.md`
- Operations: `docs/runbook/operations.md`

## License
MIT — see `LICENSE`.
