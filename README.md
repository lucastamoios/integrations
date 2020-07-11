Slack Integration for Toggl
===========================

Functionalities:
- Update status with what you are working with
- (next feature) Message to remember that you are working after a while in slack for some time without tracking. This is interesting as a way to not forget to track time but also as a way to keep a healthy relation with work.

Running the code
----------------

Download the module for migrations, create the database and run the migrations:
```bash
go get github.com/golang-migrate/migrate
createdb toggl_integrations
migrate -database "postgres://$USER:@localhost/toggl_integrations?sslmode=disable" -path db/migrations up
```

After that you should insert your credentials to the database:
```sql
INSERT INTO integrations (toggl_credentials, service_credentials, created_at) VALUES ('toggl-credentials-encoded-base64', 'xoxp-slack-credentials', NOW());
```

Integration Abstraction
---

An integration is a two-way data pipeline between Toggl and an external service.
Each integration can have multiple features which tell what data should be exchanged.

