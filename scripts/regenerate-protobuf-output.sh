set -euo pipefail
script_dirpath="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"
root_dirpath="$(dirname "${script_dirpath}")"

# ================================ CONSTANTS =======================================================
GO_MODULE="github.com/kurtosis-tech/kurtosis-go"

LIB_DIRNAME="lib"
API_DIRNAME="core_api"
BINDINGS_DIRNAME="bindings"
API_BINDINGS_GO_PKG="${GO_MODULE}/${LIB_DIRNAME}/${BINDINGS_DIRNAME}"

# =============================== MAIN LOGIC =======================================================
api_dirpath="${root_dirpath}/${LIB_DIRNAME}/${API_DIRNAME}"
input_dirpath="${api_dirpath}"
output_dirpath="${api_dirpath}/${BINDINGS_DIRNAME}"

# TODO Upgrade this to remove all generated code
if [ "${output_dirpath}/" != "/" ]; then
    if ! find ${output_dirpath} -name '*.go' -delete; then
        echo "Error: An error occurred removing the existing protobuf-generated code" >&2
        exit 1
    fi
else
    echo "Error: output dirpath must not be empty!" >&2
    exit 1
fi

for protobuf_filepath in $(find "${input_dirpath}" -name "*.proto"); do
    protobuf_filename="$(basename "${protobuf_filepath}")"

    # NOTE: When multiple people start developing on this, we won't be able to rely on using the user's local protoc because they might differ. We'll need to standardize by:
    #  1) Using protoc inside the API container Dockerfile to generate the output Go files (standardizes the output files for Docker)
    #  2) Using the user's protoc to generate the output Go files on the local machine, so their IDEs will work
    #  3) Tying the protoc inside the Dockerfile and the protoc on the user's machine together using a protoc version check
    #  4) Adding the locally-generated Go output files to .gitignore
    #  5) Adding the locally-generated Go output files to .dockerignore (since they'll get generated inside Docker)
    if ! protoc \
            -I="${input_dirpath}" \
            --go_out="plugins=grpc:${output_dirpath}" \
            `# Rather than specify the go_package in source code (which means all consumers of these protobufs would get it),` \
            `#  we specify the go_package here per https://developers.google.com/protocol-buffers/docs/reference/go-generated` \
            `# See also: https://github.com/golang/protobuf/issues/1272` \
            --go_opt="M${protobuf_filename}=${API_BINDINGS_GO_PKG};$(basename "${API_BINDINGS_GO_PKG}")" \
            "${protobuf_filepath}"; then
        echo "Error: An error occurred generating lib core files from protobuf file: ${protobuf_filepath}" >&2
        exit 1
    fi
done
