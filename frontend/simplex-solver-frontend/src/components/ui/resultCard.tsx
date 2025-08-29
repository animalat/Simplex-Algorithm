import { SimplexResponse } from "@/app/page";
import { Card, CardContent, CardHeader, CardTitle } from "./card";
import SolutionDisplay from "./solutionDisplay";

interface ResultCardProps {
    result: SimplexResponse;
}

export default function ResultCard({ result }: ResultCardProps) {
    return (
        <div className="inputCard flex justify-center">
            <Card className="w-[75vw]">
                <CardHeader>
                    <CardTitle>Simplex Result</CardTitle>
                </CardHeader>
                <CardContent>
                    <SolutionDisplay solution={result.solution} mapping={result.mapping} />
                </CardContent>
            </Card>
        </div>
    );
}
