.ONESHELL:
SHELL=/bin/bash
-include .env

init:
	cp -n .env.example .env
	go mod download

compile-proto:
	rm -rf service-definitions
	export GOPRIVATE=github.com/vibeitco && \
		go get -u github.com/vibeitco/service-definitions/ && \
		git clone git@github.com:vibeitco/service-definitions.git && \
		cd service-definitions && \
		make model i=service/${SERVICE}/v1/model.proto o=. && \
		make gateway i=service/${SERVICE}/v1/model.proto o=.
	cp service-definitions/github.com/vibeitco/${SERVICE}-service/model/* ./model/
	rm -rf service-definitions
	make fix-proto

fix-proto:
	sed -i '' 's/json:"id,omitempty"/json:"id,omitempty" bson:"_id"/g'  ./model/model.pb.go

update:
	export GOPRIVATE=github.com/vibeitco && \
		go get -u ./...

build:
	go build -o svc

build-target:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./svc .

build-docker: build-target
	docker build --tag vibeitco/${SERVICE}:${VERSION} .

run: build
	./svc

run-docker: build-docker
	docker run --rm \
		-e SERVICE=${SERVICE} \
		-e ENV=${ENV} \
		-e VERSION=${VERSION} \
		-e MONGODB_PASSWORD=${MONGODB_PASSWORD} \
		-e DOMAIN=${DOMAIN} \
		-e DOMAIN_API=${DOMAIN_API} \
		-e DOMAIN_ASSETS=${DOMAIN_ASSETS} \
		-p 2020:2020 \
		-p 8080:8080 \
		vibeitco/${SERVICE}:${VERSION}
