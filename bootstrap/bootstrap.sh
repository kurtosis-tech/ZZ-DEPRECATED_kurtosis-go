# Deletes all the extraneous files, leaving a repo containing only the example impl and necessary infrastructure needed to write a testsuite

EXAMPLE_IMPL_DIRNAME="example_impl"
README_FILENAME="README.md"

set -euo pipefail
script_dirpath="$(cd "$(dirname "${0}")" && pwd)"
root_dirpath="$(dirname "${script_dirpath}")"

read -p "VERIFICATION: This will delete nearly all files in ${root_dirpath}, leaving only what's necessary for writing a new Kurtosis Go testsuite! Are you sure you want to proceed? (Ctrl-C to abort, ENTER to continue)"

read -p "FINAL VERIFICATION: you DO want to delete files like the .git dir to bootstrap a new testsuite, correct? (Ctrl-C to abort, ENTER to continue)"

read -p "New Go module name (e.g. github.com/my-org/my-repo): " new_module_name

find "${root_dirpath}" \
    ! -name bootstrap \
    ! -name "${EXAMPLE_IMPL_DIRNAME}" \
    ! -name go.mod \
    ! -name go.sum \
    ! -name scripts \
    -mindepth 1 \
    -maxdepth 1 \
    -exec rm -rf {} \;

cp "${script_dirpath}/README.md" "${root_dirpath}/"

# Replace module names (we need the "-i '' " argument because Mac sed requires it)
existing_module_name="$(grep "module" "${root_dirpath}/go.mod" | awk '{print $2}')"
sed -i '' "s,${existing_module_name},${new_module_name},g" go.mod
# We search for old_module_name/example_impl because we don't want the old_module_name/lib entries to get renamed
sed -i '' "s,${existing_module_name}/${EXAMPLE_IMPL_DIRNAME},${new_module_name}/${EXAMPLE_IMPL_DIRNAME},g" $(find "${root_dirpath}" -type f)

rm -rf "${script_dirpath}"
echo "Bootstrap complete; view the README.md in ${root_dirpath} for next steps"
