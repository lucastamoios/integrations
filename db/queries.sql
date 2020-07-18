-- name: get-integrations
SELECT integration_id, toggl_credentials, service_credentials, created_at, deleted_at
FROM integrations
WHERE deleted_at IS NULL;

-- name: get-emoji-rules
SELECT project, emoji
FROM slack_emoji_rules
WHERE deleted_at IS NULL
    AND integration_id = $1;