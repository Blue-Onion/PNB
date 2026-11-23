const CHECK_CODE_ROUTE = "/api/check-code";

export type Report = {
  report_id: string;
  transaction: {
    function: string;
    address?: string;
    amount?: string;
    approved?: boolean;
  };
  final_risk: string;
  findings: { engine: string; severity: string; message: string }[];
  ai_analysis: {
    summary: string;
    confidence: string;
    recommendations: string[] | null;
    risk_points: { reason: string; explanation: string; severity: string }[] | null;
  };
  analysis: {
    status: string;
    simulation_performed: boolean;
    reputation_source: string;
  };
};

type CheckCodeResponse = { result: Report | null; error: string };

export async function checkCode(hexcode: string): Promise<Report> {
  const unavailable = "Could not reach the analysis API. Check that the backend is running.";
  const response = await fetch(CHECK_CODE_ROUTE, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ hexcode }),
  }).catch(() => {
    throw new Error(unavailable);
  });

  const payload: CheckCodeResponse | null = await response.json().catch(() => null);
  if (payload?.error) throw new Error(payload.error);
  if (!response.ok || !payload?.result) throw new Error(unavailable);
  return payload.result;
}
