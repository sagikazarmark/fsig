# File Signal (fsig)

[![Build Status](https://img.shields.io/travis/sagikazarmark/fsig.svg?style=flat-square)](https://travis-ci.org/sagikazarmark/fsig)
[![Go Report Card](https://goreportcard.com/badge/github.com/sagikazarmark/fsig?style=flat-square)](https://goreportcard.com/report/github.com/sagikazarmark/fsig)
[![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/sagikazarmark/fsig)


**Send signals to a child process upon file changes**

This project was born because of the need to reload applications upon Kubernetes ConfigMap changes,
but it can be used without the containerization stuff as well.

Fsig is heavily inspired by [configmap-reload](https://github.com/jimmidyson/configmap-reload)
which provides a similar use case for modern application (like Prometheus) run on Kubernetes.


## Usage

```
usage: fsig --watch=WATCH [<flags>] <signal> <cmd> [<args>...]

Send signals to a child process upon file changes

Flags:
      --help             Show context-sensitive help (also try --help-long and --help-man).
  -w, --watch=WATCH ...  Watched directory (at least one)
      --version          Show application version.

Args:
  <signal>  Signal to be sent to the child process
  <cmd>     Child process command
  [<args>]  Child process arguments
```

### Example

```bash
$ fsig -w watched/dir HUP -- ./my_program --arg
```


## Installation

Download a precompiled binary for the [latest](https://github.com/sagikazarmark/fsig/releases/latest) version.

### Ubuntu images

```dockerfile
RUN apt-get update && apt-get install -y wget

ENV FSIG_VERSION 0.3.0
RUN wget https://github.com/sagikazarmark/fsig/releases/download/v${FSIG_VERSION}/fsig_${FSIG_VERSION}_linux_amd64.tar.gz \
    && tar -C /usr/local/bin -xzvf fsig_${FSIG_VERSION}_linux_amd64.tar.gz fsig \
    && rm fsig_linux_amd64.tar.gz
```

### Alpine images

```dockerfile
RUN apk add --no-cache openssl

ENV FSIG_VERSION 0.3.0
RUN wget https://github.com/sagikazarmark/fsig/releases/download/v${FSIG_VERSION}/fsig_${FSIG_VERSION}_linux_amd64.tar.gz \
    && tar -C /usr/local/bin -xzvf fsig_${FSIG_VERSION}_linux_amd64.tar.gz fsig \
    && rm fsig_linux_amd64.tar.gz
```


## Alternatives

fsig might not always fit your use case.
The following alternatives provide similar solutions to the problem fsig tries to solve.


### Shell script


**Pros:**

- very simple
- has only two dependencies (`bash`, `inotifywait`)


**Cons:**

- works by putting processes in the background


```bash
#!/bin/bash

{
    echo "Starting [APPLICATION]"
    start_application "$@"
} &

pid=$!
watches=${WATCH_PATHS:-"/path/to/watched/file"}

echo "Setting up watches for ${watches[@]}"

{
    inotifywait -e modify,move,create,delete --timefmt '%d/%m/%y %H:%M' -m --format '%T' ${watches[@]} | while read date time; do
        echo "File change detected at ${time} on ${date}"
        kill -s HUP $pid
    done
    
    echo "Watching file changes failed, killing application"
    kill -TERM $pid
} &

wait $pid || exit 1
```

#### Install on Alpine

```bash
$ apk add bash inotify-tools
```

#### Install on Debian

```bash
$ apt-get install bash inotify-tools
```


### Entr

[entr](http://entrproject.org/) is a tool similar to `inotifywait` with the purpose of providing better user
experience when used in scripts.

**Pros:**

- tiny dependency
- very simple usage (compared to the script above)

**Cons:**

- cannot watch directories without exiting first which makes it impossible to use in a Docker container


## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.
