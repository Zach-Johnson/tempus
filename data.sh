#!/bin/sh
set -x

curl -X POST \
  http://localhost:8080/api/v1/categories \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Technique",
    "description": "Broad technical ability"
  }' | jq

curl -X POST \
  http://localhost:8080/api/v1/categories \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Independence",
    "description": "Independence exercises"
  }' | jq

curl -X POST \
  http://localhost:8080/api/v1/categories \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Songs & Creativity",
    "description": "Song practice, general creativity"
  }' | jq

# curl -X PATCH \
#   http://localhost:8080/api/v1/categories/1 \
#   -H 'Content-Type: application/json' \
#   -d '{
#   "category": {
#     "name": "Rudiments",
#     "description": "Standard drum rudiments and technique"
#   },
#   "update_mask": "name,description"
#   }'| jq

curl -X POST \
  http://localhost:8080/api/v1/tags \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Drumeo course 1.1"
  }' | jq

curl -X POST \
  http://localhost:8080/api/v1/tags \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Feet",
    "category_ids": [1]
  }' | jq

curl -X POST \
  http://localhost:8080/api/v1/tags \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Grooves",
    "category_ids": [2,3]
  }' | jq

curl -X POST \
  http://localhost:8080/api/v1/exercises \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Single Paradiddle",
    "description": "A basic drum rudiment pattern: RLRR LRLL",
    "tag_ids": [],
    "category_ids": [1]
  }' | jq

curl -X GET http://localhost:8080/api/v1/exercises/1 | jq

# List all exercises
curl -X GET http://localhost:8080/api/v1/exercises | jq

# List with pagination
curl -X GET "http://localhost:8080/api/v1/exercises?page_size=5" | jq

# Filter by category
curl -X GET "http://localhost:8080/api/v1/exercises?category_id=1" | jq

# Filter by tag
curl -X GET "http://localhost:8080/api/v1/exercises?tag_id=1"  | jq

# Filter by both category and tag
curl -X GET "http://localhost:8080/api/v1/exercises?category_id=1&tag_id=1" | jq

# Update the name and description
# curl -X PATCH \
#   http://localhost:8080/api/v1/exercises/1 \
#   -H 'Content-Type: application/json' \
#   -d '{
#     "exercise": {
#       "name": "Single Paradiddle - Basic",
#       "description": "One of the essential rudiments: RLRR LRLL"
#     },
#     "update_mask": "name,description"
#   }' | jq

# Update tags
# curl -X PATCH \
#   http://localhost:8080/api/v1/exercises/1 \
#   -H 'Content-Type: application/json' \
#   -d '{
#     "exercise": {
#       "tag_ids": [1, 2]
#     },
#     "update_mask": "tagIds"
#   }' | jq

# This is a simplified example as binary data requires proper handling
# curl -X POST \
#   http://localhost:8080/api/v1/exercises/1/images \
#   -H 'Content-Type: application/json' \
#   -d '{
#     "exercise_id": 1,
#     "image_data": "BASE64_ENCODED_IMAGE_DATA",
#     "filename": "paradiddle.png",
#     "mime_type": "image/png",
#     "description": "Pattern notation"
#   }' | jq

# curl -X DELETE http://localhost:8080/api/v1/exercise-images/1 | jq

curl -X POST \
  http://localhost:8080/api/v1/exercises/1/links \
  -H 'Content-Type: application/json' \
  -d '{
    "exercise_id": 1,
    "url": "https://example.com/tutorials/paradiddle",
    "description": "Tutorial video"
  }' | jq

# curl -X DELETE http://localhost:8080/api/v1/exercise-links/1 | jq

curl -X GET "http://localhost:8080/api/v1/exercises/1/stats" | jq

# curl -X DELETE http://localhost:8080/api/v1/exercises/1 | jq

curl -X POST \
  http://localhost:8080/api/v1/sessions \
  -H 'Content-Type: application/json' \
  -d '{
    "start_time": "2025-04-10T10:00:00Z",
    "end_time": "2025-04-10T11:30:00Z",
    "notes": "Morning practice session"
  }' | jq

curl -X POST \
  http://localhost:8080/api/v1/sessions \
  -H 'Content-Type: application/json' \
  -d '{
    "start_time": "2025-04-10T18:00:00Z",
    "end_time": "2025-04-10T20:30:00Z",
    "notes": "Evening practice session"
  }' | jq

curl -X GET http://localhost:8080/api/v1/sessions/1 | jq

# List all sessions
curl -X GET http://localhost:8080/api/v1/sessions | jq

# With pagination
curl -X GET "http://localhost:8080/api/v1/sessions?page_size=5" | jq

# Filter by date range
curl -X GET "http://localhost:8080/api/v1/sessions?start_date=2025-04-01T00:00:00Z&end_date=2025-04-30T23:59:59Z" | jq

# Filter by exercise ID
curl -X GET "http://localhost:8080/api/v1/sessions?exercise_id=1" | jq

# Update session notes
curl -X PATCH \
  http://localhost:8080/api/v1/sessions/1 \
  -H 'Content-Type: application/json' \
  -d '{
    "session": {
      "notes": "Morning practice session - focused on rudiments"
    },
    "update_mask": "notes"
  }' | jq

# Update start and end times
curl -X PATCH \
  http://localhost:8080/api/v1/sessions/1 \
  -H 'Content-Type: application/json' \
  -d '{
    "session": {
      "startTime": "2025-04-10T09:45:00Z",
      "endTime": "2025-04-10T11:45:02Z"
    },
    "update_mask": "startTime,endTime"
  }' | jq

# Delete session with ID 1
# curl -X DELETE http://localhost:8080/api/v1/sessions/1 | jq

# Create a history entry linked to a practice session
curl -X POST \
  http://localhost:8080/api/v1/history \
  -H 'Content-Type: application/json' \
  -d '{
    "exercise_id": 1,
    "session_id": 1,
    "start_time": "2025-04-10T10:00:00Z",
    "end_time": "2025-04-10T10:15:00Z",
    "bpms": [85,90],
    "time_signature": "4/4",
    "notes": "Practiced during my morning session",
    "rating": 4
  }' | jq

# List all history entries
curl -X GET http://localhost:8080/api/v1/history | jq

# Filter by exercise
curl -X GET "http://localhost:8080/api/v1/history?exercise_id=1" | jq

# Filter by session (assuming you've added this query parameter)
curl -X GET "http://localhost:8080/api/v1/history?session_id=1" | jq

# Filter by date range
curl -X GET "http://localhost:8080/api/v1/history?start_date=2025-04-01T00:00:00Z&end_date=2025-04-30T23:59:59Z" | jq

# Combined filters with pagination
curl -X GET "http://localhost:8080/api/v1/history?exercise_id=1&session_id=1&page_size=10" | jq

# Update BPM, rating
curl -X PATCH \
  http://localhost:8080/api/v1/history/1 \
  -H 'Content-Type: application/json' \
  -d '{
    "history": {
      "bpms": [95],
      "rating": 5
    },
    "update_mask": "bpms,rating"
  }' | jq

# curl -X DELETE http://localhost:8080/api/v1/history/1 | jq
