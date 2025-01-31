## How-to

The service and required environment can be started with [docker-compose.yaml](docker-compose.yaml). It also contains
[environment variables](#environment-variables), including API keys which you need to set

There is also [docker-compose-load-test.yaml](docker-compose-load-test.yaml) which runs pretty simple load test

## Components

### Redis

Redis is used for caching data (API responses) and storing requests per minute: users RPM and API RPM. Available at
`:6379`

### Prometheus

Prometheus service for collecting and storing metrics. Configured
with [prometheus.yml](monitoring/prometheus/prometheus.yml). Available at `:9090`

Has two metrics:
* requests_total - total number of requests to APIs and the service
* api_response_time - time from starting request to API to getting response

### Grafana

Used for visualising metrics from prometheus. Has a few simple dashboards defined
in [ai_load_balancer_dashboard.json](monitoring/grafana/dashboards/ai_load_balancer_dashboard.json). Configuration files
are located at [grafana](monitoring/grafana). Available at `:3000`, default credentials: `admin/admin`

### k6

Load testing tool. Deployed in the [docker-compose-load-test.yaml](docker-compose-load-test.yaml). Test itself is
defined in [stress-test.js](test/stress-test.js) and it is pretty simple. It just throws requests as up to 50 different
users. And it actually has only one meaningful check - it checks that when we hit APIs' rate limits there are no more
successful requests. Test result are available in the logs (though can be outputted to file, grafana, etc.)

### Mock server

Mock server for testing purposes. Mocks APIs' responses with random request processing time. Deployed in
the [docker-compose-load-test.yaml](docker-compose-load-test.yaml). Available at `:8081`

### The service itself

Available at `:8080`. Has three endpoints: `/metrics`, `/swagger/index.html` (swagger API doc) and `/generate`, which is
described in the swagger. Request body for `/generate` has two fields: prompt and model. Model can be either "gpt" or
"meta-llama", they are mapped according to [model-to-clients.json](model-to-clients.json) where for each model ("gpt",
"meta-llama") available API and its actual model name are defined

`/generate` request flow:

1. Check if the user hit its request rate limit
2. Validate request body
3. Check if the response for this model and prompt is in cache (return it if it was there)
4. Find a set of suitable APIs for the requested model (using [model-to-clients.json](model-to-clients.json))
5. Using round-robin approach try to find available API, comparing current requests-per-minute count and maximum RPM for
   the API
    1. There is also a sample code in the [balancer.go](internal/api/balancer.go) for choosing API based on the current
       RPM and API average response time
6. Make a request to API and return response

Sample cURL requests (formatted for cmd):

* `curl localhost:8080/generate --header "Content-Type: application/json" -d "{\"prompt\":\"your prompt\", \"model\":\"gpt\"}"`
* `curl localhost:8080/generate --header "Content-Type: application/json" -d "{\"prompt\":\"your prompt\", \"model\":\"meta-llama\"}"`

About `/generate` error handling. All errors are accounted for in the service, one way or another, but it is not
reflected in the response of `/generate`. There of course should be proper system with error codes and error messages
for different situations, but it is time-consuming and I assumed that it is fine as it is probably not in the aim of
this assignment


## Environment variables

| Name                           | Description                                                               |
|--------------------------------|---------------------------------------------------------------------------|
| CACHE_RESPONSE_EXPIRATION_TIME | Expiration time for cached API response (type `time.Duration`)            |
| HUGGINGFACE_KEY                | HuggingFace API key                                                       |
| HUGGINGFACE_RPM                | HuggingFace maximum request per minute                                    |
| HUGGINGFACE_TIMEOUT            | HuggingFace request timeout (type `time.Duration`)                        |
| HUGGINGFACE_URL                | HuggingFace URL                                                           |
| MODEL_TO_CLIENTS               | Path to JSON file with mappings of _service's model_ to API and its model |
| OPENAI_KEY                     | OpenAI API key                                                            |
| OPENAI_RPM                     | OpenAI maximum request per minute                                         |
| OPENAI_TIMEOUT                 | OpenAI request timeout (type `time.Duration`)                             |
| OPENAI_URL                     | OpenAI URL                                                                |
| OPENROUTER_KEY                 | OpenRouter API key                                                        |
| OPENROUTER_RPM                 | OpenRouter maximum request per minute                                     |
| OPENROUTER_TIMEOUT             | OpenRouter request timeout (type `time.Duration`)                         |
| OPENROUTER_URL                 | OpenRouter URL                                                            |
| REDIS_URL                      | Redis URL                                                                 |
| REDIS_PASSWORD                 | Redis password                                                            |
| USER_RPM                       | Maximum requests per minute for user                                      |
