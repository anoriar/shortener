#!/usr/bin/make
include .env
export $(shell sed 's/=.*//' .env)


.PHONY: start
start:
	docker-compose -f docker-compose.yaml up -d --force-recreate

.PHONY: fmt-import
fmt-import:
	goimports -local github.com/anoriar/shortener -w ./

.PHONY: generate-proto
generate-proto:
	./proto/generate-proto

.PHONY: godoc-generate
godoc-generate:
	godoc-generate

