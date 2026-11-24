"use client";

import { useState, type FormEvent } from "react";
import { useFetch } from "@/hooks/use-fetch";
import { checkCode } from "@/services/check-code";
import { demoPresets } from "@/services/demo-presets";
import { SecurityReport } from "./components/security-report";

function ReportSkeleton() {
  return (
    <div className="report report-skeleton" aria-hidden="true">
      <div className="panel">
        <div className="skeleton-bar skeleton-badge" />
        <div className="skeleton-bar skeleton-title" />
        <div className="skeleton-bar skeleton-summary" />
      </div>
      <div className="panel">
        <div className="skeleton-bar skeleton-heading" />
        <div className="metrics">
          {["function", "address", "permission"].map((field) => (
            <div className="metric" key={field}>
              <div className="skeleton-bar skeleton-label" />
              <div className="skeleton-bar skeleton-value" />
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

export default function Home() {
  const [hexcode, setHexcode] = useState("");
  const { run, result, error, pending, reset, dismissError } = useFetch(checkCode);

  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!pending) void run(hexcode.trim());
  }

  return (
    <main className="workspace">
      <header className="site-header">
        <span className="wordmark"><span aria-hidden="true" className="brand-mark">◇</span>txShield</span>
        <span className="header-label">Ethereum transaction analyzer</span>
      </header>

      <section className="panel input-panel" aria-labelledby="page-title">
        <div aria-hidden="true" className="input-icon">&gt;_</div>
        <h1 id="page-title">Review a transaction</h1>
        <p className="muted">Paste Ethereum calldata to inspect permissions and security findings before signing.</p>
        <form onSubmit={submit}>
          <label htmlFor="hexcode">Transaction calldata (hex)</label>
          <textarea
            id="hexcode"
            placeholder="0x…"
            value={hexcode}
            onChange={(event) => {
              setHexcode(event.target.value);
              reset();
            }}
            disabled={pending}
            rows={3}
            autoCapitalize="off"
            autoComplete="off"
            spellCheck={false}
            required
          />
          <fieldset className="presets" disabled={pending}>
            <legend>Demo presets</legend>
            <div className="preset-list">
              {demoPresets.map((preset) => (
                <button
                  key={preset.label}
                  type="button"
                  aria-pressed={hexcode === preset.calldata}
                  onClick={() => {
                    setHexcode(preset.calldata);
                    reset();
                  }}
                >
                  {preset.label}
                </button>
              ))}
            </div>
          </fieldset>
          <button className="submit" type="submit" disabled={pending}>
            {pending ? "Analyzing…" : "Analyze calldata"}
            <span aria-hidden="true">{pending ? "…" : "→"}</span>
          </button>
        </form>
        <div className="sr-only" role="status">
          {pending ? "Analysis in progress" : result ? "Analysis complete. The report is below." : ""}
        </div>
      </section>

      {pending ? <ReportSkeleton /> : result && <SecurityReport key={result.report.report_id} {...result} />}
      {error && (
        <div className="toast">
          <div role="alert">
            <strong>Analysis failed</strong>
            <p>{error}</p>
          </div>
          <button type="button" aria-label="Dismiss error" onClick={dismissError}>×</button>
        </div>
      )}
    </main>
  );
}
