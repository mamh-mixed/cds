-- +migrate Up
ALTER TABLE "integration_model" ADD COLUMN issue_tracker boolean DEFAULT false;

-- +migrate Down
ALTER TABLE "integration_model" DROP COLUMN issue_tracker;
