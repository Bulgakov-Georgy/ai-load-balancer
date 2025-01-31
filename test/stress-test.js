import http from 'k6/http';
import {check, sleep} from 'k6';
import {Counter} from 'k6/metrics';

const httpReqSucceededCounter = new Counter('http_req_succeeded');

export let options = {
    stages: [
        { duration: '10s', target: 50 },  // Ramp up to 50 users
        { duration: '30s', target: 50 },  // Hold 50 users
        { duration: '10s', target: 0 },   // Ramp down
    ],
    thresholds: {
        http_req_duration: ['p(95)<5000'], // 95% of requests < 5000ms
        http_req_succeeded: ['count<=210'],   // 210 is a sum of rate limits per minute for the test
    },
};

// Generate random IPs to bypass per-user rate limiting
function getUserIp() {
    return `${Math.floor(Math.random() * 255)}.${Math.floor(Math.random() * 255)}.${Math.floor(Math.random() * 255)}.${Math.floor(Math.random() * 255)}`;
}

// Get random model for the request
function getRandomModel() {
    return Math.random() > 0.7 ? "gpt" : "meta-llama";
}

export default function () {
    let url = 'http://ai-load-balancer:8080/generate';
    let userIp = getUserIp();
    let model = getRandomModel()

    let payload = JSON.stringify({ prompt: (Math.random() * 10e5).toString(), model: model });
    let params = {
        headers: {
            'Content-Type': 'application/json',
            'X-Forwarded-For': `${userIp}`,
        },
    };

    let res = http.post(url, payload, params);

    if (res.status === 200) {
        httpReqSucceededCounter.add(1);
    }

    check(res, {
        'is status 200': (r) => r.status === 200,
        'is rate limited (429)': (r) => r.status === 429,
        'response time < 5000ms': (r) => r.timings.duration < 5000,
    });

    sleep(1);
}
