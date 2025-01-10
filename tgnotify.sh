# secrets
TG_CHAT_ID=985959967
TG_BOT_TOKEN=7694863935:AAEMK4MmLacTzBEc6-UVdg91M8B2JmyNvT0

#github actions params
PROJECT="test-project" ## ${{needs.build.result}}
BUILD_RESULT="failure" ## ${{needs.build.result}}
DEPLOY_SERVER="google.com" # ${{inputs.server}}
IS_BUILD="true" # ${{inputs.build}}
IS_DEPLOY="true" # ${{inputs.deploy}}
START_TIME=$(date +%s) # $(date -d "${{ github.run_started_at }}" +%s)
LOGS="hello\nworld"
GIT_BRANCH="main"
GIT_SERVER="https://github.com" # ${{github.server_url}}
GIT_REPO="xakepp35/pkg" # ${{github.repository}}
GIT_RUN="123" # ${{github.run_id}}
GIT_ACTOR="xakepp35" ## ${{github.actor}}
DOCKER_TAG="qwertmax/okto-local-server:feature-swagger-12345678"
NOTIFICATION_HEADER="🚀 *Build and Deploy* 🚀"

GIT_REPO_URL="${GIT_SERVER}/${GIT_REPO}/"
GIT_ACTOR_URL="${GIT_SERVER}/${GIT_ACTOR}/"
GIT_RUN_URL="${GIT_REPO_URL}actions/runs/${GIT_RUN}"
GIT_BRANCH_URL="${GIT_REPO_URL}tree/${GIT_BRANCH}"

STATUS="✅ Success"
if [[ "${BUILD_RESULT}" != "success" ]]; then
	STATUS="❌ Failure"
	# LOGS=$(tail -n 20 $GITHUB_WORKSPACE/_temp/logs || echo "No logs available.")
fi

#START_TIME=$(date -d "${{ github.run_started_at }}" +%s)
END_TIME=$(date +%s)
DURATION=$(( END_TIME - START_TIME ))
MINUTES=$(( DURATION / 60 ))
SECONDS=$(( DURATION % 60 ))
TXT="$NOTIFICATION_HEADER\n\n"
# if [[ "$BUILD_RESULT" != "success" ]]; then
#   TXT="${TXT}- *Build:* $IS_BUILD\n"
# fi
TXT="${TXT}*Status*: ${STATUS}\n\n"
TXT="${TXT}*Repo:* [${GIT_REPO}](${GIT_REPO_URL})\n"
TXT="${TXT}*Actor:* [${GIT_ACTOR}](${GIT_ACTOR_URL})\n"
TXT="${TXT}*Branch:* [${GIT_BRANCH}](${GIT_BRANCH_URL})\n"
TXT="${TXT}*Run:* [$GIT_RUN]($GIT_RUN_URL)\n"

TXT="${TXT}*Project*: ${PROJECT}\n"
if [[ "${IS_BUILD}" == "true" ]]; then
	TXT="${TXT}*Build time:* $MINUTES min $SECONDS sec\n"
fi
if [[ "${IS_DEPLOY}" = "true" ]]; then
	TXT="${TXT}*Deploy Server:* ${DEPLOY_SERVER}\n"
fi
TXT="${TXT}*Docker tag:* $DOCKER_TAG\n"

if [[ "$BUILD_RESULT" != "success" ]]; then
  TXT="${TXT}\n*Logs:*\n\`\`\`\n${LOGS}\n\`\`\`"
fi
TXT=$(printf "%b" "$TXT")
echo $TXT
curl -X POST "https://api.telegram.org/bot${TG_BOT_TOKEN}/sendMessage" \
  -d parse_mode="Markdown" \
  -d chat_id="${TG_CHAT_ID}" \
  -d disable_web_page_preview=true \
  -d text="$TXT"
