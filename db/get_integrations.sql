-- name: get-integrations
SELECT *
FROM integrations
WHERE deleted_at IS NULL;