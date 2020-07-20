-- name: get-integrations
SELECT integration_id, toggl_credentials, service_credentials, created_at, deleted_at
FROM integrations
WHERE deleted_at IS NULL;

-- name: get-integrations-for-user
SELECT integration_id, toggl_credentials, service_credentials, created_at, deleted_at
FROM integrations
WHERE deleted_at IS NULL
    AND toggl_credentials = $1;

-- name: get-emoji-rules
SELECT project, emoji
FROM slack_emoji_rules
WHERE deleted_at IS NULL
    AND integration_id = $1;

-- name: create-integration
INSERT INTO integrations (name, toggl_credentials, service_credentials, created_at)
VALUES ($1, $2, $3, NOW());