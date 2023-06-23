-- +migrate Up
ALTER TABLE hatchery ADD COLUMN model_type VARCHAR(255);
ALTER TABLE v2_workflow_run ADD COLUMN header JSONB;
ALTER TABLE v2_workflow_run_job ADD COLUMN header JSONB;
ALTER TABLE v2_workflow_run_job ADD COLUMN region VARCHAR(255);
ALTER TABLE v2_workflow_run_job ADD COLUMN model_type VARCHAR(255);
ALTER TABLE v2_workflow_run_job ADD COLUMN project_key VARCHAR(255);
SELECT create_foreign_key_idx_cascade('FK_v2_workflow_run_job_project', 'v2_workflow_run_job', 'project', 'project_key', 'projectkey');

CREATE TABLE v2_workflow_run_job_info (
  "id"                  uuid PRIMARY KEY,
  "workflow_run_id"     uuid NOT NULL,
  "workflow_run_job_id" uuid NOT NULL,
  "issued_at"           TIMESTAMP WITH TIME ZONE DEFAULT LOCALTIMESTAMP,
  "level"               VARCHAR(50),
  "message"             TEXT
);
SELECT create_foreign_key_idx_cascade('FK_v2_workflow_run_job_info_run', 'v2_workflow_run_job_info', 'v2_workflow_run', 'workflow_run_id', 'id');
SELECT create_foreign_key_idx_cascade('FK_v2_workflow_run_job_info_job', 'v2_workflow_run_job_info', 'v2_workflow_run_job', 'workflow_run_job_id', 'id');

-- +migrate Down
ALTER TABLE hatchery DROP COLUMN model_type;
ALTER TABLE v2_workflow_run_job DROP COLUMN header;
ALTER TABLE v2_workflow_run_job DROP COLUMN project_key;
ALTER TABLE v2_workflow_run_job DROP COLUMN region;
ALTER TABLE v2_workflow_run_job DROP COLUMN model_type;
ALTER TABLE v2_workflow_run DROP COLUMN header;
DROP TABLE v2_workflow_run_job_info;

