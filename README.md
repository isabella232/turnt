# Turnt

_A simple test client for [Turnstile][]_

Turnt is a cross-platform client for Turnstile that automatically generates
the relevant Rapid7 Authentication headers for quickly testing requests
against services behind a Turnstile server.

Turnt compiles to a single binary with no dependencies
and can be `curl`ed from Github directly.

## Usage
```
$ turnt --help

A helper app to make requests against Turnstile

Usage:
  turnt [url] [flags]

Flags:
      --digest string     Digest signing scheme (default "SHA256")
  -H, --header string     HTTP request headers (default "{}")
  -h, --help              help for turnt
  -u, --identity string   Identity key for the request
  -X, --method string     HTTP request method (default "GET")
  -d, --payload string    HTTP request payload (default "{}")
  -p, --secret string     Secret key for the request

$ turnt -u key -p secret http://localhost:8085/test -d '{"foo": "bar"}'

INFO Using method method=GET
INFO Using URI uri=/test
INFO Using host host=localhost:8085
INFO Using date date=1503595334
INFO Using identity identity=key
INFO Using digest digest=SHA256=Qm/ATwS/j9tYMdw3u7bc9w9jo34FpoxupfY+ha5Xk3Y=
INFO Using signature signature=GyZAQTDki06+KnAiuFZmPvZq5ZAh/9CYXfLv8xZg+KQ=
INFO Using authorization authorization=Rapid7-HMAC-V1-SHA256 a2V5Okd5WkFRVERraTA2K0tuQWl1RlptUHZacTVaQWgvOUNZWGZMdjh4WmcrS1E9
INFO Response:
{
  "headers": {
    "host": "localhost:8085",
    "user-agent": "Go-http-client/1.1",
    "authorization": "Rapid7-HMAC-V1-SHA256 a2V5Okd5WkFRVERraTA2K0tuQWl1RlptUHZacTVaQWgvOUNZWGZMdjh4WmcrS1E9",
    "date": "1503595334",
    "digest": "SHA256=Qm/ATwS/j9tYMdw3u7bc9w9jo34FpoxupfY+ha5Xk3Y=",
    "accept-encoding": "gzip",
    "content-type": "application/octet-stream",
    "r7-correlation-id": "08cbca5d-ad9e-409c-97bc-2dd0dee33cf5",
    "connection": "close",
    "content-length": "14",
    "accept-charset": "utf-8"
  },
  "method": "GET",
  "url": "/test",
  "data": "{\"foo\": \"bar\"}"
}
```

## Development
We use [dep][] to manage dependencies.
You can install it via

```bash
$ go get -u github.com/golang/dep/cmd/dep
```

or, on macOS

```bash
$ brew install dep
$ brew upgrade dep
```
Once you clone the repo, make sure to run `dep ensure` to pull down
the project's (minimal) dependencies.

[Turnstile][] provides a [dummy service][] that echoes the request
it receives. That, combined with a running Turnstile server in front
of it provides a test-bench for development.

## Building
You can build Turnt with any Golang build tool. We prefer
using [Gox][]. It's simple to use:

```bash
$ go get github.com/mitchellh/gox
$ gox

Number of parallel builds: 7

-->   freebsd/amd64: github.com/rapid7/turnt
-->       linux/arm: github.com/rapid7/turnt
-->      netbsd/arm: github.com/rapid7/turnt
-->    darwin/amd64: github.com/rapid7/turnt
-->     freebsd/386: github.com/rapid7/turnt
-->       linux/386: github.com/rapid7/turnt
-->     linux/amd64: github.com/rapid7/turnt
-->   windows/amd64: github.com/rapid7/turnt
-->     openbsd/386: github.com/rapid7/turnt
-->   openbsd/amd64: github.com/rapid7/turnt
-->     windows/386: github.com/rapid7/turnt
-->      netbsd/386: github.com/rapid7/turnt
-->     freebsd/arm: github.com/rapid7/turnt
-->    netbsd/amd64: github.com/rapid7/turnt
-->      darwin/386: github.com/rapid7/turnt
```

[Turnstile]: https://github.com/rapid7/turnstile
[Gox]: https://github.com/mitchellh/gox
[dep]: https://github.com/golang/dep
[dummy service]: https://github.com/rapid7/turnstile/blob/master/bin/svc