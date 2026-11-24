export const RISK_LEVELS = [
  "UNKNOWN",
  "SAFE",
  "LOW",
  "REVIEW",
  "MEDIUM",
  "HIGH",
] as const;
export type RiskLevel = (typeof RISK_LEVELS)[number];

export type Report = {
  report_id: string;
  transaction: {
    function: "transfer" | "approve" | "setApprovalForAll" | "unknown";
    address?: string;
    amount?: string;
    approved?: boolean;
  };
  final_risk: RiskLevel;
  findings: {
    engine: "PERMISSION_CHECKER" | "SCAM_ADDRESS_CHECKER";
    severity: RiskLevel;
    message: string;
  }[];
  ai_analysis: {
    summary: string;
    confidence: string;
    recommendations: string[] | null;
    risk_points:
      | { reason: string; explanation: string; severity: string }[]
      | null;
  };
  analysis: {
    status: "COMPLETED" | "AI_UNAVAILABLE";
    simulation_performed: boolean;
    reputation_source: "LOCAL_TEST_DATABASE";
  };
};

type CheckCodeResponse =
  | { Success: true; Data: Report }
  | { Success: false; Data: string };

export async function checkCode(
  input: string,
): Promise<{ report: Report; raw: string }> {
  const calldata = input.trim();
  if (!/^0x(?:[\da-fA-F]{2}){4,}$/.test(calldata)) {
    throw new Error(
      "Enter 0x-prefixed hex with a four-byte selector and complete byte pairs.",
    );
  }

  const baseUrl = process.env.BACKEND_URL?.trim().replace(/\/+$/, "") ?? "";
  const controller = new AbortController();
  // Give the backend's 30-second AI request time to return a local report.
  const timeout = setTimeout(() => controller.abort(), 40_000);
  let response: Response;
  let raw: string;
  console.log(`${baseUrl}/api/check/${encodeURIComponent(calldata)}`);

  try {
    response = await fetch(
      `${baseUrl}/api/check/${encodeURIComponent(calldata)}`,
      {
        headers: { Accept: "application/json" },
        cache: "no-store",
        signal: controller.signal,
      },
    );
    raw = await response.text();
  } catch {
    if (controller.signal.aborted) {
      throw new Error("The analysis timed out. Try again.");
    }
    throw new Error(
      "Could not reach the analysis API. Check the backend URL and confirm the server is running.",
    );
  } finally {
    clearTimeout(timeout);
  }

  let payload: CheckCodeResponse | null = null;
  try {
    payload = JSON.parse(raw);
  } catch {
    // The HTTP status can still explain errors from a proxy or server.
  }
  if (payload?.Success === false && typeof payload.Data === "string") {
    throw new Error(payload.Data);
  }
  if (!response.ok)
    throw new Error(`The analysis API returned HTTP ${response.status}.`);
  if (
    payload?.Success !== true ||
    !payload.Data?.transaction ||
    !Array.isArray(payload.Data.findings) ||
    !payload.Data.ai_analysis ||
    !payload.Data.analysis ||
    !RISK_LEVELS.includes(payload.Data.final_risk)
  ) {
    throw new Error("The analysis API returned an unexpected response.");
  }
  return { report: payload.Data, raw };
}
