package cdsclient

import (
	"context"
	"fmt"
	"github.com/ovh/cds/sdk"
	"net/url"
)

func (c *client) WorkflowV2Run(projectKey, vcsIdentifier, repoIdentifier, wkfName string, mods ...RequestModifier) (*sdk.V2WorkflowRun, error) {
	var run sdk.V2WorkflowRun
	path := fmt.Sprintf("/v2/project/%s/vcs/%s/repository/%s/workflow/%s/run", projectKey, url.PathEscape(vcsIdentifier), url.PathEscape(repoIdentifier), wkfName)
	_, _, _, err := c.RequestJSON(context.Background(), "POST", path, nil, &run, mods...)
	if err != nil {
		return nil, err
	}
	return &run, nil
}

func (c *client) WorkflowV2RunStatus(projectKey, vcsIdentifier, repoIdentifier, wkfName string, runNumber int64) (*sdk.V2WorkflowRun, error) {
	var run sdk.V2WorkflowRun
	path := fmt.Sprintf("/v2/project/%s/vcs/%s/repository/%s/workflow/%s/run/%d", projectKey, url.PathEscape(vcsIdentifier), url.PathEscape(repoIdentifier), wkfName, runNumber)
	_, _, _, err := c.RequestJSON(context.Background(), "GET", path, nil, &run)
	if err != nil {
		return nil, err
	}
	return &run, nil
}
