package jira

import (
	"context"

	"github.com/pkg/errors"
	"gopkg.in/andygrunwald/go-jira.v1"
)

type JiraClient interface {
	AddWatcher(ctx context.Context, issueID string, userName string) (*jira.Response, error)
	AddWorklogRecord(ctx context.Context, issueID string, record *jira.WorklogRecord) (*jira.WorklogRecord, *jira.Response, error)
	AddComment(ctx context.Context, issueID string, comment *jira.Comment) (*jira.Comment, *jira.Response, error)
	FieldGetList(ctx context.Context) ([]jira.Field, *jira.Response, error)
	IssueAddComment(ctx context.Context, issueID string, comment *jira.Comment) (*jira.Comment, *jira.Response, error)
	IssueCreate(ctx context.Context, issue *jira.Issue) (*jira.Issue, *jira.Response, error)
	IssueGet(ctx context.Context, issueID string, options *jira.GetQueryOptions) (*jira.Issue, *jira.Response, error)
	IssueSearchPages(ctx context.Context, jql string, options *jira.SearchOptions, f func(jira.Issue) error) error
	IssueUpdate(ctx context.Context, issue *jira.Issue) (*jira.Issue, *jira.Response, error)
	UpdateAssignee(ctx context.Context, issueID string, assignee *jira.User) (*jira.Response, error)
}

type JiraClientImpl struct {
	C *jira.Client
}

func (j *JiraClientImpl) IssueCreate(ctx context.Context, i *jira.Issue) (*jira.Issue, *jira.Response, error) {
	ires, res, err := j.C.Issue.Create(i)
	return ires, res, errors.WithStack(err)
}

func (j *JiraClientImpl) IssueUpdate(ctx context.Context, i *jira.Issue) (*jira.Issue, *jira.Response, error) {
	ires, res, err := j.C.Issue.Update(i)
	return ires, res, errors.WithStack(err)
}

func (j *JiraClientImpl) FieldGetList(ctx context.Context) ([]jira.Field, *jira.Response, error) {
	ires, res, err := j.C.Field.GetList()
	return ires, res, errors.WithStack(err)
}

func (j *JiraClientImpl) IssueGet(ctx context.Context, issueID string, options *jira.GetQueryOptions) (*jira.Issue, *jira.Response, error) {
	ires, res, err := j.C.Issue.Get(issueID, options)
	return ires, res, errors.WithStack(err)
}

func (j *JiraClientImpl) IssueSearchPages(ctx context.Context, jql string, options *jira.SearchOptions, f func(jira.Issue) error) error {
	return errors.WithStack(j.C.Issue.SearchPages(jql, options, f))
}

func (j *JiraClientImpl) IssueAddComment(ctx context.Context, issueID string, comment *jira.Comment) (*jira.Comment, *jira.Response, error) {
	ires, res, err := j.C.Issue.AddComment(issueID, comment)
	return ires, res, errors.WithStack(err)
}

func (j *JiraClientImpl) AddWorklogRecord(ctx context.Context, issueID string, record *jira.WorklogRecord) (*jira.WorklogRecord, *jira.Response, error) {
	resrecord, res, err := j.C.Issue.AddWorklogRecord(issueID, record)
	return resrecord, res, errors.WithStack(err)
}

func (j *JiraClientImpl) AddComment(ctx context.Context, issueID string, comment *jira.Comment) (*jira.Comment, *jira.Response, error) {
	resc, res, err := j.C.Issue.AddComment(issueID, comment)
	return resc, res, errors.WithStack(err)
}

func (j *JiraClientImpl) AddWatcher(ctx context.Context, issueID string, userName string) (*jira.Response, error) {
	res, err := j.C.Issue.AddWatcher(issueID, userName)
	return res, errors.WithStack(err)
}

func (j *JiraClientImpl) UpdateAssignee(ctx context.Context, issueID string, assignee *jira.User) (*jira.Response, error) {
	res, err := j.C.Issue.UpdateAssignee(issueID, assignee)
	return res, errors.WithStack(err)
}

type JiraWebhook struct {
	Comment struct {
		Author struct {
			Active     bool `json:"active"`
			AvatarUrls struct {
				One6x16  string `json:"16x16"`
				Four8x48 string `json:"48x48"`
			} `json:"avatarUrls"`
			DisplayName  string `json:"displayName"`
			EmailAddress string `json:"emailAddress"`
			Name         string `json:"name"`
			Self         string `json:"self"`
		} `json:"author"`
		Body         string `json:"body"`
		Created      string `json:"created"`
		ID           string `json:"id"`
		Self         string `json:"self"`
		UpdateAuthor struct {
			Active     bool `json:"active"`
			AvatarUrls struct {
				One6x16  string `json:"16x16"`
				Four8x48 string `json:"48x48"`
			} `json:"avatarUrls"`
			DisplayName  string `json:"displayName"`
			EmailAddress string `json:"emailAddress"`
			Name         string `json:"name"`
			Self         string `json:"self"`
		} `json:"updateAuthor"`
		Updated string `json:"updated"`
	} `json:"comment"`
	ID           int64     `json:"id"`
	Issue        JiraIssue `json:"issue"`
	Timestamp    int64     `json:"timestamp"`
	User         JiraUser  `json:"user"`
	WebhookEvent string    `json:"webhookEvent"`
}

type JiraUser struct {
	Active     bool `json:"active"`
	AvatarUrls struct {
		One6x16  string `json:"16x16"`
		Four8x48 string `json:"48x48"`
	} `json:"avatarUrls"`
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
	Key          string `json:"key"`
	Name         string `json:"name"`
	Self         string `json:"self"`
}

type JiraIssue struct {
	Fields JiraFields `json:"fields"`
	ID     string     `json:"id"`
	Key    string     `json:"key"`
	Self   string     `json:"self"`
}

type JiraFields struct {
	Created     string          `json:"created"`
	Description string          `json:"description"`
	Labels      []string        `json:"labels"`
	Components  []JiraComponent `json:"components"`
	Priority    JiraPriority    `json:"priority"`
	Summary     string          `json:"summary"`
}

type JiraComponent struct {
	Description string `json:"description"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Self        string `json:"self"`
}

type JiraPriority struct {
	IconURL string `json:"iconUrl"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Self    string `json:"self"`
}
