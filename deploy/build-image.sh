#!/bin/sh -x

REGION=$1
ACCOUNT=$2
REPOSITORY=$3
PROJECT=$4
HASH=$5
BUILD=$6

REPO_URL="$ACCOUNT.dkr.ecr.$REGION.amazonaws.com/$REPOSITORY"

aws ecr get-login-password --region "$REGION" \
| docker login \
  -u AWS \
  --password-stdin \
  "$REPO_URL"

TAG="$PROJECT":"$HASH"-"$BUILD"

ls -las

docker build -t "$REPO_URL"/"$TAG" -f deploy/"$PROJECT"/Dockerfile .

docker push "$REPO_URL"/"$TAG"

