FROM golang:1.15-alpine AS builder
WORKDIR /build
# Copy and download dependencies using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o testsuite example_impl/main.go

# TODO Get rid of tee/LOG_FILEPATH in favor of using a Docker logging driver in the initializer
# Note how we use a variable custom to this particular test suite, GO_EXAMPLE_SERVICE_IMAGE
CMD set -euo pipefail && ./testsuite \
    --metadata-filepath=${METADATA_FILEPATH} \
    --service-image=${GO_EXAMPLE_SERVICE_IMAGE} \
    --test=${TEST} \
    --log-level=${LOG_LEVEL} \
    --services-relative-dirpath=${SERVICES_RELATIVE_DIRPATH} \
    --kurtosis-api-ip=${KURTOSIS_API_IP} 2>&1 | tee ${LOG_FILEPATH}
