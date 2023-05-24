# promsnoop-harness

A simple test harness to measure the overhead of uprobe-based monitoring of prometheus metrics.

## Prerequisites

- Golang
- [promsnoop](https://github.com/acmel/libbpf-bootstrap/tree/prometheusnoop)
- [k6](https://k6.io/docs/get-started/installation/)
- bpftrace

## Usage

The binary is a simple HTTP server with a several routes:

- `/baseline` - has no metrics in the request path
- `/prom/1`, `/prom/10`, `/prom/100` - has 1, 10, or 100 metrics in the request path
- `/usdt/1`, `/usdt/10`, `/usdt/100` - has 1, 10, or 100 metrics in the request path

The server uses the [Gin](https://github.com/gin-gonic/gin) web framework.

To load the webserver and provide metrics for comparison we're using `k6`.
K6 will run a script that hits one of the available HTTP endpoints, issuing one request per second for the duration of the test.

Under these conditions, it should be possible to measure the overhead of uprobes and usdt probes.

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
k6 run ./tests/baseline.js
```

To attach uprobes with bpftrace:

```bash
sudo bpftrace -v ./uprobe.bt -p $(sudo pgrep harness)
```

Note: bpftrace didn't like the symbol name for the uprobe.
You may have to `readelf --symbols harness -W | grep '(*counter).Inc'` to find the correct offset to attach the probe to and update `uprobe.bt` accordingly.

To attach USDT probes with bpftrace:

```bash
sudo bpftrace -v ./usdt.bt -p $(sudo pgrep harness)
```
