set -euo pipefail

script_dirpath="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"
root_dirpath="$(dirname "${script_dirpath}")"

echo "Running unit tests..."
if ! go test "${root_dirpath}/..."; then
    echo "Tests failed!"
    exit 1
fi
echo "Tests succeeded"

