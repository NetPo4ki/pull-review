import http from 'k6/http';
import { sleep, fail } from 'k6';

export const options = { vus: 5, duration: '30s' };
const base = 'http://localhost:8080';

export function setup() {
  for (let i = 0; i < 60; i++) {
    const r = http.get(`${base}/ready`);
    if (r.status === 200) return;
    sleep(1);
  }
  fail('API not ready');
}

export default function () {
  const suffix = `${__VU}`;
  const team = `lt-team-${suffix}-${__ITER}`;
  const u1 = `lt-u1-${suffix}`;
  const u2 = `lt-u2-${suffix}`;
  const u3 = `lt-u3-${suffix}`;

  http.get(`${base}/health`);

  http.post(`${base}/team/add`, JSON.stringify({
    team_name: team,
    members: [
      { user_id: u1, username: 'A', is_active: true },
      { user_id: u2, username: 'B', is_active: true },
      { user_id: u3, username: 'C', is_active: true },
    ],
  }), { headers: { 'Content-Type': 'application/json' } });

  const prId = `lt-pr-${suffix}-${Date.now()}`;
  http.post(`${base}/pullRequest/create`, JSON.stringify({
    pull_request_id: prId,
    pull_request_name: 'load',
    author_id: u1,
  }), { headers: { 'Content-Type': 'application/json' } });

  http.post(`${base}/pullRequest/merge`, JSON.stringify({
    pull_request_id: prId,
  }), { headers: { 'Content-Type': 'application/json' } });

  sleep(1);
}