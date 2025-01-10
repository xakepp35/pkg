# secrets
TG_CHAT_ID=985959967
TG_BOT_TOKEN=7694863935:AAEMK4MmLacTzBEc6-UVdg91M8B2JmyNvT0

#github actions params
PROJECT="okto-local-server" ## ${{needs.build.result}}
BUILD_RESULT="success" ## ${{needs.build.result}}
ACTOR="xakepp35" ## 
DEPLOY_SERVER="mars.stage.okto.ru" # ${{inputs.server}}
IS_BUILD="true" # ${{inputs.build}}
IS_DEPLOY="true" # ${{inputs.deploy}}
START_TIME=$(date +%s) # $(date -d "${{ github.run_started_at }}" +%s)
LOGS=""
GIT_BRANCH="branch"

#TODO pass everything is needed, 
WORKFLOW_URL= "https://github.com/"#${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }}



STATUS="✅ Success"
if [[ "${{ BUILD_RESULT }}" != "success" ]]; then
	STATUS="❌ Failure"
	# LOGS=$(tail -n 20 $GITHUB_WORKSPACE/_temp/logs || echo "No logs available.")
fi

#START_TIME=$(date -d "${{ github.run_started_at }}" +%s)
END_TIME=$(date +%s)
DURATION=$(( END_TIME - START_TIME ))
MINUTES=$(( DURATION / 60 ))
SECONDS=$(( DURATION % 60 ))
TEXT="🚀 *Build and Deploy* 🚀\n"
TEXT="${TEXT}\n"
if [[ "${{ BUILD_RESULT }}" != "success" ]]; then
TEXT="${TEXT}- *Build:* ${IS_BUILD}\n"
fi

TEXT="${TEXT}\**Status*: ${STATUS}\n"
TEXT="${TEXT}\**Project*: ${PROJECT}\n"
TEXT="${TEXT}*Actor:* ${ACTOR}\n"
TEXT="${TEXT}- *Branch:* ${GIT_BRANCH}\n"
if BUILD then
	TEXT="${TEXT}- *Build time:* ${GIT_BRANCH}\n"
if DEPLOY then
	TEXT="${TEXT}- *Server:* ${DEPLOY_SERVER}\n"
fi
TEXT="${TEXT}\n"
TEXT="${TEXT}- *Docker tag:* ${TODO}\n"
TEXT="${TEXT}- *Cost:* ${MINUTES} min ${SECONDS} sec\n"
TEXT="${TEXT}- *Run:* [Workflow Link](${WORKFLOW_URL})\n"


if [[ "$BUILD_RESULT" != "success" ]]; then
  TEXT="${TEXT}\nLogs:\n\`\`\`\n${LOGS}\n\`\`\`"
fi

TEXT=$(printf "%b" "$TEXT")

echo $TEXT
curl -X POST "https://api.telegram.org/bot${TG_BOT_TOKEN}/sendMessage" \
  -d parse_mode="Markdown" \
  -d chat_id="${TG_CHAT_ID}" \
  -d text="$TEXT"
