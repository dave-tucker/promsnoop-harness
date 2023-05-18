# promsnoop-harness

A simple test harness to measure the overhead of uprobe-based monitoring of prometheus metrics.

## Prerequisites

- Golang
- [promsnoop](https://github.com/acmel/libbpf-bootstrap/tree/prometheusnoop)
- [k6](https://k6.io/docs/get-started/installation/)

## Usage

The binary is a simple HTTP server with a single route `/products/{id}` and returns a JSON response like `{ "productID": "1234" }`.
The server uses the [Gin](https://github.com/gin-gonic/gin) web framework and is instrumented using [gin-metrics](https://github.com/penglongli/gin-metrics).
The library sets 7 metrics after each HTTP request is handled - 6 counters and 1 histogram.

To load the webserver and provide metrics for comparison we're using `k6`.
K6 will run a script that hits our HTTP endpoint with a configurable number of concurrent users.

Under these conditions, it should be possible to measure the overhead of promsnoop by comparing the performance of the web server with and without promsnoop enabled.

Build the test binary:

```bash
go build -o harness main.go
```

Launch the webserver:

```bash
./harness -cpuprofile cpu.prof
```

Run the k6 script:

```bash
k6 run --vus 1000 --duration 5s script.js
```

To attach promsnoop to the webserver, run the following command:

```bash
sudo promsnoop -b harness
```
