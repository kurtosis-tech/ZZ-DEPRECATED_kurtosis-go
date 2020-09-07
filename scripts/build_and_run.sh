set -euo pipefail
script_dirpath="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"

# ====================== CONSTANTS =======================================================
DOCKER_ORG="kurtosistech"
EXAMPLE_IMAGE="kurtosis-go-example"
INITIALIZER_IMAGE="kurtosistech/kurtosis-core_initializer:develop"  # Parameterize the tag?

# ====================== ARG PARSING =======================================================
show_help() {
    echo "${0}:"
    echo "  -h      Displays this message"
    echo "  -b      Executes only the build step, skipping the run step"
    echo "  -r      Executes only the run step, skipping the build step"
    echo "  -d      Extra args to pass to 'docker run' (e.g. '--env MYVAR=somevalue')"
}

do_build=true
do_run=true
extra_docker_args=""
while getopts "brd:" opt; do
    case "${opt}" in
        h)
            show_help
            exit 0
            ;;
        b)
            do_run=false
            ;;
        r)
            do_build=false
            ;;
        d)
            extra_docker_args="${OPTARG}"
            ;;
    esac
done

# ====================== MAIN LOGIC =======================================================
git_branch="$(git rev-parse --abbrev-ref HEAD)"
docker_tag="$(echo "${git_branch}" | sed 's,[/:],_,g')"

root_dirpath="$(dirname "${script_dirpath}")"
if "${do_build}"; then
    echo "Running unit tests..."
    if ! go test "${root_dirpath}/..."; then
        echo "Tests failed!"
        exit 1
    else
        echo "Tests succeeded"
    fi

    echo "Building example Go implementation image..."
    docker build -t "${DOCKER_ORG}/${EXAMPLE_IMAGE}:${docker_tag}" -f "${root_dirpath}/example_impl/Dockerfile" "${root_dirpath}"
fi

if "${do_run}"; then
    go_suite_execution_volume="go-example-suite_${docker_tag}_$(date +%s)"
    docker volume create "${go_suite_execution_volume}"
    docker run \
        --mount "type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock" \
        --mount "type=volume,source=${go_suite_execution_volume},target=/suite-execution" \
        --env 'CUSTOM_ENV_VARS_JSON={"GO_EXAMPLE_SERVICE_IMAGE":"nginxdemos/hello"}' \
        --env "TEST_SUITE_IMAGE=${DOCKER_ORG}/${EXAMPLE_IMAGE}:${docker_tag}" \
        --env "SUITE_EXECUTION_VOLUME=${go_suite_execution_volume}" \
        ${extra_docker_args} \
        "${INITIALIZER_IMAGE}"
fi
