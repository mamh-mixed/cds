export interface WorkflowGenerateResponse {
    readonly errors: string[];
    readonly workflow: any;
}

export interface WorkflowGenerateRequest {
    filePath: string;
    params: {[key: string]: string};
}
