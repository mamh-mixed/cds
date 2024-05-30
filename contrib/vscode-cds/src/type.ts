import { RequestType } from 'vscode-messenger-common';

export const WorkflowRefresh = { method: 'workflow-refresh'};
export const WorkflowTemplate = { method: 'workflow-template' };

export type WorkflowTemplate = {
  workflowTemplate: any;
};

export type WorkflowData = {
  workflow: any;
};

export type GenerateWorkflowData = {
  parameters: {[key: string]: string} 
};

export type GenerateWorkflowDataResponse = {
  errors: string;
  workflow: any;
}

export type Parameter = {key: string};

export const GenerateWorkflow:  RequestType<GenerateWorkflowData, GenerateWorkflowDataResponse> = { 
  method: 'generateWorkflow'
};