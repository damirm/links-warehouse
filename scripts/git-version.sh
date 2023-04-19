#!/bin/bash

CURRENT_DATE=${CURRENT_DATE:-$(date +%Y%m%d)}
CURRENT_REVISION=${CURRENT_REVISION:-$(git rev-parse --short HEAD)}
TODAY_COMMITS_COUNTER=$(git rev-list --count --since=yesterday --before=today HEAD)

echo -n "${CURRENT_DATE}-${TODAY_COMMITS_COUNTER}-${CURRENT_REVISION}"

