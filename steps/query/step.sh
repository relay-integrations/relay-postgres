#!/bin/bash
set -euo pipefail

NI="${NI:-ni}"
CLOUD_SQL="${CLOUD_SQL:-cloud_sql_proxy}"

GOOGLE=$(ni get -p {.google})
if [ -n "${GOOGLE}" ]; then
  $NI gcp config -d "/workspace/.gcp"
fi

INSTANCE=$($NI get -p {.instance})
if [ -n "${INSTANCE}" ]; then
  ${CLOUD_SQL} -instances=${INSTANCE}=tcp:5432 -credential_file=/workspace/.gcp/credentials.json -term_timeout=60m &
fi

STATEMENT=$(ni get -p {.statement})
if [ -n "${STATEMENT}" ]; then
  relay-postgres-query
fi

if [ -n "${INSTANCE}" ]; then
  pkill -HUP ${CLOUD_SQL} &
fi
