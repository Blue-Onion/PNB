import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "txShield — Ethereum Transaction Analyzer",
  description: "Review Ethereum transaction calldata before signing.",
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
