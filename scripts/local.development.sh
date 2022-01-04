#!/bin/bash

export COMPOSE_IGNORE_ORPHANS=True

export BACKEND_IMAGE=learn-go-restful-api-backend-go
export BACKEND_IMAGE_TAG=development
export BACKEND_CONTAINER=learn-go-restful-api-backend-go-development
export BACKEND_HOST=learn-go-restful-api-backend-go.service
export BACKEND_STAGE=development

docker build -t "$BACKEND_IMAGE:$BACKEND_IMAGE_TAG" .
docker-compose -f ./manifest/docker-compose.development.yaml up -d --build
