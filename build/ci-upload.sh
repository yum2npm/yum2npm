#!/bin/sh

. build/ci-release-prep.sh

for GOOS in $(find bin/ -type d -mindepth 1 -maxdepth 1 -exec basename '{}' ';'); do
    for GOARCH in $(find "bin/$GOOS" -type d -mindepth 1 -maxdepth 1 -exec basename '{}' ';'); do
        for asset in $(find "bin/$GOOS/$GOARCH" -type f -executable -mindepth 1 -maxdepth 1 -exec basename '{}' ';'); do
            curl \
                --header "JOB-TOKEN: $CI_JOB_TOKEN" \
                --upload-file bin/${GOOS}/${GOARCH}/${asset} \
                "${assets_url_base}/${asset}-${GOOS}-${GOARCH}"
        done
    done
done
