#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
ROOT_DIR="${SCRIPT_DIR}/../"

pushd ${ROOT_DIR}
DHC_SECRETS_PATH=${DHC_SECRETS_PATH:-"/opt/secrets/t1"}
SECRETS_PATH="$DHC_SECRETS_PATH/../cluster-api-provider-openstack"
for f in $(pushd "$SECRETS_PATH" >/dev/null; find . -type d -path '*/\.*' -prune -o -not -name '.*' -type f -print; popd >/dev/null); do
    mkdir -p $(dirname ${f})
    cp "$SECRETS_PATH/$f" "$f"
done
popd
