import { Component, ViewEncapsulation } from '@angular/core';
import { Messenger, VsCodeApi } from 'vscode-messenger-webview';
import { load, LoadOptions } from 'js-yaml';
import { HOST_EXTENSION } from 'vscode-messenger-common';
import { GenerateWorkflow, Parameter, WorkflowData, WorkflowRefresh, WorkflowTemplate } from '../../../../src/type';

export declare function acquireVsCodeApi(): VsCodeApi;
const vsCodeApi = acquireVsCodeApi();

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class AppComponent {
  title = 'cds.workflow.preview';

  workflow: string = '';
  workflowError: string = '';

  workflowTemplate: any = {};
  templateError: string = '';
  templatesInputs: {[key: string]: string} = {};
  collapsedParameters: boolean = false;
  
  viewMessenger: Messenger;

  constructor() {
    this.viewMessenger = new Messenger(vsCodeApi);
    this.viewMessenger.onNotification(WorkflowRefresh, e => {
      this.workflow =  (e as WorkflowData).workflow;
    });
    this.viewMessenger.onNotification(WorkflowTemplate, e => {

      let data = (e as WorkflowTemplate).workflowTemplate;
      if (data && data !== '') {
        this.receivedWorkflowTemplate(data);
      }
  });
    this.viewMessenger.start();
  }

  toggleTemplateParameters(): void {
    this.collapsedParameters = !this.collapsedParameters;
  }

  receivedWorkflowTemplate(data: any): void {
    try {
      this.templateError = '';
      this.workflowTemplate = load(data, <LoadOptions>{
          onWarning: () => {}
      });
      let oldParams = structuredClone(this.templatesInputs);
      if (this.workflowTemplate['parameters']) {
        // Add new params
        this.workflowTemplate['parameters'].forEach((p: Parameter) => {
          if (!this.templatesInputs[p.key]) {
            this.templatesInputs[p.key] = oldParams[p.key]?oldParams[p.key]: '';
          }
        });
      }
    } catch (e: any) {
      this.templateError = e.message;
    }
  }

  generateWorkflow(): void {
    this.getWorkflowFromExtension();
  }

  async getWorkflowFromExtension() {
    const generatedWorkflow = await this.viewMessenger.sendRequest(GenerateWorkflow, HOST_EXTENSION, {parameters: this.templatesInputs});
    this.workflow = generatedWorkflow.workflow;
    this.workflowError = generatedWorkflow.errors;
  }
}
