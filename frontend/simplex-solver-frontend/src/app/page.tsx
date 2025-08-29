"use client"

import { useState } from "react"
import InputCard from "@/components/ui/inputCard"
import ResultCard from "@/components/ui/resultCard";

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
      {simplexResult && <ResultCard result={simplexResult} />}
    </div>
  );
}
