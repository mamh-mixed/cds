package internal

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ovh/cds/engine/test"
	"github.com/ovh/cds/sdk"
	"github.com/ovh/cds/sdk/cdsclient/mock_cdsclient"
	"github.com/rockbears/log"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getIntegrationHandler(t *testing.T) {
	log.Factory = log.NewTestingWrapper(t)
	// Create test directory for current test
	fs := afero.NewOsFs()
	basedir := "test-" + test.GetTestName(t) + "-" + sdk.RandomString(10) + "-" + fmt.Sprintf("%d", time.Now().Unix())
	t.Logf("Creating worker basedir at %s", basedir)
	require.NoError(t, fs.MkdirAll(basedir, os.FileMode(0755)))

	// Setup test worker
	wk := &CurrentWorker{basedir: afero.NewBasePathFs(fs, basedir)}
	wk.currentJob.wJob = &sdk.WorkflowNodeJobRun{ID: 1}
	wk.currentJob.projectKey = "proj"
	wk.currentJob.workflowName = "wkf"
	wk.currentJob.runNumber = 1

	// Prepare mock client for cds workers
	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })
	m := mock_cdsclient.NewMockWorkerInterface(ctrl)
	wk.client = m

	jiraCfg := sdk.JIRAIntegration.DefaultConfig.Clone()
	jiraCfg.SetValue("url", "my_url")
	jiraCfg.SetValue("username", "my_username")
	jiraCfg.SetValue("password", "my_password")

	m.EXPECT().WorkflowRunGet("proj", "wkf", int64(1)).Return(&sdk.WorkflowRun{
		Number: 1,
		Workflow: sdk.Workflow{
			Integrations: []sdk.WorkflowProjectIntegration{
				{
					ProjectIntegration: sdk.ProjectIntegration{
						Name:   "jira",
						Config: jiraCfg,
					},
				},
			},
		},
	}, nil)

	m.EXPECT().ProjectIntegrationGet("proj", "jira", true).Return(
		sdk.ProjectIntegration{
			Name:   "jira",
			Model:  sdk.JIRAIntegration,
			Config: jiraCfg,
		}, nil,
	)

	req, err := http.NewRequest(http.MethodGet, "/integrations/jira", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	getIntegrationHandler(context.TODO(), wk)(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	if w.Code != 200 {
		cdsError := sdk.DecodeError(w.Body.Bytes())
		t.Log(cdsError.Error())
		t.FailNow()
	}
}
