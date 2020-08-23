FROM golang:1.15-alpine AS builder
WORKDIR /build
# Copy and download dependencies using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o kurtosis-go-example example_impl/main.go

CMD set -euo pipefail && ./kurtosis-go-example \
    --test-names-filepath=${TEST_NAMES_FILEPATH} \
    --test=${TEST} \
    --kurtosis-api-ip=${KURTOSIS_API_IP} 2>&1 | tee ${LOG_FILEPATH}