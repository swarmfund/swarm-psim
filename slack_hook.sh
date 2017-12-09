#!/bin/bash

LINK="<https://gitlab.com/swarmfund/psim/pipelines|PSIM>"
COMMIT_SHA=$(git rev-parse --short HEAD)
COMMIT_MESSAGE=$(git log -1 --pretty=%B)

start_build_mod="
\"fallback\":\"$LINK Build Started!\",
\"pretext\":\"$LINK Build Started!\",
\"color\":\"#0000ff\"
"
success_mod="
\"fallback\":\"$LINK Build Success!\",
\"pretext\":\"$LINK Build Success!\",
\"color\":\"#008000\"
"

failure_mod="
\"fallback\":\"$LINK Build Failed!\",
\"pretext\":\"$LINK Build Failed!\",
\"color\":\"#D00000\"
"
mod=$start_build_mod

case "$1" in
    "start-build")
        mod=$start_build_mod
        ;;
    "success")
        mod=$success_mod
        ;;
    "failed")
        mod=$failure_mod
        ;;
    *)
    echo "Unknown hook type!"
    exit 1
esac

msg="{
    \"channel\":\"#reports_ci_builds\",
    \"username\":\"Build Report\",
    \"icon_emoji\": \":squirrel:\",
    \"attachments\":[
      {
         $mod,
         \"fields\":[
            {
               \"title\":\"Commit:\",
               \"value\":\"$COMMIT_SHA: $COMMIT_MESSAGE\",
               \"short\":false
            }
         ]
      }
   ]
}"

curl -X POST --data-urlencode "payload=$msg" https://hooks.slack.com/services/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