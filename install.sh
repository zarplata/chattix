#!/bin/bash

set -eou pipefail

CONST_CHAT_SLACK="slack"
CONST_CHAT_MATTERMOST="mattermost"

# DEFAULT VALUES
DEFAULT_CHAT="$CONST_CHAT_SLACK"
DEFAULT_CHAT_URL="https://slack.com/api/chat.postMessage"
DEFAULT_CHAT_API_TOKEN="supersecrettoken"

DEFAULT_ZABBIX_ADDRESS="localhost"

RED='\033[0;31m'
BRED='\033[1;31m'
YELLOW='\033[1;33m'
BOLD='\033[1m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

warn() {
    printf "${YELLOW}[WARN]${NC} $1\n" >&2
}

err() {
    printf "${RED}[ERRO]${NC} $1\n" >&2
}

cri() {
    printf "${BRED}[CRIT]${NC} $1\n" >&2
}

info() {
    printf "${CYAN}[INFO]${NC} $1\n" >&2
}

log() {
    printf "$1\n" >&2
}


set_chat() {
    local text 
    local selected_messenger

text=$(cat  <<-EOL
${CYAN}[1] Choose the type of messenger.${NC}
${BOLD}-> (1) ${DEFAULT_CHAT}${NC}
(2) mattermost

Write the number of selected messenger or
any key if you want to use default selection
EOL
)
    log "$text"

    IFS="$\n" read -r selected_messenger

    if [ "$selected_messenger" = "2" ]; then
        DEFAULT_CHAT=$CONST_CHAT_MATTERMOST
    fi
    
    log "Your choice is: ${BOLD}$DEFAULT_CHAT${NC}\n"
}

set_chat_url() {
    local text
    local address

text=$(cat  <<-EOL
${CYAN}[2] Write address of the ${DEFAULT_CHAT} messenger\n${NC}
EOL
)

    log "$text"

    if [ "$DEFAULT_CHAT" = "$CONST_CHAT_SLACK" ]; then
        text=$(cat  <<-EOL
For slack the address has already configured.
The address is: ${BOLD}${YELLOW}${DEFAULT_CHAT_URL}${NC}${NC}
And you dont't need to change it.
Howewer if you want change the behaivor of
the webhook script you may change it in 
result configuration file after the 
installation will have finished.
EOL
    )
    log "$text"
    fi

    if [ "$DEFAULT_CHAT" = "$CONST_CHAT_MATTERMOST" ]; then
        printf "Messenger address: " >&2
        IFS="$\n" read -r address
        printf "$address\n"

    fi


}

set_chat_token() {
    local text
    local token

text=$(cat  <<-EOL
${CYAN}[3] Set up access API token for ${DEFAULT_CHAT} messenger\n${NC}
EOL
)
    log "$text"

text=$(cat  <<-EOL
Write the access API token which provided by ${DEFAULT_CHAT}.
For ${CONST_CHAT_SLACK} you can find it on page ${YELLOW}https://api.slack.com/apps/APPID/oauth?${NC}
For ${CONST_CHAT_MATTERMOST} it's not necessary to provide the access token 
because messages will send through webhook mechanism.
EOL
)

    log "$text"

printf "Access token: " >&2

    IFS="$\n" read -r token
    if [[ -z "${token// }" ]]; then
        DEFAULT_CHAT_API_TOKEN=""
        return
    fi

    DEFAULT_CHAT_API_TOKEN=$token
}

set_chat
set_chat_url
set_chat_token
