#!/bin/bash

export COMPOSE_IGNORE_ORPHANS=True

export BACKEND_IMAGE=learn-go-restful-api-backend-go
export BACKEND_IMAGE_TAG=production
export BACKEND_CONTAINER=learn-go-restful-api-backend-go-production
export BACKEND_HOST=learn-go-restful-api-backend-go.service
export BACKEND_STAGE=production

docker build -t "$BACKEND_IMAGE:$BACKEND_IMAGE_TAG" .
docker-compose -f ./manifest/docker-compose.production.yaml up -d --build
