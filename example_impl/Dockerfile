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

CMD ./testsuite \
    `# ------------- Custom variable particular to this test suite -----------------` \
    --service-image=${EXAMPLE_SERVICE_IMAGE} \
    `# ------------------Kurtosis-specific environment variables -------------------` \
    `# See: https://github.com/kurtosis-tech/kurtosis-docs/blob/master/testsuite-details.md` \
    --kurtosis-api-ip=${KURTOSIS_API_IP} \
    --log-level=${LOG_LEVEL} \
    --metadata-filepath=${METADATA_FILEPATH} \
    --services-relative-dirpath=${SERVICES_RELATIVE_DIRPATH} \
    --test=${TEST}
