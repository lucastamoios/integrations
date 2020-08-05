-- name: get-all-integrations
SELECT integration_id, toggl_credentials, service_credentials, created_at, deleted_at
FROM integrations
WHERE deleted_at IS NULL;

-- name: get-integrations-for-user
SELECT integration_id, toggl_credentials, service_credentials, created_at, deleted_at
FROM integrations
WHERE deleted_at IS NULL
    AND toggl_credentials = $1;

-- name: get-rules
SELECT project, emoji, do_not_disturb
FROM slack_rules
WHERE deleted_at IS NULL
    AND integration_id = $1;

-- name: create-integration
INSERT INTO integrations (name, toggl_credentials, service_credentials, created_at)
VALUES ($1, $2, $3, NOW());

-- name: create-new-rule
INSERT INTO slack_rules (integration_id, project, client, description, emoji, do_not_disturb)
VALUES ($1, $2, $3, $4, $5, $6);