# Deletes all the extraneous files, leaving a repo containing only the example impl and necessary infrastructure needed to write a testsuite

EXAMPLE_IMPL_DIRNAME="example_impl"
README_FILENAME="README.md"

# Constants 
GO_MOD_FILENAME="go.mod"
GO_MOD_MODULE_KEYWORD="module "  # The key we'll look for when replacing the module name in go.mod
BUILDSCRIPT_FILENAME="build_and_run.sh"
DOCKER_IMAGE_VAR_KEYWORD="SUITE_IMAGE=" # The variable we'll look for in the Docker file for replacing the Docker image name

set -euo pipefail
script_dirpath="$(cd "$(dirname "${0}")" && pwd)"
root_dirpath="$(dirname "${script_dirpath}")"
buildscript_filepath="${root_dirpath}/scripts/${BUILDSCRIPT_FILENAME}"
go_mod_filepath="${root_dirpath}/${GO_MOD_FILENAME}"

# ============== Validation =================================================================
# Validation, to save us in case someone changes stuff in the future
if [ "$(grep "${GO_MOD_MODULE_KEYWORD}" "${go_mod_filepath}" | wc -l)" -ne 1 ]; then
    echo "Validation failed: Could not find exactly one line in ${GO_MOD_FILENAME} with keyword '${GO_MOD_MODULE_KEYWORD}' for use when replacing with the user's module name" >&2
    exit 1
fi
if [ "$(grep "^${DOCKER_IMAGE_VAR_KEYWORD}" "${buildscript_filepath}" | wc -l)" -ne 1 ]; then
    echo "Validation failed: Could not find exactly one line in ${buildscript_filepath} starting with keyword '${DOCKER_IMAGE_VAR_KEYWORD}' for use when replacing with the user's Docker image name" >&2
    exit 1
fi

# ============== Inputs & Verification =================================================================
read -p "VERIFICATION: This will clear the state of this repo, leaving a fresh, bootstrapped testsuite. Proceed? (Ctrl-C to abort, ENTER to continue)"
new_module_name=""
while [ -z "${new_module_name}" ]; do
    read -p "Go module name to give this repo (e.g. github.com/your-org/your-repo): " new_module_name
done
docker_image_name=""
while [ -z "${docker_image_name}" ]; do
    echo "Your testsuite will be encapsulated in a Docker image, which must conform to the Docker image naming rules at https://docs.docker.com/engine/reference/commandline/tag/#extended-description"
    read -p "Name to give the Docker image of this testsuite (e.g. your-docker-repo/your-image-name): " docker_image_name
done


# ============== Main Code =================================================================
find "${root_dirpath}" \
    ! -name bootstrap \
    ! -name "${EXAMPLE_IMPL_DIRNAME}" \
    ! -name "${GO_MOD_FILENAME}" \
    ! -name go.sum \
    ! -name scripts \
    -mindepth 1 \
    -maxdepth 1 \
    -exec rm -rf {} \;

cp "${script_dirpath}/README.md" "${root_dirpath}/"

# Replace module names in code (we need the "-i '' " argument because Mac sed requires it)
existing_module_name="$(grep "module" "${go_mod_filepath}" | awk '{print $2}')"
sed -i '' "s,${existing_module_name},${new_module_name},g" ${go_mod_filepath}
# We search for old_module_name/example_impl because we don't want the old_module_name/lib entries to get renamed
sed -i '' "s,${existing_module_name}/${EXAMPLE_IMPL_DIRNAME},${new_module_name}/${EXAMPLE_IMPL_DIRNAME},g" $(find "${root_dirpath}" -type f)

# Replace Docker image name in buildscript
sed -i '' "s,^${DOCKER_IMAGE_VAR_KEYWORD}.*,${DOCKER_IMAGE_VAR_KEYWORD}\"${docker_image_name}\"," "${buildscript_filepath}"

rm -rf "${script_dirpath}"
echo "Bootstrap complete; view the README.md in ${root_dirpath} for next steps"
