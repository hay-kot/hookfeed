# hookfeed

Hookfeed is a single tenant feed for incoming webhooks in nearly any format.

[Go Reference](https://pkg.go.dev/github.com/hay-kot/hookfeed)

## Install

```bash
go get -u github.com/hay-kot/hookfeed
```

## Features

- [ ] Ecosystem Compatibility
  - [ ] Ntfy
  - [ ] Gotify
  - [ ] Push Over
- [ ] Support basic templating for messages
- [ ] Markdown message support
- [ ] Lua based middleware system
  - [ ] Register nearly any endpoint pattern
  - [ ] Write custom Lua middleware to handle incoming API requests
- [ ] HTML form support
  - [ ] Use hookfeed to add basic form inputs to your website
- [ ] Display raw request
- [ ] Rich message support
  - [ ] First class properties
    - [ ] Job Duration
    - [ ] Tags
    - [ ] Log Lines
  - [ ] Send any JSON data to display as metadata
  - [ ] View RAW Headers and JSON
- [ ] User Authentication
- [ ] Health Checks
  - [ ] Expect a message every `n` amount of time
  - [ ] Alerting when healthcheck failes
  - [ ] Expose endpoint for healthcheck status (gatus integration)
  - [ ] Webhook support for external services

## Examples
