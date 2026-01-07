import React, { useState } from "react";
import { Layout } from "./components/Layout";
import { ConsentDashboard } from "./pages/ConsentDashboard";
import { BreakGlass } from "./pages/BreakGlass";
import { Reconciliation } from "./pages/Reconciliation";
import { Attestations } from "./pages/Attestations";

type NavKey = "consent" | "breakglass" | "reconciliation" | "attestations";

export default function App() {
  const [nav, setNav] = useState<NavKey>("consent");
  const [token, setToken] = useState("");
  const [dpop, setDpop] = useState("");
  const [revHandle, setRevHandle] = useState("");

  return (
    <Layout current={nav} onNav={setNav}>
      <div className="panel" style={{ display: "flex", gap: 12, flexWrap: "wrap" }}>
        <div style={{ display: "grid", gap: 6, minWidth: 220 }}>
          <label style={{ color: "var(--muted)" }}>Capability Token</label>
          <input aria-label="Capability token" style={inputStyle} value={token} onChange={(e) => setToken(e.target.value)} />
        </div>
        <div style={{ display: "grid", gap: 6, minWidth: 220 }}>
          <label style={{ color: "var(--muted)" }}>DPoP Proof</label>
          <input aria-label="DPoP" style={inputStyle} value={dpop} onChange={(e) => setDpop(e.target.value)} />
        </div>
        <div style={{ display: "grid", gap: 6, minWidth: 220 }}>
          <label style={{ color: "var(--muted)" }}>Revocation Handle</label>
          <input aria-label="Revocation handle" style={inputStyle} value={revHandle} onChange={(e) => setRevHandle(e.target.value)} />
        </div>
      </div>

      {nav === "consent" && <ConsentDashboard token={token} dpop={dpop} revHandle={revHandle} />}
      {nav === "breakglass" && <BreakGlass token={token} dpop={dpop} />}
      {nav === "reconciliation" && <Reconciliation dpop={dpop} />}
      {nav === "attestations" && <Attestations />}
    </Layout>
  );
}

const inputStyle: React.CSSProperties = {
  background: "#0b1220",
  border: "1px solid var(--border)",
  color: "var(--text)",
  padding: "8px 10px",
  borderRadius: 10,
  outline: "none",
};
