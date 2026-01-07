import React from "react";

const rows = [
  { cid: "CID_MR_123", status: "Issued", agg: "aggSig123", ts: "10:20" },
  { cid: "CID_MR_456", status: "Dispensed", agg: "aggSig456", ts: "09:55" },
];

export function Attestations() {
  return (
    <div className="grid">
      <div className="panel">
        <h2>Prescription Attestations</h2>
        <table className="table" aria-label="Prescription attestations">
          <thead>
            <tr>
              <th>MedicationRequest CID</th>
              <th>Status</th>
              <th>Aggregate Sig</th>
              <th>Time</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((r) => (
              <tr key={r.cid}>
                <td>{r.cid}</td>
                <td>
                  <span className="badge success">{r.status}</span>
                </td>
                <td className="badge">{r.agg}</td>
                <td>{r.ts}</td>
                <td>
                  <button className="btn ghost">Verify</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
