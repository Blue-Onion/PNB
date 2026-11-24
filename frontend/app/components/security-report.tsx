import { useState } from "react";
import { RISK_LEVELS, type Report } from "@/services/check-code";

function Address({ value }: { value?: string }) {
  const [copyStatus, setCopyStatus] = useState("");

  async function copyAddress() {
    if (!value) return;
    try {
      await navigator.clipboard.writeText(value);
      setCopyStatus("Copied");
    } catch {
      setCopyStatus("Copy failed");
    }
  }

  return (
    <>
      <code>{value ?? "Unavailable"}</code>
      {value && (
        <button className="copy" onClick={copyAddress} type="button">
          {copyStatus || "Copy"}
        </button>
      )}
      <span className="sr-only" role="status">{copyStatus}</span>
    </>
  );
}

export function SecurityReport({ report, raw }: { report: Report; raw: string }) {
  const { transaction, ai_analysis: ai, analysis, final_risk: risk } = report;
  const permission = report.findings.find((finding) => finding.engine === "PERMISSION_CHECKER");
  const reputation = report.findings.find((finding) => finding.engine === "SCAM_ADDRESS_CHECKER");
  const findings = [...report.findings].sort((first, second) =>
    RISK_LEVELS.indexOf(second.severity) - RISK_LEVELS.indexOf(first.severity),
  );
  const hasAI = analysis.status === "COMPLETED";
  const proceed = risk === "SAFE" || risk === "LOW";
  const amount = transaction.amount;
  const permissionText = transaction.approved === undefined
    ? amount ?? "Unavailable"
    : transaction.approved ? "Grant collection-wide access" : "Revoke operator access";
  const engines = [
    {
      name: "Hex Decoder",
      status: transaction.function === "unknown" ? "Unknown selector" : "Decoded",
    },
    {
      name: "Permission Checker",
      status: permission?.severity ?? "No finding",
    },
    {
      name: "Scam Address Checker",
      status: reputation ? "Local test match" : transaction.address ? "No local test match" : "No address to check",
    },
    {
      name: "Transaction Simulator",
      status: analysis.simulation_performed ? "Performed" : "Skipped",
    },
    {
      name: "Risk Scoring + AI Explanation",
      status: analysis.status,
    },
  ];

  return (
    <article className="report" aria-labelledby="report-title">
      <header className="verdict panel" data-risk={risk}>
        <span className="risk-badge">{risk}</span>
        <h2 id="report-title">{proceed ? "PROCEED" : "REJECT TRANSACTION"}</h2>
        <p>{findings[0]?.message ?? (hasAI && ai.summary ? ai.summary : "No warning was returned.")}</p>
      </header>

      <section className="panel report-section" aria-labelledby="breakdown-title">
        <h3 id="breakdown-title">Decoded breakdown</h3>
        <dl className="metrics">
          <div className="metric">
            <dt>Function</dt>
            <dd><code>{transaction.function}</code></dd>
          </div>
          <div className="metric">
            <dt>Spender / Target Address</dt>
            <dd><Address value={transaction.address} /></dd>
          </div>
          <div className="metric">
            <dt>Amount / Permission</dt>
            <dd>{permissionText}</dd>
            {amount !== undefined && <small>Raw token units</small>}
          </div>
        </dl>
      </section>

      <section className="panel report-section" aria-labelledby="pipeline-title">
        <h3 id="pipeline-title">Engine pipeline</h3>
        <ol className="engines">
          {engines.map((engine, index) => (
            <li key={engine.name}>
              <span className="engine-number">Engine {index + 1}</span>
              <strong>{engine.name}</strong>
              <span className="engine-status">{engine.status}</span>
            </li>
          ))}
        </ol>
      </section>

      <section className="panel report-section" aria-labelledby="findings-title">
        <h3 id="findings-title">Security findings</h3>
        <ul className="findings">
          {findings.map((finding) => (
            <li key={finding.engine}>
              <span className="risk-badge">{finding.severity}</span>
              <span>{finding.message}</span>
            </li>
          ))}
        </ul>
        <h4>AI explanation</h4>
        <p className="muted">{hasAI ? ai.summary : "AI analysis is unavailable. Local findings are shown above."}</p>
        {hasAI && !!ai.recommendations?.length && (
          <ul className="recommendations">
            {ai.recommendations.map((recommendation, index) => <li key={index}>{recommendation}</li>)}
          </ul>
        )}
      </section>

      <details className="panel raw-response">
        <summary>Raw API response</summary>
        <pre tabIndex={0} aria-label="Raw response JSON">{raw}</pre>
      </details>
    </article>
  );
}
