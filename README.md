# Bigtable emulator data dumper

## What?

This api reads data from [firebase bigtable emulator](https://cloud.google.com/bigtable/docs/emulator) using the [go/bigtable package](https://pkg.go.dev/cloud.google.com/go/bigtable) and dumps it all out in json format.

## Why?

The Google Firebase Emulator Suite for local development has a nice UI for Authentication, Firestore and Storage to view your local data. However, there is unfortunately no such UI for Bigtable at the moment.

Output is done via json since Bigtable is a NoSQL database, structured output in a custom format would only make it more difficult to read.

## How?

1. Copy and fill env vars

`cp .env.example .env`

2. Build

`docker compose build`

3. Run

`docker compose up -d`

4. View

Open in browser (preferably firefox, since it has a nice json treeview out of the box)

http://localhost:8765/ for all the data for every table,

http://localhost:8765/:tableName for all the data for a single table by name

or pipe straight into jq:

`curl localhost:8765 | jq`
