set -x
set -euo pipefail

base_revision="${1:-}"

if [ -z "${base_revision}" ]; then
    echo "No base revision supplied" >&2
    exit 1
fi

if git diff --exit-code ${base_revision}...HEAD CHANGELOG.md; then
    echo "PR has no CHANGELOG entry. Please update the CHANGELOG!"
    return_code=1
else
    return_code=0
fi
exit "${return_code}"
