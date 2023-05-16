package sdk

// Context available on ascode action
// "cds": CDSContext - information about workflow instantiation
// "git": GitContext - information about the repository used in a job

type CDSContext struct {
	// Workflow
	Event           map[string]interface{} `json:"event,omitempty"`
	Version         string                 `json:"version,omitempty"`
	RunID           string                 `json:"run_id,omitempty"`
	RunNumber       string                 `json:"run_number,omitempty"`
	RunAttempt      string                 `json:"run_attempt,omitempty"`
	WorkflowRef     string                 `json:"workflow_ref,omitempty"`
	WorkflowSha     string                 `json:"workflow_sha,omitempty"`
	TriggeringActor string                 `json:"triggering_actor,omitempty"`
	// Job
	Job string `json:"triggering_actor,omitempty"`
	// Worker
	Workspace string `json:"workspace,omitempty"`
}

type GitContext struct {
	Hash       string `json:"hash,omitempty"`
	HashShort  string `json:"hash_short,omitempty"`
	Repository string `json:"repository,omitempty"`
	Branch     string `json:"branch,omitempty"`
	Tag        string `json:"tag,omitempty"`
	Author     string `json:"author,omitempty"`
	Message    string `json:"message,omitempty"`
	URL        string `json:"url,omitempty"`
	Server     string `json:"server,omitempty"`
	EventName  string `json:"event_name,omitempty"`
	Connection string `json:"connection,omitempty"`
	SSHKey     string `json:"ssh_key,omitempty"`
	PGPKey     string `json:"pgp_key,omitempty"`
	HttpUser   string `json:"http_user,omitempty"`
}
