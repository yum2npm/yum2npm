#!/bin/sh

tag_name="$(echo "${CI_COMMIT_TITLE}" | cut -d" " -f2 | tr -d ']')"
tag_message="$(echo "${CI_COMMIT_TITLE}" | cut -d" " -f3-)"
assets_url_base="${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/packages/generic/yum2npm/${tag_name}"
