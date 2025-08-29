import { Dialog, DialogContent, DialogDescription, DialogPortal, DialogTitle, DialogTrigger } from "@radix-ui/react-dialog";
import { Button } from "./button";
import { DialogHeader } from "./dialog";
import { ResizableHandle, ResizablePanel, ResizablePanelGroup } from "./resizable";

export function ButtonDialog() {
    const defaultInput = "let x1;\nlet x2;\nmax 3 * x1 + 4 * x2;\ns.t. x1 + x2 <= 5;\nx1 >= 0;\nx2 >= 0;";
    const defaultDescription = "<- Declarations\n\n<- Objective (max or min)\n<- Constraints"
    return (
        <Dialog>
            <form>
                <DialogTrigger asChild>
                    <Button
                        variant="ghost"
                        onClick={() => {}}
                    >
                        Example
                    </Button>
                </DialogTrigger>
                <DialogPortal>
                    <span className="fixed inset-0 flex items-center justify-center z-50">
                        <DialogContent
                            aria-label="Input Example"
                            className="
                                w-[90vw] h-[90vh]
                                p-2 rounded-lg
                                bg-white shadow-lg
                                border-4 border-gray-300
                                text-gray-600
                                buttonDialog
                                overflow-auto
                            "
                        >
                            <DialogHeader>
                                <DialogTitle className="text-black font-semibold">Input Example</DialogTitle>
                                <DialogDescription className="flex">
                                    <span className="p-0.5 whitespace-pre-wrap exampleText ml-auto">{defaultInput}</span>
                                    <span className="w-px bg-gray-300 mx-4"></span>
                                    <span className="p-0.5 whitespace-pre-wrap exampleText mr-auto">{defaultDescription}</span>
                                </DialogDescription>
                            </DialogHeader>
                        </DialogContent>
                    </span>
                </DialogPortal>
            </form>
        </Dialog>
    )
}
