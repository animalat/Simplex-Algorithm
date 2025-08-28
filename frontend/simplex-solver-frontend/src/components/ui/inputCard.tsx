import { Button } from "@/components/ui/button"
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { TextArea } from "./TextArea"

export default function InputCard() {
    const defaultInput = "let   x1;\nlet   x2;\nmax 3 * x1 + 4 * x2;\ns.t.   x1 + x2 <= 5;\n       x1 >= 0;\n       x2 >= 0;"
    return (
        <div className="flex justify-center">
            <Card className="w-[75vw]">
                <CardHeader>
                    <CardTitle>Linear Program Solver</CardTitle>
                    <CardDescription>Enter a Linear Program below to solve</CardDescription>
                    <CardAction>
                        <Button variant="ghost">Example</Button>
                    </CardAction>
                </CardHeader>
                <CardContent>
                    <form>
                        <TextArea
                            id="simplexInput"
                            className="h-[400px]"
                            defaultValue={defaultInput}
                        />
                    </form>
                </CardContent>
                <CardFooter>
                    <Button type="submit" className="w-full">
                        Solve
                    </Button>
                </CardFooter>
            </Card>
        </div>
    )
}