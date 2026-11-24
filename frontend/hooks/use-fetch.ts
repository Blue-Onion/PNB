import { useState } from "react";

export function useFetch<Input, Result>(request: (input: Input) => Promise<Result>) {
  const [result, setResult] = useState<Result | null>(null);
  const [error, setError] = useState("");
  const [pending, setPending] = useState(false);

  async function run(input: Input) {
    setPending(true);
    setError("");
    setResult(null);
    try {
      setResult(await request(input));
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "The request failed. Try again.");
    } finally {
      setPending(false);
    }
  }

  function reset() {
    setResult(null);
    setError("");
  }

  return { run, result, error, pending, reset, dismissError: () => setError("") };
}
