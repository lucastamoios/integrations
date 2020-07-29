Slack Integration for Toggl
===========================

Functionalities:
- [x] Update status with what you are working with
- [ ] Set "do not disturb" for given projects
- [ ] (next feature) Message to remember that you are working after a while in slack for some time without tracking. This is interesting as a way to not forget to track time but also as a way to keep a healthy relation with work.
- [ ] Send time tracking command via slack

Running the code
----------------

Create the database and run the migrations:

```
$ make db
$ make migrate
```

After that you should insert your credentials to the database:

```
$ make db-shell
postgres=# \c toggl_integrations;
toggl_integrations=# INSERT INTO integrations (toggl_credentials, service_credentials, created_at) VALUES ('toggl-credentials-encoded-base64', 'xoxp-slack-credentials', NOW());
```

Integration Abstraction
---

An integration is a two-way data pipeline between Toggl and an external service.
Each integration can have multiple features which tell what data should be exchanged.
