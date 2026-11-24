const demoAddressWord = "0".repeat(24) + "a".repeat(40);
const calldata = (selector: string, argument: string) =>
  `0x${selector}${demoAddressWord}${argument.padStart(64, "0")}`;

export const demoPresets = [
  { label: "Unlimited Approve", calldata: calldata("095ea7b3", "f".repeat(64)) },
  { label: "NFT Master Key", calldata: calldata("a22cb465", "1") },
  { label: "Safe Revoke", calldata: calldata("a22cb465", "0") },
  { label: "Standard Transfer", calldata: calldata("a9059cbb", "64") },
];
