import React from "react";

type NavKey = "consent" | "breakglass" | "reconciliation" | "attestations";

interface LayoutProps {
  current: NavKey;
  onNav: (k: NavKey) => void;
  children: React.ReactNode;
}

const navItems: { key: NavKey; label: string }[] = [
  { key: "consent", label: "Consent Dashboard" },
  { key: "breakglass", label: "Break-Glass" },
  { key: "reconciliation", label: "Reconciliation" },
  { key: "attestations", label: "Prescriptions" },
];

export function Layout({ current, onNav, children }: LayoutProps) {
  return (
    <div className="app-shell">
      <nav className="nav" aria-label="Main">
        <h1>DMR Control</h1>
        {navItems.map((item) => (
          <button
            key={item.key}
            className={current === item.key ? "active" : ""}
            onClick={() => onNav(item.key)}
            aria-current={current === item.key}
          >
            {item.label}
          </button>
        ))}
      </nav>
      <main className="content">{children}</main>
    </div>
  );
}
