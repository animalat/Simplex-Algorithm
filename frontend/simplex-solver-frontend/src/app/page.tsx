"use client"

import { useState } from "react"
import InputCard from "@/components/ui/inputCard"

export interface SimplexResponse {
  solution: number[];
  resultType: string;
  certificate: number[];
  mapping: Record<number, string>;
}

export default function Home() {
  const [simplexResult, setSimplexResult] = useState<SimplexResponse | null>(null);
  return (
    <div>
      <InputCard onSimplexResponse={setSimplexResult} />
    </div>
  );
}
