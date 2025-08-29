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
                    <div className="fixed inset-0 flex items-center justify-center z-50">
                        <DialogContent
                            aria-label="Input Example"
                            className="
                                w-[60vw] h-[600px]
                                p-6 rounded-lg
                                bg-white shadow-lg
                                border-4 border-gray-300
                                text-gray-600
                            "
                        >
                            <DialogHeader>
                                <DialogTitle>Input Example</DialogTitle>
                                <DialogDescription>
                                    <ResizablePanelGroup
                                        direction="horizontal"
                                        className="rounded-lg border md:min-w-[450px]"
                                    >
                                        <ResizablePanel defaultSize={50}>
                                            <div className="h-full p-6 whitespace-pre-wrap exampleText">
                                                {defaultInput}
                                            </div>
                                        </ResizablePanel>
                                        <ResizableHandle />
                                        <ResizablePanel defaultSize={50}>
                                            <div className="h-full p-6 whitespace-pre-wrap exampleText">
                                                {defaultDescription}
                                            </div>
                                        </ResizablePanel>
                                    </ResizablePanelGroup>
                                </DialogDescription>
                            </DialogHeader>
                        </DialogContent>
                    </div>
                </DialogPortal>
            </form>
        </Dialog>
    )
}