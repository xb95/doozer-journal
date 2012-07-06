doozer-journal
==================

Snapshots, mutation journaling and recovery of doozerd coordinator state.

## Usage

```
Usage: doozer-journal [globals] command

Globals:
  -file   location of backup file (./doozerd.log)
  -uri    doozerd cluster URI     (doozer:?ca=localhost:8046)

Commands:
  journal    takes an initial snapshot & journals mutations
  restore    replays journal
  snapshot   makes a snapshot and exits
```

### Conventions

This repository follows the code conventions dictated by [gofmt](http://golang.org/cmd/gofmt/).
