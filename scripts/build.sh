set -euo pipefail

DOCKER_ORG="kurtosistech"
EXAMPLE_IMAGE="kurtosis-go-example"

script_dirpath="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"
root_dirpath="$(dirname "${script_dirpath}")"

echo "Running unit tests..."
if ! go test "${root_dirpath}/..."; then
    echo "Tests failed!"
    exit 1
else
    echo "Tests succeeded"
fi

echo "Building example Go implementation image..."
docker build -t "${DOCKER_ORG}/${EXAMPLE_IMAGE}" -f "${root_dirpath}/example_impl/Dockerfile" "${root_dirpath}"
