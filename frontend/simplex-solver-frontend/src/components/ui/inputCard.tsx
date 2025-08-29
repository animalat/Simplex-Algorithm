"use client"

import { useState } from "react"
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
import { Textarea } from "@/components/ui/textarea"
import { SimplexResponse } from "@/app/page"
import axios from "axios"
import { ButtonDialog } from "./buttonDialog"

interface InputCardProps {
    onSimplexResponse: (data: SimplexResponse) => void;
}

export default function InputCard({ onSimplexResponse }: InputCardProps) {
    const defaultInput = "let x1;\nlet x2;\nmax 3 * x1 + 4 * x2;\ns.t. x1 + x2 <= 5;\nx1 >= 0;\nx2 >= 0;";
    const URL = "http://localhost:8080/solve"

    const [value, setValue] = useState(defaultInput);
    const handleSolve = () => {
        axios.post(URL, value, {
            headers: {
                'Content-Type': 'text/plain'
            }
        }).then((response) => {
            onSimplexResponse(response.data);
        }).catch((error) => {
            if (error.response) {
                console.log("Status: ", error.response.status)
                console.log("Error message: ", error.response.data)
            } else if (error.request) {
                console.log("No response: ", error.request)
            } else {
                console.log("No response, no request: ", error)
            }
        });
    }
    
    return (
        <div className="inputCard flex justify-center py-10">
            <Card className="w-[75vw]">
                <CardHeader>
                    <CardTitle>Linear Program Solver</CardTitle>
                    <CardDescription>Enter a Linear Program below to solve</CardDescription>
                    <CardAction>
                        <ButtonDialog />
                    </CardAction>
                </CardHeader>
                <CardContent>
                    <form>
                        <Textarea
                            id="simplexInput"
                            className="h-[400px]"
                            defaultValue={defaultInput}
                            onChange={(e) => setValue(e.target.value)}
                            data-gramm="false"
                        />
                    </form>
                </CardContent>
                <CardFooter>
                    <Button
                        type="submit"
                        className="w-full"
                        onClick={handleSolve}
                    >
                        Solve
                    </Button>
                </CardFooter>
            </Card>
        </div>
    );
}