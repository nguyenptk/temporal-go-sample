# temporal-go-sample
A Go ft. Docker sample for Temporal workflow with signal handling

## Launch the Temporal

Temporal and Postgres are defined in the docker-compose.yaml file so we just need to run:

```bash
$ docker-compose up --build -d
```

Then 4 modules will be up: Postgres, Temporal, Temporal-admin-tools, and Temporal-ui.

## Launch the Go server

The Go server will manage the workflows of Temporal and there are two types:

- Workflow
- Activity

The workflow will call activities that are implemented as standalone functions. To start the Go server, just run
the simple command:

```bash
$ go run main.go
```

Then the Go server will listen on port 6000 and register all the workflows and activities. We can access the Temporal admin UI on http://localhost:8080 to see the workflows and debug them.

## Test the workflow

The Go server exposes an endpoint `/do-something` to triage a workflow. We can call it via curl:

```bash
curl -i -X POST "http://localhost:6000/do-something"
```

Then the Go server will submit the workflow with 5 activities, and each activity will wait for the signal handling from the previous one to
continue working on its logic.

And then in the Temporal admin UI, we will see a registered workflow `DO_SOMETHING[uuid]`, and in detail, all event histories will be listed down there.

Notes: in real service, we should replace the UUID by a generating UUID function instead of hard coding `uuid` like the sample. Error handling, logging, and metrics collection are omitted in this sample for simplicity.
