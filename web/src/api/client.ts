export interface FetchRecordPayload {
  token: string;
  dpop: string;
  cid: string;
  rev_handle: string;
  arweave_tx?: string;
}

const base = import.meta.env.VITE_GATEWAY_URL || "http://localhost:8080";

async function post<T>(path: string, body: unknown, dpop?: string): Promise<T> {
  const headers: Record<string, string> = { "Content-Type": "application/json" };
  if (dpop) headers["DPoP"] = dpop;
  const res = await fetch(`${base}${path}`, {
    method: "POST",
    headers,
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(`request failed: ${res.status} ${text}`);
  }
  return res.json();
}

export function validateCapability(token: string, dpop: string, rev_handle: string) {
  return post<{ valid: boolean; rev_state: string }>("/v1/validate-capability", {
    token,
    dpop,
    rev_handle,
  }, dpop);
}

export function fetchRecord(payload: FetchRecordPayload) {
  return post<{ ciphertext: string; cid: string }>("/v1/fetch-record", payload, payload.dpop);
}

export function revoke(rev_handle: string, reason: string, dpop?: string) {
  return post<{ tx_hash: string }>("/v1/revoke", { rev_handle, reason }, dpop);
}

export function breakGlassActivate(case_id: string, policy_ref: string, officer_signatures: string[], dpop?: string) {
  return post<{ capability: string; rev_handle: string; ttl_seconds: number }>(
    "/v1/breakglass/activate",
    { case_id, policy_ref, officer_signatures },
    dpop,
  );
}

export function resolveConflict(resourceId: string, chosen: "A" | "B", dpop?: string) {
  // Placeholder endpoint; adjust to actual reconciliation API when available.
  return post<{ ok: boolean }>("/v1/reconciliation/resolve", { resourceId, chosen }, dpop);
}
