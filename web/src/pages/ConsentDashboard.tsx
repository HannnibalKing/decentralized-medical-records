import React, { useState } from "react";
import { revoke, validateCapability } from "../api/client";

interface Delegate {
  name: string;
  role: string;
  scope: string;
  expires: string;
  status: "active" | "expiring";
  revHandle: string;
}

const delegates: Delegate[] = [
  { name: "Dr. Lee", role: "Clinician", scope: "Meds + Labs", expires: "Mar 03", status: "active", revHandle: "rev-handle-lee" },
  { name: "A. Rivera", role: "Caregiver", scope: "Vitals", expires: "Jan 17", status: "expiring", revHandle: "rev-handle-riv" },
];

export function ConsentDashboard({ token, dpop, revHandle }: { token: string; dpop: string; revHandle: string }) {
  const [loading, setLoading] = useState<string | null>(null);
  const [message, setMessage] = useState<string>("");
  const [validation, setValidation] = useState<string>("");

  const handleRevoke = async (h: string) => {
    try {
      setLoading(h);
      setMessage("");
      await revoke(h, "user_action", dpop);
      setMessage(`Revoked ${h}`);
    } catch (err: any) {
      setMessage(err.message || "Revoke failed");
    } finally {
      setLoading(null);
    }
  };

  const handleValidate = async () => {
    if (!token || !revHandle) {
      setValidation("Token and rev_handle required");
      return;
    }
    try {
      setValidation("Validating...");
      const res = await validateCapability(token, dpop, revHandle);
      setValidation(`${res.rev_state}`);
    } catch (err: any) {
      setValidation(err.message || "Validation failed");
    }
  };
  return (
    <div className="grid">
      <div className="panel">
        <h2>Active Delegates</h2>
        <div className="grid two">
          {delegates.map((d) => (
            <div className="card" key={d.name}>
              <div style={{ display: "flex", justifyContent: "space-between" }}>
                <div>
                  <div style={{ fontWeight: 700 }}>{d.name}</div>
                  <div style={{ color: "var(--muted)", fontSize: 13 }}>{d.role}</div>
                </div>
                <span className={`badge ${d.status === "expiring" ? "warn" : "success"}`}>
                  {d.status === "expiring" ? "Expiring soon" : "Active"}
                </span>
              </div>
              <div style={{ color: "var(--muted)", fontSize: 13 }}>Scope: {d.scope}</div>
              <div style={{ color: "var(--muted)", fontSize: 13 }}>Expires: {d.expires}</div>
              <div style={{ display: "flex", gap: 8 }}>
                <button
                  className="btn"
                  aria-label={`Revoke ${d.name}`}
                  onClick={() => handleRevoke(d.revHandle || revHandle)}
                  disabled={loading === d.revHandle}
                >
                  {loading === d.revHandle ? "Revoking..." : "Revoke"}
                </button>
                <button className="btn ghost" aria-label={`Extend ${d.name}`}>Extend</button>
              </div>
            </div>
          ))}
        </div>
      </div>

      <div className="panel">
        <h2>Tokens & Revocation</h2>
        <div className="grid two">
          <div className="card">
            <div className="kpi">12</div>
            <div style={{ color: "var(--muted)" }}>Active capability tokens</div>
            <button className="btn ghost">Revoke all</button>
          </div>
          <div className="card">
            <div className="kpi">4</div>
            <div style={{ color: "var(--muted)" }}>Expiring in next 24h</div>
            <button className="btn ghost">Notify holders</button>
          </div>
        </div>
        <div style={{ marginTop: 10, display: "flex", gap: 8, alignItems: "center" }}>
          <button className="btn ghost" onClick={handleValidate}>Validate token</button>
          {validation && <span className="badge" aria-live="polite">{validation}</span>}
        </div>
      </div>

      <div className="panel">
        <h2>Audit Highlights</h2>
        <table className="table" aria-label="Recent audit entries">
          <thead>
            <tr>
              <th>Actor</th>
              <th>Action</th>
              <th>Resource</th>
              <th>Time</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>Dr. Lee</td>
              <td>Read</td>
              <td>MedicationRequest/123</td>
              <td>10:22</td>
            </tr>
            <tr>
              <td>A. Rivera</td>
              <td>Read</td>
              <td>Observation/445</td>
              <td>09:58</td>
            </tr>
            <tr>
              <td>Gateway</td>
              <td>Revoke</td>
              <td>cap-handle-xyz</td>
              <td>09:10</td>
            </tr>
          </tbody>
        </table>
      </div>
      {message && <div className="badge" aria-live="polite">{message}</div>}
    </div>
  );
}
