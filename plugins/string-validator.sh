#!/bin/bash
if [ -z "${VALIDATION_PATTERN}" ]; then
  echo "Required variable 'VALIDATION_PATTERN' is not set."
  exit 1
fi

# indicator whether the script should fail with a non-zero exit code or just print a warning
FAIL_ON_MISMATCH="${FAIL_ON_MISMATCH:=true}"

case "${GIKS_HOOK_TYPE}" in
  "commit-msg")
    if [[ "${1}" =~ ${VALIDATION_PATTERN} ]]; then
      exit 0
    fi

    if [ "${FAIL_ON_MISMATCH}" == "false" ]; then
      echo "WARNING: provided string '${1}' does not match required pattern '${VALIDATION_PATTERN}'."
      exit 0
    fi

    echo "ERROR: provided string '${1}' does not match required pattern '${VALIDATION_PATTERN}'."
    exit 1
  ;;
  *)
    echo "WARNING: string validator plugin does support commit hook type '${GIKS_HOOK_TYPE}'"
    ;;
esac