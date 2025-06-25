#!/bin/sh

. build/ci-release-prep.sh

assets="$(
    for GOOS in $(find bin/ -type d -mindepth 1 -maxdepth 1 -exec basename '{}' ';'); do
        for GOARCH in $(find "bin/$GOOS" -type d -mindepth 1 -maxdepth 1 -exec basename '{}' ';'); do
            for asset in $(find "bin/$GOOS/$GOARCH" -type f -executable -mindepth 1 -maxdepth 1 -exec basename '{}' ';'); do
                echo "${asset}-${GOOS}-${GOARCH}"
            done
        done
    done
)"

release-cli create \
    --name "yum2npm v${tag_name}" \
    --tag-name "${CI_COMMIT_TAG}" \
    --ref "${CI_COMMIT_SHA}" \
    $(
        for asset in ${assets}; do
            echo -n --assets-link '{"name":"'"${asset}"'","url":"'"${assets_url_base}/${asset}"'"}' " "
        done
    )
