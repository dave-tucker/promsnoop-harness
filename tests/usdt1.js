import { sleep } from 'k6';
import http from 'k6/http';
import { summarize } from './summary.js';

export const options = {
  vus: 1,
  iterations: 100,
};

export default function () {
  http.get('http://127.0.0.1:8080/usdt/1');
  sleep(1);
}

export function handleSummary(data) {
  return summarize("usdt1", data)
}
