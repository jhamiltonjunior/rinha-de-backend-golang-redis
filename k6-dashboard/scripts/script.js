import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 100,
  duration: '2m',
};

export default function () {
  const res = http.get('https://test-api.k6.io/public/crocodiles/1/');
  check(res, { 'status is 200': (r) => r.status === 200 });
  sleep(1);
}