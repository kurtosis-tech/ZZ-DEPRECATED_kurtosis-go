FROM golang:1.15-alpine AS builder
WORKDIR /build
# Copy and download dependencies using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o testsuite.bin testsuite/main.go

CMD ./testsuite.bin \
    `# ------------- Custom variable particular to this test suite -----------------` \
    --api-service-image=${API_SERVICE_IMAGE} \
    --datastore-service-image=${DATASTORE_SERVICE_IMAGE} \
    `# ------------------Kurtosis-specific environment variables -------------------` \
    `# See: http://docs.kurtosistech.com/testsuite-details.html
    --kurtosis-api-ip=${KURTOSIS_API_IP} \
    --log-level=${LOG_LEVEL} \
    --metadata-filepath=${METADATA_FILEPATH} \
    --services-relative-dirpath=${SERVICES_RELATIVE_DIRPATH} \
    --test=${TEST}
