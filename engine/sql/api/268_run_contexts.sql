-- +migrate Up
ALTER TABLE workflow_node_run ADD COLUMN contexts JSONB;

-- +migrate Down
ALTER TABLE workflow_node_run DROP COLUMN contexts;

