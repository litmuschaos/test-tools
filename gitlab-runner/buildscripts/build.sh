#!/bin/bash

set -e

CHANGED_FILES=`git diff --name-only master...${TRAVIS_COMMIT}`
INFRA=True
DIR="gitlab-runner"

for CHANGED_FILE in $CHANGED_FILES; do
  if [[ "$CHANGED_FILE" =~ $DIR ]]; then
    INFRA=False
    break
  fi
done

if [[ $INFRA == True ]]; then
  echo "NOT building Gitlab-runner infra image."
  
else
  echo "building Gitlab-runner infra image."
  cd ..
  docker build -t atulabhi/kops:v22 .
fi
