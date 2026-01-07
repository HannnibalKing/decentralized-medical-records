import React, { useState } from "react";
import { resolveConflict } from "../api/client";

const conflicts = [
  {
    id: "Obs-889",
    field: "Result",
    sourceA: "HgA1c 6.9% (Hospital A)",
    sourceB: "HgA1c 7.1% (Hospital B)",
    attA: "tx 0xabc",
    attB: "tx 0xdef",
  },
  {
    id: "Med-120",
    field: "Dosage",
    sourceA: "5mg daily",
    sourceB: "10mg daily",
    attA: "tx 0xaaa",
    attB: "tx 0xbbc",
  },
];

export function Reconciliation({ dpop }: { dpop: string }) {
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState<string | null>(null);

  const handleResolve = async (id: string, choice: "A" | "B") => {
    setLoading(id + choice);
    setMessage("");
    try {
      await resolveConflict(id, choice, dpop);
      setMessage(`Resolved ${id} with ${choice}`);
    } catch (err: any) {
      setMessage(err.message || "Resolve failed");
    } finally {
      setLoading(null);
    }
  };
  return (
    <div className="grid">
      <div className="panel">
        <h2>Conflicts</h2>
        <table className="table" aria-label="Conflicting entries">
          <thead>
            <tr>
              <th>Resource</th>
              <th>Field</th>
              <th>Source A</th>
              <th>Source B</th>
              <th>Resolve</th>
            </tr>
          </thead>
          <tbody>
            {conflicts.map((c) => (
              <tr key={c.id}>
                <td>{c.id}</td>
                <td>{c.field}</td>
                <td>
                  <div>{c.sourceA}</div>
                  <div className="badge">{c.attA}</div>
                </td>
                <td>
                  <div>{c.sourceB}</div>
                  <div className="badge">{c.attB}</div>
                </td>
                <td style={{ display: "flex", gap: 8 }}>
                  <button className="btn" onClick={() => handleResolve(c.id, "A")} disabled={loading === c.id + "A"}>
                    {loading === c.id + "A" ? "Resolving..." : "Accept A"}
                  </button>
                  <button className="btn ghost" onClick={() => handleResolve(c.id, "B")} disabled={loading === c.id + "B"}>
                    {loading === c.id + "B" ? "Resolving..." : "Accept B"}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="panel">
        <h2>Merge Summary</h2>
        <div className="grid two">
          <div className="card">
            <div className="kpi">2</div>
            <div style={{ color: "var(--muted)" }}>Conflicts remaining</div>
          </div>
          <div className="card">
            <div className="kpi">5</div>
            <div style={{ color: "var(--muted)" }}>Entries merged</div>
          </div>
        </div>
        <button className="btn" style={{ marginTop: 12 }}>Re-anchor merged bundle</button>
      </div>
      {message && <div className="badge" aria-live="polite">{message}</div>}
    </div>
  );
}
