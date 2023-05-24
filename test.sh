#!/bin/bash

set -ex

pkill harness || true
./harness > /dev/null 2>&1 &

k6 run ./tests/baseline.js

sudo pkill bpftrace || true
sudo bpftrace -v ./usdt.bt -p $(sudo pgrep harness) > /dev/null 2>&1 &

k6 run ./tests/usdt1.js
k6 run ./tests/usdt10.js
k6 run ./tests/usdt100.js

sudo pkill bpftrace || true
sudo bpftrace -v ./uprobe.bt -p $(sudo pgrep harness) > /dev/null 2>&1 &

k6 run ./tests/prom1.js
k6 run ./tests/prom10.js
k6 run ./tests/prom100.js

sudo pkill bpftrace || true
killall harness || true

rm results.csv || true
cp baseline.csv results.csv
awk 'FNR>1' prometheus1.csv prometheus10.csv prometheus100.csv usdt1.csv usdt10.csv usdt100.csv  >> results.csv
