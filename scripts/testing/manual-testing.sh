#!/bin/bash

if [ $# -eq 0 ]; then
    echo "Usage: $0 <messages_count>"
    exit 1
fi

COUNT=$1

MODES=("[\"big\",\"small\"]" "[\"small\"]" "[\"big\"]")

for ((i=1; i<=COUNT; i++)); do
    ID=$(uuidgen)
    TROPHIES=$((RANDOM % 2301 + 200))
    PREMIUM=$([ $((RANDOM % 2)) -eq 1 ] && echo "true" || echo "false")
    MODE_INDEX=$((RANDOM % 3))
    GAME_MODE=${MODES[$MODE_INDEX]}

    MESSAGE="{\"id\":\"$ID\",\"trophies\":$TROPHIES,\"is_premium\":$PREMIUM,\"game_mode\":$GAME_MODE}"

    nats pub "create.player" "$MESSAGE"

    echo "Send: $MESSAGE"
done