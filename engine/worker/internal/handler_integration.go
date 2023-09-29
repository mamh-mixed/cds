package internal

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ovh/cds/engine/worker/pkg/workerruntime"
	"github.com/ovh/cds/sdk"
	"github.com/pkg/errors"
	"github.com/rockbears/log"
)

func getIntegrationHandler(ctx context.Context, wk *CurrentWorker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		integratioName := vars["integration"]

		if integratioName == "" { // workaround just for unit test, in the real life mux let us use var
			log.Debug(ctx, "getting integration name by path %q", r.URL.Path)
			splittedPath := strings.SplitN(r.URL.Path, "/", 3)
			if len(splittedPath) == 3 {
				integratioName = splittedPath[2]
			}
		}

		ctx := workerruntime.SetJobID(ctx, wk.currentJob.wJob.ID)
		ctx = workerruntime.SetStepOrder(ctx, wk.currentJob.currentStepIndex)
		ctx = workerruntime.SetStepName(ctx, wk.currentJob.currentStepName)

		workflowRun, err := wk.client.WorkflowRunGet(wk.currentJob.projectKey, wk.currentJob.workflowName, wk.currentJob.runNumber)
		if err != nil {
			log.ErrorWithStackTrace(ctx, err)
			writeError(w, r, err)
			return
		}

		var integration *sdk.WorkflowProjectIntegration
		for i := range workflowRun.Workflow.Integrations {
			if workflowRun.Workflow.Integrations[i].ProjectIntegration.Name == integratioName {
				integration = &workflowRun.Workflow.Integrations[i]
				break
			}
		}

		if integration == nil {
			err := errors.Errorf("integration %q not found", integratioName)
			log.ErrorWithStackTrace(ctx, err)
			writeError(w, r, err)
			return
		}

		projectIntegration, err := wk.client.ProjectIntegrationGet(wk.currentJob.projectKey, integratioName, true)
		if err != nil {
			log.ErrorWithStackTrace(ctx, err)
			writeError(w, r, err)
			return
		}

		writeJSON(w, projectIntegration, http.StatusOK)
	}
}
