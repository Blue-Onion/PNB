import { useState } from "react";
import type { Report } from "@/services/check-code";

const transactions: Record<string, { title: string; address: string }> = {
  approve: { title: "Token Approval", address: "Spender address" },
  transfer: { title: "Token Transfer", address: "Recipient address" },
  setApprovalForAll: { title: "NFT Operator Approval", address: "Operator address" },
};

function RiskBadge({ level, overall = false }: { level: string; overall?: boolean }) {
  const label = overall && ["HIGH", "MEDIUM", "LOW"].includes(level) ? level + " RISK" : level;
  return <span className="risk-badge" data-risk={level}>{label}</span>;
}

function Address({ label, value }: { label: string; value: string }) {
  const [copyStatus, setCopyStatus] = useState("");

  async function copyAddress() {
    try {
      await navigator.clipboard.writeText(value);
      setCopyStatus("Copied");
    } catch {
      setCopyStatus("Copy failed");
    }
  }

  return (
    <div className="address">
      <div className="eyebrow">{label}</div>
      <div className="address-row">
        <code>{value}</code>
        <button className="copy" onClick={copyAddress} type="button">Copy</button>
      </div>
      <span className="muted" role="status">{copyStatus}</span>
    </div>
  );
}

export function SecurityReport({ report }: { report: Report }) {
  const { transaction, ai_analysis: ai, analysis } = report;
  const details = transactions[transaction.function];
  const hasAI = analysis.status !== "AI_UNAVAILABLE";
  const findings = [...report.findings].sort((first, second) => {
    const levels = ["HIGH", "MEDIUM", "REVIEW", "LOW", "UNKNOWN", "SAFE"];
    return levels.indexOf(first.severity) - levels.indexOf(second.severity);
  });

  return (
    <article className="panel report" aria-labelledby="report-title">
      <header className="report-header">
        <div>
          <p className="eyebrow">txShield · Transaction security report</p>
          <h2 id="report-title">{details?.title ?? "Unknown Transaction Function"}</h2>
        </div>
        <RiskBadge level={report.final_risk} overall />
      </header>

      <section className="report-section">
        <h3 className="muted summary-heading">AI-generated summary</h3>
        <p>{hasAI && ai.summary ? ai.summary : "AI analysis is unavailable. Local findings are shown below."}</p>
        <dl className="metrics">
          {transaction.amount !== undefined && (
            <div className="metric">
              <dt>{transaction.function === "approve" ? "Approval amount" : "Transfer amount"}</dt>
              <dd>{BigInt(transaction.amount).toLocaleString("en-US")}</dd>
              <small>Raw token units</small>
            </div>
          )}
          {transaction.approved !== undefined && (
            <div className="metric">
              <dt>Operator permission</dt>
              <dd>{transaction.approved ? "Grant approval" : "Revoke approval"}</dd>
            </div>
          )}
          {hasAI && ai.confidence && (
            <div className="metric">
              <dt>AI confidence</dt>
              <dd className="confidence">{ai.confidence.toLowerCase()}</dd>
              <small>Self-reported</small>
            </div>
          )}
        </dl>
        {transaction.address && (
          <Address label={details?.address ?? "Address"} value={transaction.address} />
        )}
      </section>

      <section className="report-section">
        <h3>Security findings</h3>
        <ul className="findings">
          {findings.map((finding) => (
            <li key={finding.engine} className="finding">
              <div>
                <h4>{finding.engine === "SCAM_ADDRESS_CHECKER" ? "Local scam database match" : "Permission check"}</h4>
                <p className="muted">{finding.message}</p>
              </div>
              <RiskBadge level={finding.severity} />
            </li>
          ))}
        </ul>
      </section>

      {hasAI && ai.recommendations && ai.recommendations.length > 0 && (
        <section className="report-section">
          <h3>Recommendations</h3>
          <ul className="recommendations">
            {ai.recommendations.map((recommendation, index) => <li key={index}>{recommendation}</li>)}
          </ul>
        </section>
      )}

      <aside className="limitations">
        <h3>Analysis limitations</h3>
        <p className="muted">
          {analysis.reputation_source === "LOCAL_TEST_DATABASE" && "Address reputation uses a local test database. "}
          {!analysis.simulation_performed && "Transaction simulation was not performed."}
        </p>
      </aside>
    </article>
  );
}
