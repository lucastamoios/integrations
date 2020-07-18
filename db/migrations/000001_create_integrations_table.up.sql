CREATE TABLE integrations (
    integration_id BIGSERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(255),
    toggl_credentials VARCHAR(1024) NOT NULL,
    service_credentials VARCHAR(1024) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE TABLE slack_emoji_rules (
    slack_emoji_rules_id BIGSERIAL NOT NULL PRIMARY KEY,
    integration_id INTEGER REFERENCES integrations(integration_id),
    project VARCHAR(255),
    client VARCHAR(255),
    description VARCHAR(255),
    emoji VARCHAR(255) NOT NULL,
    deleted_at TIMESTAMP
)