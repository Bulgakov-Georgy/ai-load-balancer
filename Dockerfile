FROM golang:1.23

#ENV CACHE_RESPONSE_EXPIRATION_TIME=5m \
#    HUGGINGFACE_KEY="" \
#    HUGGINGFACE_RPM=100 \
#    HUGGINGFACE_TIMEOUT=1m \
#    HUGGINGFACE_URL=https://api-inference.huggingface.co/models/ \
#    MODEL_TO_CLIENTS=model-to-clients.json \
#    OPENAI_KEY="" \
#    OPENAI_RPM=60 \
#    OPENAI_TIMEOUT=1m \
#    OPENAI_URL=https://api.openai.com/v1/chat/completions \
#    OPENROUTER_KEY="" \
#    OPENROUTER_RPM=50 \
#    OPENROUTER_TIMEOUT=1m \
#    OPENROUTER_URL=https://openrouter.ai/api/v1/chat/completions \
#    REDIS_URL=localhost:6379 \
#    REDIS_PASSWORD="" \
#    USER_RPM=100

WORKDIR /app

COPY . .

RUN go build -o ai-load-balancer ./cmd/main.go

EXPOSE 8080

CMD ["./ai-load-balancer"]
