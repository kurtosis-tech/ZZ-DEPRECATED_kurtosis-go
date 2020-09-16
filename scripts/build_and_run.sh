set -euo pipefail
script_dirpath="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"

# ====================== CONSTANTS =======================================================
SUITE_IMAGE="kurtosistech/kurtosis-go-example"
KURTOSIS_CORE_CHANNEL="master"
INITIALIZER_IMAGE="kurtosistech/kurtosis-core_initializer:${KURTOSIS_CORE_CHANNEL}"
API_IMAGE="kurtosistech/kurtosis-core_api:${KURTOSIS_CORE_CHANNEL}"
PARALLELISM=3

BUILD_ACTION="build"
RUN_ACTION="run"
BOTH_ACTION="all"
HELP_ACTION="help"

# ====================== ARG PARSING =======================================================
show_help() {
    echo "${0} <action> <extra Docker args...>"
    echo ""
    echo "  Actions:"
    echo "    help    Displays this messages"
    echo "    build   Executes only the build step, skipping the run step"
    echo "    run     Executes only the run step, skipping the build step"
    echo "    all     Executes both build and run steps"
    echo ""
    echo "  Example:"
    echo "    ${0} all --env PARALLELISM=4"
    echo ""
}

action="${1:-}"
shift 1

do_build=true
do_run=true
case "${action}" in
    ${HELP_ACTION})
        show_help
        exit 0
        ;;
    ${BUILD_ACTION})
        do_build=true
        do_run=false
        ;;
    ${RUN_ACTION})
        do_build=false
        do_run=true
        ;;
    ${BOTH_ACTION})
        do_build=true
        do_run=true
        ;;
    *)
        echo "Error: First argument must be one of '${HELP_ACTION}', '${BUILD_ACTION}', '${RUN_ACTION}', or '${BOTH_ACTION}'" >&2
        exit 1
        ;;
esac

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
    docker build -t "${SUITE_IMAGE}:${docker_tag}" -f "${root_dirpath}/example_impl/Dockerfile" "${root_dirpath}"
fi

if "${do_run}"; then
    suite_execution_volume="go-example-suite_${docker_tag}_$(date +%s)"
    docker volume create "${suite_execution_volume}"

    # NOTE: spaces here will confuse Docker!
    custom_env_vars_json_flag="CUSTOM_ENV_VARS_JSON={\"GO_EXAMPLE_SERVICE_IMAGE\":\"nginxdemos/hello\"}"
    docker run \
        --mount "type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock" \
        --mount "type=volume,source=${suite_execution_volume},target=/suite-execution" \
        --env "${custom_env_vars_json_flag}" \
        --env "TEST_SUITE_IMAGE=${SUITE_IMAGE}:${docker_tag}" \
        --env "SUITE_EXECUTION_VOLUME=${suite_execution_volume}" \
        --env "KURTOSIS_API_IMAGE=${API_IMAGE}" \
        --env "PARALLELISM=${PARALLELISM}" \
        "${@}" \
        "${INITIALIZER_IMAGE}"
fi
