# runlog

CI containers need to log _somewhere_. This repo is for the code that creates
that somewhere.

## Quick Start

```
source env/local
docker-compose up
echo $TEST_LOG | nc localhost 12345
```

## Architecture

To handle ingest, Logstash is being used because it's already a well-proven
solution and allows storage in different backends.
