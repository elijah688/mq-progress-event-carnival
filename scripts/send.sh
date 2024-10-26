#!/bin/bash

EXCHANGE_NAME="my_topic"
ROUTING_KEY="my.routing.key"

# Generate 3 UUIDs
uuid1=$(uuidgen)
uuid2=$(uuidgen)
uuid3=$(uuidgen)

# Initialize progress for each UUID
progress1=0
progress2=0
progress3=0

# Loop until all progress values reach 1.00
while true; do
    # Check if all UUIDs have completed progress
    if (($(echo "$progress1 > 1.00" | bc -l) && $(echo "$progress2 > 1.00" | bc -l) && $(echo "$progress3 > 1.00" | bc -l))); then
        echo "3 processes completed"
        break
    fi

    # Draw a random number between 0 and 2 to select a UUID
    random_num=$((RANDOM % 3))

    # Set the UUID and progress based on the random number
    case $random_num in
    0)
        id=$uuid1
        progress=$progress1
        ;;
    1)
        id=$uuid2
        progress=$progress2
        ;;
    2)
        id=$uuid3
        progress=$progress3
        ;;
    esac

    # Skip if this UUID has already completed
    if (($(echo "$progress > 1.00" | bc -l))); then
        continue
    fi

    payload_content=$(jq -n \
        --arg id "$id" \
        --arg name "Save as draft filtered locations massively" \
        --arg user "elijah.iliev@brambles.com@geoloc1" \
        --arg state "Running" \
        --arg startTime "$(date -u +"%Y-%m-%dT%H:%M:%S.%NZ")" \
        --arg finishedTime "0001-01-01T00:00:00Z" \
        --arg duration "12 secs" \
        --arg errorMessage "" \
        --argjson percentageComplete "$progress" \
        '{
            id: $id,
            name: $name,
            user: $user,
            state: $state,
            startTime: $startTime,
            finishedTime: $finishedTime,
            duration: $duration,
            errorMessage: $errorMessage,
            percentageComplete: $percentageComplete
        }' | jq .)

    json_payload=$(jq -n \
        --arg routing_key "$ROUTING_KEY" \
        --arg payload "$payload_content" \
        '{
            properties: {},
            routing_key: $routing_key,
            payload: $payload,
            payload_encoding: "string"
        }')

    curl -u guest:guest \
        -H "content-type:application/json" \
        -X POST -d "$json_payload" \
        http://localhost:15672/api/exchanges/%2f/$EXCHANGE_NAME/publish

    # Update progress for the selected UUID
    case $random_num in
    0) progress1=$(echo "$progress1 + 0.01" | bc) ;;
    1) progress2=$(echo "$progress2 + 0.01" | bc) ;;
    2) progress3=$(echo "$progress3 + 0.01" | bc) ;;
    esac

    # Wait 500 milliseconds before sending the next message
    # sleep 0.2
done
