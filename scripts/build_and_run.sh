set -euo pipefail
script_dirpath="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"

# ====================== CONSTANTS =======================================================
# TODO Extract these constants out into their own file so the bootstrap can easily replace it
# The name of to give the Docker image containing the testsuite
KURTOSIS_DOCKERHUB_ORG="kurtosistech"
SUITE_IMAGE="${KURTOSIS_DOCKERHUB_ORG}/kurtosis-go-example"

GO_MODULE="github.com/kurtosis-tech/kurtosis-go"
LIB_CORE_API_RELATIVE_DIRPATH="lib_core/api"
LIB_CORE_API_GENERATED_CODE_DIRNAME="generated"

BUILD_ACTION="build"
RUN_ACTION="run"
BOTH_ACTION="all"
HELP_ACTION="help"

# ====================== ARG PARSING =======================================================
show_help() {
    echo "${0} <action> [<extra kurtosis.sh script args...>]"
    echo ""
    echo "  This script will optionally a) build your Kurtosis testsuite into a Docker image and/or b) run it via a call to the kurtosis.sh script"
    echo ""
    echo "  To select behaviour, choose from the following actions:"
    echo ""
    echo "    help    Displays this messages"
    echo "    build   Executes only the build step, skipping the run step"
    echo "    run     Executes only the run step, skipping the build step"
    echo "    all     Executes both build and run steps"
    echo ""
    echo "  To see the args the kurtosis.sh script accepts for the 'run' phase, call '$(basename ${0}) all --help'"
    echo ""
}

if [ "${#}" -eq 0 ]; then
    show_help
    exit 1     # Exit with error code so we dont't get spurious CI passes
fi

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
    if ! [ -f "${root_dirpath}"/.dockerignore ]; then
        echo "Error: No .dockerignore file found in root; this is required so Docker caching works properly" >&2
        exit 1
    fi

    echo "Generating lib core files from protobuf files..."
    lib_core_api_dirpath="${root_dirpath}/${LIB_CORE_API_RELATIVE_DIRPATH}"
    generated_api_code_dirpath="${lib_core_api_dirpath}/${LIB_CORE_API_GENERATED_CODE_DIRNAME}"
    if [ "${generated_api_code_dirpath}/" != "/" ]; then
        if ! find ${generated_api_code_dirpath} -name '*.go' -delete; then
            echo "Error: An error occurred removing the protobuf-generated files" >&2
            exit 1
        fi
    else
        echo "Error: lib core generated API code dirpath must not be empty!" >&2
        exit 1
    fi
    for protobuf_file in ${lib_core_api_dirpath}/*.proto; do
        # TODO clean this thing up
        if ! protoc -I="${lib_core_api_dirpath}" \
                --go_out=plugins=grpc:"${generated_api_code_dirpath}" \
                --go_opt=module=${GO_MODULE}/${LIB_CORE_API_RELATIVE_DIRPATH}/${LIB_CORE_API_GENERATED_CODE_DIRNAME} \
                "${protobuf_file}"; then
            echo "Error: An error occurred generating lib core files from protobuf file: ${protobuf_file}" >&2
            exit 1
        fi
    done
    echo "Successfully generated lib core files from protobuf files"

    echo "Running unit tests..."
    # TODO Extract this go-specific logic out into a separate script so we can copy/paste the build_and_run.sh between various languages
    if ! go test "${root_dirpath}/..."; then
        echo "Tests failed!"
        exit 1
    fi
    echo "Tests succeeded"

    echo "Building ${SUITE_IMAGE} Docker image..."
    docker build -t "${SUITE_IMAGE}:${docker_tag}" -f "${root_dirpath}/testsuite/Dockerfile" "${root_dirpath}"
fi

if "${do_run}"; then
    # ======================================= Custom Docker environment variables ========================================================
    # NOTE: Replace these with whatever custom properties your service needs
    api_service_image="${KURTOSIS_DOCKERHUB_ORG}/example-microservices_api"
    datastore_service_image="${KURTOSIS_DOCKERHUB_ORG}/example-microservices_datastore"
    # Docker only allows you to have spaces in the variable if you escape them or use a Docker env file
    custom_env_vars_json='{
        "API_SERVICE_IMAGE": "'${api_service_image}'",
        "DATASTORE_SERVICE_IMAGE": "'${datastore_service_image}'"
    }'
    # ====================================== End custom Docker environment variables =====================================================

    bash "${script_dirpath}/kurtosis.sh" --custom-env-vars "${custom_env_vars_json}" "${@}" "${SUITE_IMAGE}:${docker_tag}"
fi
