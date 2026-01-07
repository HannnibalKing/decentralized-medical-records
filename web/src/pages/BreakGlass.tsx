import React, { useState } from "react";
import { breakGlassActivate } from "../api/client";

const emergencies = [
  { caseId: "ER-4421", scope: "Emergency bundle", ttl: "24m", status: "active" },
  { caseId: "ER-4310", scope: "Vitals only", ttl: "expired", status: "revoked" },
];

export function BreakGlass({ token, dpop }: { token: string; dpop: string }) {
  const [caseId, setCaseId] = useState("ER-1234");
  const [scope, setScope] = useState("Emergency bundle");
  const [ttl, setTtl] = useState("30");
  const [sig, setSig] = useState("");
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);

  const onActivate = async () => {
    setLoading(true);
    setMessage("");
    try {
      const res = await breakGlassActivate(caseId, scope, sig ? [sig] : [], dpop);
      setMessage(`Issued rev_handle ${res.rev_handle} TTL ${res.ttl_seconds}s`);
    } catch (err: any) {
      setMessage(err.message || "Activation failed");
    } finally {
      setLoading(false);
    }
  };
  return (
    <div className="grid">
      <div className="panel">
        <h2>Activate Break-Glass</h2>
        <div className="grid two">
          <div className="card">
            <div style={{ fontWeight: 700 }}>Case ID</div>
            <input aria-label="Case ID" placeholder="ER-1234" style={inputStyle} value={caseId} onChange={(e) => setCaseId(e.target.value)} />
            <div style={{ fontWeight: 700 }}>Scope</div>
            <select aria-label="Scope" style={inputStyle} value={scope} onChange={(e) => setScope(e.target.value)}>
              <option>Emergency bundle</option>
              <option>Vitals only</option>
              <option>Meds + allergies</option>
            </select>
            <div style={{ fontWeight: 700 }}>TTL (minutes)</div>
            <input aria-label="TTL" placeholder="30" style={inputStyle} value={ttl} onChange={(e) => setTtl(e.target.value)} />
            <button className="btn" onClick={onActivate} disabled={loading}>
              {loading ? "Activating..." : "Activate (requires 2-of-3)"}
            </button>
          </div>
          <div className="card">
            <div style={{ fontWeight: 700 }}>Officer Signatures</div>
            <div className="badge warn">Pending multisig</div>
            <p style={{ color: "var(--muted)" }}>Paste aggregated signature payload.</p>
            <input aria-label="Officer signatures" placeholder="sigpayload" style={inputStyle} value={sig} onChange={(e) => setSig(e.target.value)} />
          </div>
        </div>
        {message && <div className="badge" aria-live="polite">{message}</div>}
      </div>

      <div className="panel">
        <h2>Recent Emergency Access</h2>
        <table className="table" aria-label="Break-glass events">
          <thead>
            <tr>
              <th>Case</th>
              <th>Scope</th>
              <th>TTL</th>
              <th>Status</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            {emergencies.map((e) => (
              <tr key={e.caseId}>
                <td>{e.caseId}</td>
                <td>{e.scope}</td>
                <td>{e.ttl}</td>
                <td>
                  <span className={`badge ${e.status === "active" ? "success" : "danger"}`}>
                    {e.status}
                  </span>
                </td>
                <td>
                  <button className="btn ghost">Revoke</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

const inputStyle: React.CSSProperties = {
  background: "#0b1220",
  border: "1px solid var(--border)",
  color: "var(--text)",
  padding: "10px 12px",
  borderRadius: 10,
  outline: "none",
};
