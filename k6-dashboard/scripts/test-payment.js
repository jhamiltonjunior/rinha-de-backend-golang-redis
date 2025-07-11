import { uuidv4 } from "https://jslib.k6.io/k6-utils/1.4.0/index.js";
import { check, sleep } from "k6";
import http from "k6/http";

export const options = {
  stages: [
    { duration: "30s", target: 10 },   
    { duration: "1m", target: 10 },    
    { duration: "30s", target: 0 },    
  ],
  thresholds: {
    http_req_duration: ["p(95) < 500"], 
    http_req_failed: ["rate < 0.1"],    
  },
};

export default function () {
  
  const correlationId = uuidv4();
  const amount = Math.round((Math.random() * 100 + 1) * 100) / 100; 

  const payload = {
    correlationId: correlationId,
    amount: amount
  };

  const params = {
    headers: {
      "Content-Type": "application/json",
    },
  };

  const response = http.post("http://nginx/payments", JSON.stringify(payload), params);

  check(response, {
    "status is 200-299": (r) => r.status >= 200 && r.status < 300,
    "response time < 500ms": (r) => r.timings.duration < 500,
    "has valid response": (r) => r.body && r.body.length > 0,
  });

  
  console.log(`Request sent - CorrelationId: ${correlationId}, Amount: ${amount}, Status: ${response.status}`);

  
  sleep(1);
}


export function setup() {
  console.log("Starting load test for localhost:9999");
  console.log("Testing POST requests with random correlationId and amount");
}


export function teardown(data) {
  console.log("Load test completed");
}
