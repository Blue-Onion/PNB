"use client";

import { useState, type FormEvent } from "react";
import { useFetch } from "@/hooks/use-fetch";
import { checkCode } from "@/services/check-code";
import { SecurityReport } from "./components/security-report";

export default function Home() {
  const [hexcode, setHexcode] = useState("");
  const { run, result, error, pending, dismissError } = useFetch(checkCode);

  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!pending) void run(hexcode.trim());
  }

  return (
    <main className="workspace">
      <header className="site-header">
        <span className="wordmark">txShield</span>
        <span className="header-label">Ethereum transaction analyzer</span>
      </header>

      <section className="panel input-panel" aria-labelledby="page-title">
        <h1 id="page-title">Review a transaction</h1>
        <p className="muted">Paste Ethereum calldata to inspect permissions and security findings before signing.</p>
        <form onSubmit={submit}>
          <label htmlFor="hexcode">Transaction calldata (hex)</label>
          <textarea
            id="hexcode"
            placeholder="0x…"
            value={hexcode}
            onChange={(event) => setHexcode(event.target.value)}
            rows={3}
            autoCapitalize="off"
            autoComplete="off"
            spellCheck={false}
            required
          />
          <button className="submit" type="submit" disabled={pending || !hexcode.trim()}>
            {pending ? "Analyzing…" : "Analyze calldata"}
          </button>
        </form>
        <div className="sr-only" role="status">
          {pending ? "Analysis in progress" : result ? "Analysis complete. The report is below." : ""}
        </div>
      </section>

      {result && <SecurityReport report={result} />}
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
