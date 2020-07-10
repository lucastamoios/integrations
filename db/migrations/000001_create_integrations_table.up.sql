CREATE TABLE integrations (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    toggl_credentials VARCHAR(1024) NOT NULL,
    service_credentials VARCHAR(1024) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);