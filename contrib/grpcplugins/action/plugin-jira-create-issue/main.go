package main

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rockbears/log"
	gojira "gopkg.in/andygrunwald/go-jira.v1"

	"github.com/ovh/cds/contrib/grpcplugins"
	"github.com/ovh/cds/contrib/integrations/jira"
	"github.com/ovh/cds/sdk"
	"github.com/ovh/cds/sdk/grpcplugin/actionplugin"
)

/* Inside contrib/grpcplugins/action
$ make build plugin-arsenal-delete-alternative
$ make publish plugin-arsenal-delete-alternative
*/

type jiraPlugin struct {
	actionplugin.Common
}

func (e *jiraPlugin) Manifest(ctx context.Context, _ *empty.Empty) (*actionplugin.ActionPluginManifest, error) {
	return &actionplugin.ActionPluginManifest{
		Name:        "Jira Create Issue Plugin",
		Author:      "François Samin",
		Description: "Create a JIRA issue",
		Version:     sdk.VERSION,
	}, nil
}

func (e *jiraPlugin) Run(ctx context.Context, q *actionplugin.ActionQuery) (*actionplugin.ActionResult, error) {
	log.Factory = log.NewStdWrapper(log.StdWrapperOptions{DisableTimestamp: true, Level: log.LevelInfo})
	log.UnregisterField(log.FieldCaller, log.FieldSourceFile, log.FieldSourceLine, log.FieldStackTrace)

	var (
		integrationName = getStringOption(q, "cds.integration.issue_tracker")
	)
	if integrationName == "" {
		return fail("missing jira issue_tracker integration on workflow")
	}

	projectIntegration, err := grpcplugins.GetProjectIntegration(e.HTTPPort, integrationName)
	if err != nil {
		failErr(err)
	}

	url := projectIntegration.Config["url"].Value
	username := projectIntegration.Config["username"].Value
	password := projectIntegration.Config["password"].Value

	jiraAuth := gojira.BasicAuthTransport{
		Username: username,
		Password: password,
	}

	goJiraClient, err := gojira.NewClient(jiraAuth.Client(), url)
	if err != nil {
		failErr(err)
	}

	jiraClient := &jira.JiraClientImpl{C: goJiraClient}

	var (
		typ         = getStringOption(q, "type")
		projectKey  = getStringOption(q, "projectKey")
		summary     = getStringOption(q, "summary")
		description = getStringOption(q, "description")
	)

	issue := gojira.Issue{
		Fields: &gojira.IssueFields{
			Summary:     summary,
			Description: description,
			Type: gojira.IssueType{
				Name: typ,
			},
			Project: gojira.Project{
				Key: projectKey,
			},
		},
	}

	i, _, err := jiraClient.IssueCreate(ctx, &issue)
	if err != nil {
		failErr(err)
	}

	log.Info(ctx, "issue %s created", i.ID)

	return &actionplugin.ActionResult{
		Status: sdk.StatusSuccess,
	}, nil
}

func getStringOption(q *actionplugin.ActionQuery, keys ...string) string {
	for _, k := range keys {
		if v, exists := q.GetOptions()[k]; exists {
			return v
		}
	}
	return ""
}

func fail(format string, args ...interface{}) (*actionplugin.ActionResult, error) {
	return failErr(fmt.Errorf(format, args...))
}

func failErr(err error) (*actionplugin.ActionResult, error) {
	fmt.Println("Error:", err)
	return &actionplugin.ActionResult{
		Details: err.Error(),
		Status:  sdk.StatusFail,
	}, nil
}

func main() {
	e := jiraPlugin{}
	if err := actionplugin.Start(context.Background(), &e); err != nil {
		panic(err)
	}
}
