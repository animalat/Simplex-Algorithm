import Latex from "react-latex-next";

interface SolutionDisplayProps {
    solution: number[];
    mapping: Record<number, string>;
}

export default function SolutionDisplay({ solution, mapping }: SolutionDisplayProps) {
    const lhs = solution.map((_, idx) => (`\\text{${mapping[idx]}}`)).join(' \\\\ ');
    const rhs = solution.map(val => `${val}`).join(' \\\\ ');
    return (
        <Latex> 
            {String.raw`\[
                \begin{pmatrix} ${lhs} \end{pmatrix} =
                \begin{pmatrix} ${rhs} \end{pmatrix}
            \]`}
        </Latex>
    );
}