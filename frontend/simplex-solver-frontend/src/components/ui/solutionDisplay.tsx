import Latex from "react-latex-next";

interface SolutionDisplayProps {
    solution: number[];
    resultType: string;
    certificate: number[];
    mapping: Record<number, string>;
}

export default function SolutionDisplay({ solution, resultType, certificate, mapping }: SolutionDisplayProps) {
    const createLHS = (vector: number[]) => {
        const lhs = vector.map((_, idx) => (`\\text{${mapping[idx]}}`)).join(' \\\\ ');
        return lhs
    }
    const createRHS = (vector: number[]) => {
        const rhs = vector.map(val => `${val}`).join(' \\\\ ');
        return rhs
    }

    return (
        <div>
            Result is {resultType}
            {resultType != "infeasible" && (
                <div className="mt-4">
                    <p>Certificate:</p>
                    <Latex>
                        {String.raw`\[
\begin{pmatrix} ${createLHS(solution)} \end{pmatrix} =
\begin{pmatrix} ${createRHS(solution)} \end{pmatrix}
                        \]`}
                    </Latex>
                </div>
            )}
            {resultType == "unbounded" && (
                <div className="mt-4">
                    <p>Certificate:</p>
                    <Latex>
                        {String.raw`\[
\begin{pmatrix} ${createLHS(certificate)} \end{pmatrix} =
\begin{pmatrix} ${createRHS(certificate)} \end{pmatrix}
                        \]`}
                    </Latex>
                </div>
            )}
        </div>
    );
}
