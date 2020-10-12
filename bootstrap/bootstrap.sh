# Deletes all the extraneous files, leaving a repo containing only the example impl and necessary infrastructure needed to write a testsuite

set -euo pipefail
script_dirpath="$(cd "$(dirname "${0}")" && pwd)"
root_dirpath="$(dirname "${script_dirpath}")"

read -p "VERIFICATION: This will delete nearly all files in ${root_dirpath}, leaving only what's necessary for writing a new Kurtosis Go testsuite! Are you sure you want to proceed? (Ctrl-C to abort, ENTER to continue)"

read -p "FINAL VERIFICATION: you DO want to delete files like the .git dir to bootstrap a new testsuite, correct? (Ctrl-C to abort, ENTER to continue)"



find "${root_dirpath}" \
    ! -name bootstrap \
    ! -name example_impl \
    ! -name go.mod \
    ! -name go.sum \
    ! -name scripts \
    -mindepth 1 \
    -maxdepth 1 \
    -exec echo {} \;
