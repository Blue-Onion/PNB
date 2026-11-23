import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    const backend = new URL(
      "/api/check-code",
      process.env.BACKEND_URL || "http://localhost:8080",
    );
    return [{ source: "/api/check-code", destination: backend.href }];
  },
  experimental: {
    agentFeedback: true,
  },
  cacheComponents: true,
  partialPrefetching: true,
};

export default nextConfig;
