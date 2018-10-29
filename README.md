# runlog

CI containers need to log _somewhere_. This repo is for the code that creates
that somewhere.

## Quick Start

```
source env/local

# build the server
run build-server

# run the server with logstash ingest
docker-compose up

# build the client
run build-client

# get logs
./runlogq get log 1

# test log ingest/stream (if it doesn't work right away,
# give it a minute because logstash takes a bit to start)
echo $TEST_LOG | nc localhost 12345
```

## Architecture

To handle ingest, Logstash is being used because it's already a well-proven
solution and allows storage in different backends.

Querying is being kept separate from the ingest. This is because this system
is designed to be ingest heavy.

Both the ingest service and the query service need to share the same file
system so that the query service can stream logs while they are being ingested.
