#!/bin/sh
set -e

if [ -z "$TODO_BACKEND_URL" ]; then
    echo "Missing TODO_BACKEND_URL"
    exit 1
fi

article_url="$(
  wget -SO /dev/null "https://en.wikipedia.org/wiki/Special:Random" 2>&1 | awk '/^[[:space:]]*[Ll]ocation:/ {print "https:" $2}'
)"

if [ -z "$article_url" ]; then
  echo "Failed to get random Wikipedia article"
  exit 1
fi

wget -O /dev/null --post-data="title=Read article $article_url" "$TODO_BACKEND_URL/todos"
