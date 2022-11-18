#!/bin/sh -x

REGION=$1
ACCOUNT=$2
SCOPE=$3
PROJECT=$4
HASH=$5
BUILD=$6

REGISTRY="$ACCOUNT.dkr.ecr.$REGION.amazonaws.com"
REPOSITORY="$SCOPE"/"$PROJECT"
REPOSITORY_URL="$REGISTRY"/"$REPOSITORY"

aws ecr get-login-password --region "$REGION" \
| docker login \
  -u AWS \
  --password-stdin \
  "$REGISTRY"

TAG="$HASH"-"$BUILD"

ls -las

docker build -t "$REPOSITORY_URL":"$TAG" -f deploy/"$PROJECT"/Dockerfile .

docker push "$REPOSITORY_URL":"$TAG"

