export function summarize(name, data) {
    const avg_latency = data.metrics['http_req_duration{expected_response:true}'].values.avg * 1000;
    const p90_latency = data.metrics['http_req_duration{expected_response:true}'].values['p(90)'] * 1000;
    const p95_latency = data.metrics['http_req_duration{expected_response:true}'].values['p(95)'] * 1000;
    const body = `Test, Average (µs), P90 (µs), P95 (µs)\n${name}, ${avg_latency.toFixed(2)}, ${p90_latency.toFixed(2)}, ${p95_latency.toFixed(2)}\n`;
    const filename = `${name}.csv`;
    return {
        [filename]: body,
    };
}
