#!/bin/bash
echo "hook-type: ${GIKS_HOOK_TYPE}"
if [ -z "${VALIDATION_PATTERN}" ]; then
  echo "Required variable 'VALIDATION_PATTERN' is not set."
  exit 1
fi

echo "Validating with pattern '${VALIDATION_PATTERN}'"